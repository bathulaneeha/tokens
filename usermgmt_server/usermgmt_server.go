package main

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"

	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
)

// structure of token
type Token struct {
	Id                string
	Name              string
	DomainLow         uint64
	DomainMid         uint64
	DomainHigh        uint64
	StatePartialValue uint64
	StateFinalValue   uint64
	Message           string
}

// defining variables requried
var (
	port   = flag.Int("port", 5030, "The server port")
	tokens = []Token{}
)

type server struct {
	pb.UnimplementedTokenServiceServer
}

// gives the minimum x value based on the intervals
func argmin_x(name string, n1 uint64, n2 uint64) uint64 {
	k := Hash(name, n1)
	l := n1
	for i := 0; n1 < n2-1; i++ {
		n1 = n1 + 1
		temp := Hash(name, n1)
		if temp < k {
			k = temp
			l = n1
		}
	}
	return l
}

// hash function to calculate hash value for name
func Hash(name string, nonce uint64) uint64 {
	hasher := sha256.New()
	hasher.Write([]byte(fmt.Sprintf("%s %d", name, nonce)))
	return binary.BigEndian.Uint64(hasher.Sum(nil))
}

// printing the list of tokens
func GetAllTokens(ctx context.Context) {
	log.Printf("Current tokens is:")
	for _, token := range tokens {
		fmt.Printf("ID: %v; \n	Name: %v; \n	DomainLow: %v; \n	DomainMid: %v; \n	DomainHigh: %v; \n	StatePartialValue: %v; \n	StateFinalValue: %v \n\n", token.Id, token.Name, token.DomainLow, token.DomainMid, token.DomainHigh, token.StatePartialValue, token.StateFinalValue)
	}
}

func GetAllTokensIds(ctx context.Context) {
	log.Printf("All tokens are:")
	for _, token := range tokens {
		fmt.Printf("ID: %v\n",token.Id)
	}
}

// function to create tokens
func (s *server) CreateToken(ctx context.Context, in *pb.Token) (*pb.Token, error) {
	fmt.Printf("Recieved method call:\n to create ID: %v\n", in.GetId())
	oldleng := len(tokens)
	for i := 0; i < oldleng; i++ {
		if tokens[i].Id == in.GetId() {
			log.Printf("Creating ID " + in.GetId() + " failed as it already exists")
			fmt.Printf("\n")
			return nil, errors.New("id " + in.GetId() + " already exists")
		}
	}
	onetoken := Token{
		Id: in.GetId(),
	}
	tokens = append(tokens, onetoken)
	log.Printf("Creating ID " + in.GetId() + " successed.\n")
	GetAllTokens(ctx)
	GetAllTokensIds(ctx)
	return &pb.Token{
		Id:      tokens[len(tokens)-1].Id,
		Message: "Created ID " + in.GetId() + " successed.",
	}, nil
}

// function to write on tokens
func (s *server) WriteToken(ctx context.Context, in *pb.Token) (*pb.Token, error) {
	log.Printf("Recieved method call:\n to write ID: %v,\nName: %v, \nDomainLow: %v, \nDomainMid: %v, \nDomainHigh: %v\n", in.GetId(), in.GetName(), in.GetDomainLow(), in.GetDomainMid(), in.GetDomainHigh())
	oldleng := len(tokens)
	for i := 0; i < oldleng; i++ {
		if tokens[i].Id == in.GetId() {
			tokens[i].Name = in.GetName()
			tokens[i].DomainLow = in.GetDomainLow()
			tokens[i].DomainMid = in.GetDomainMid()
			tokens[i].DomainHigh = in.GetDomainHigh()
			tokens[i].StatePartialValue = argmin_x(in.GetName(), in.GetDomainLow(), in.GetDomainMid())
			tokens[i].StateFinalValue = 0
			log.Printf("Writing ID " + in.GetId() + " successed.\n")
			GetAllTokens(ctx)
			return &pb.Token{
				Id:                tokens[i].Id,
				Name:              tokens[i].Name,
				DomainLow:         tokens[i].DomainLow,
				DomainMid:         tokens[i].DomainMid,
				DomainHigh:        tokens[i].DomainHigh,
				StatePartialValue: tokens[i].StatePartialValue,
				StateFinalValue:   tokens[i].StateFinalValue,
				Message:           "Wrote ID " + tokens[i].Id + " successed.",
			}, nil
		}
	}
	log.Printf("Writing ID " + in.GetId() + " failed.\n")
	GetAllTokens(ctx)
	return nil, errors.New("id " + in.GetId() + " was not found")
}

// function to read tokens
func (s *server) ReadToken(ctx context.Context, in *pb.Token) (*pb.Token, error) {
	log.Printf("Recieved method call:\n to read ID: %v\n", in.GetId())
	oldleng := len(tokens)
	for i := 0; i < oldleng; i++ {
		if tokens[i].Id == in.GetId() {

			tempfinal := argmin_x(tokens[i].Name, tokens[i].DomainMid, tokens[i].DomainHigh)
			tokens[i].StateFinalValue = tokens[i].StatePartialValue
			if tempfinal <= tokens[i].StatePartialValue {
				tokens[i].StateFinalValue = tempfinal
			}

			log.Printf("Reading ID " + in.GetId() + " successed and StateFinalValue updated.\n")
			GetAllTokens(ctx)
			return &pb.Token{
				Id:                tokens[i].Id,
				Name:              tokens[i].Name,
				DomainLow:         tokens[i].DomainLow,
				DomainMid:         tokens[i].DomainMid,
				DomainHigh:        tokens[i].DomainHigh,
				StatePartialValue: tokens[i].StatePartialValue,
				StateFinalValue:   tokens[i].StateFinalValue,
				Message:           "Read ID " + in.GetId() + " successed.",
			}, nil
		}
	}
	log.Printf("Reading ID " + in.GetId() + " failed.\n")
	fmt.Println()
	return nil, errors.New("id " + in.GetId() + " was not found")
}

// function to drop token
func (s *server) DropToken(ctx context.Context, in *pb.Token) (*pb.Token, error) {
	log.Printf("Recieved method call:\n to drop ID: %v\n", in.GetId())
	oldleng := len(tokens)
	for i := 0; i < oldleng; i++ {
		if tokens[i].Id == in.GetId() {
			tokens[i] = tokens[oldleng-1]
			tokens = tokens[:oldleng-1]

			log.Printf("Dropping ID " + in.GetId() + " successed.\n")
			GetAllTokens(ctx)

			return &pb.Token{
				Message: "Dropped " + in.GetId() + " successed",
			}, nil
		}
	}
	log.Printf("Dropping ID " + in.GetId() + " failed.\n")
	GetAllTokens(ctx)
	return nil, errors.New("id " + in.GetId() + " was not found")
}

// main function to start execution
func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Server: failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTokenServiceServer(s, &server{})
	log.Printf("Server: listening at port %v\n\n", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Server: failed to serve: %v", err)
	}
}
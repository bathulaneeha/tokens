package main

//libraries imported
import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"
	pb "example.com/go-usermgmt-grpc/usermgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//initial values of all the pointers are assigned as
var (
	createP = flag.Bool("create", false, "create operation pointer")
	readP   = flag.Bool("read", false, "read operation pointer")
	writeP  = flag.Bool("write", false, "write operation pointer")
	dropP   = flag.Bool("drop", false, "drop operation pointer")
	idP     = flag.String("id", "-1", "id pointer of token")
	nameP   = flag.String("name", "", "name pointer of token")
	lowP    = flag.Uint64("low", 0, "The low value  of your token")
	midP    = flag.Uint64("mid", 0, "The mid value  of your token")
	highP   = flag.Uint64("high", 0, "The high value  of your token")
	hostP   = flag.String("host", "localhost", "The host (string) to connect to. Default is localhost")
	portP   = flag.String("port", "50051", "The port to connect to")
)
//function to check for flags
func checkflags() {
	if *createP {
		if *readP || *writeP|| *dropP {
			log.Fatalf("Client: Flag error. Only one operation is allowed at same time.")
		}
	} else if *readP {
		if *createP || *writeP|| *dropP {
			log.Fatalf("Client: Flag error. Only one operation is allowed at same time.")
		}
	} else if *writeP{
		if *createP || *readP || *dropP {
			log.Fatalf("Client: Flag error. Only one operation is allowed at same time.")
		}
	} else if *dropP {
		if *createP || *readP || *writeP{
			log.Fatalf("Client: Flag error. Only one operation is allowed at same time.")
		}
	}

}

func main() {
	flag.Parse()
	log.Printf("Your input is:")
	if *createP{
	fmt.Println("Operation:\nCreate: ")
	} else if *readP{
		fmt.Println("Operation:\nRead: ")
	}else if *writeP{
		fmt.Println("Operation:\nWrite: ")
	}else if *dropP{
		fmt.Println("Operation:\nDrop: ")
	}else{
		fmt.Println("INVALID OPERATION")
	}
	fmt.Println("Parameters: ID:", *idP, ";\n Host:", *hostP, ";\n Port: ", *portP, ";\n Low: ", *lowP, ";\n Mid: ", *midP, "; \nHigh: ", *highP, ";\n Name: ", *nameP)
	fmt.Println()
	checkflags()
	//create connection with server
	addr := *hostP + ":" + *portP
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Client: did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewTokenServiceClient(conn)
	// Call server methods to deal with token
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if *createP {
		r, err := c.CreateToken(ctx, &pb.Token{
			Id: *idP,
		})
		if err != nil {
			log.Fatalf("Client: failed to call server CreateToken(): %v", err)
		}
		log.Printf("Client received: \nID: %v, Message: %v\n", r.GetId(), r.GetMessage())
	}

	//based on method calls, paramaters are passed
	if *writeP {
		r, err := c.WriteToken(ctx, &pb.Token{
			Id:         *idP,
			Name:       *nameP,
			DomainLow:  *lowP,
			DomainMid:  *midP,
			DomainHigh: *highP,
		})
		if err != nil {
			log.Fatalf("Client: failed to call server WriteToken(): %v", err)
		}
		log.Printf("Client received: Message: %v, \nID: %v,\nName: %s, \nDomainLow: %v, \nDomainMid: %v, \nDomainHigh: %v, \nStatePartialValue: %v, \nStateFinalValue: %v", r.GetMessage(), r.GetId(), r.GetName(), r.GetDomainLow(), r.GetDomainMid(), r.GetDomainHigh(), r.GetStatePartialValue(), r.GetStateFinalValue())
	}else if *readP {
		r, err := c.ReadToken(ctx, &pb.Token{
			Id: *idP,
		})
		if err != nil {
			log.Fatalf("Client: failed to call server ReadToken(): %v", err)
		}
		log.Printf("Client received: Message: %v, \nID: %v,\nName: %s, \nDomainLow: %v, \nDomainMid: %v, \nDomainHigh: %v, \nStatePartialValue: %v, \nStateFinalValue: %v", r.GetMessage(), r.GetId(), r.GetName(), r.GetDomainLow(), r.GetDomainMid(), r.GetDomainHigh(), r.GetStatePartialValue(), r.GetStateFinalValue())
	}else if *dropP {
		r, err := c.DropToken(ctx, &pb.Token{
			Id: *idP,
		})
		if err != nil {
			log.Fatalf("Client: failed to call server DropToken(): %v", err)
		}
		log.Printf("Client received: %v \n", r.GetMessage())
	}
}
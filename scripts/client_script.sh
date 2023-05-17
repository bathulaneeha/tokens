go run usermgmt_client/usermgmt_client.go -create -id 1234 -host localhost -port 50051

go run usermgmt_client/usermgmt_client.go t -write -id 1234 -name abc -low 0 -mid 10 -high 100 -host localhost -port 50051

go run usermgmt_client/usermgmt_client.go -read -id 1234 -host localhost -port 50051

go run usermgmt_client/usermgmt_client.go -drop 1234 -host localhost -port 50051
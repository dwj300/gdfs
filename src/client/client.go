package main

import (
	"flag"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "gdfs/src/contracts"
)

var serverAddr = flag.String("server_addr", "127.0.0.1:9000", "The server address in the format of host:port")

func main() {
    	grpclog.Println("Client starting!")
	flag.Parse()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	client := pb.NewBlobsClient(conn)
	data, err := client.ReadBlob(context.Background(), &pb.Filename{Filename:"dummy"})
	grpclog.Println(string(data.Data))
	data1, err := client.ListBlobs(context.Background(), &pb.Empty{})
	grpclog.Println(data1.String())
}

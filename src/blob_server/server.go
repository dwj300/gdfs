package main

import (
	"crypto/rand"
	"io"
	"fmt"
	"os"
	"io/ioutil"
	"errors"
	"flag"
	"net"
	"log"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "gdfs/src/contracts"
x)

var (
	dataDir = flag.String("dataDir", "/tmp/data/", "Data directory for blob server")
	port = flag.Int("port", 9000, "The server port")
)

type blobServer struct {}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (s *blobServer) CreateBlob(ctx context.Context, putData *pb.PutData) (*pb.Empty, error) {
	grpclog.Println("CreateBlob called!")
	name := *dataDir + putData.Filename
	if _, err := os.Stat(name); err == nil {
		var err = errors.New("blob server: file already exists")
		return &pb.Empty{}, err
	}
	f, err := os.Create(name)
	check(err)
	defer f.Close()
	_, err = f.Write(putData.Data)
	check(err)
	return &pb.Empty{}, err
}

func (s *blobServer) ReadBlob(ctx context.Context, filename *pb.Filename) (*pb.Data, error) {
	grpclog.Println("ReadBlob called!")
	name := *dataDir + filename.Filename
	if _, err := os.Stat(name); err == nil {
		dat, err := ioutil.ReadFile(name)
		check(err)
		return &pb.Data{dat}, nil
	}
	var err = errors.New("blob server: file not found")
	return &pb.Data{}, err
}

func (s *blobServer) UpdateBlob(ctx context.Context, putData *pb.PutData) (*pb.Empty, error) {
	grpclog.Println("UpdateBlob called!")
	name := *dataDir + putData.Filename
	if _, err := os.Stat(name); err == nil {
		f, err := os.OpenFile(name, os.O_WRONLY, 0644)
		_, err = f.Write(putData.Data)
		return &pb.Empty{}, err
	}
	var err = errors.New("blob server: file doesn't exist")
	return &pb.Empty{}, err
}

func (s *blobServer) DeleteBlob(ctx context.Context, filename *pb.Filename) (*pb.Empty, error) {
	grpclog.Println("Delete called!")
	name := *dataDir + filename.Filename
	if _, err := os.Stat(name); err == nil {
		err = os.Remove(name)
		return &pb.Empty{}, err
	}
	var err = errors.New("blob server: file doesn't exist")
	return &pb.Empty{}, err
}

func (s *blobServer) ListBlobs(ctx context.Context, empty *pb.Empty) (*pb.BlobList, error) {
	grpclog.Println("ListBlobs called!")
	files, _ := ioutil.ReadDir(*dataDir)
	blobList := new(pb.BlobList)
	for _, f := range files {
		blob := &pb.Blob{Filename:f.Name()}
		blobList.Blobs = append(blobList.Blobs, blob)
	}
	return blobList, nil
}

func main() {
	//grpclog.Println(newUUID())
	grpclog.Println("Blob server starting!")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterBlobsServer(grpcServer, new(blobServer))
	grpcServer.Serve(lis)
}

// newUUID generates a random UUID according to RFC 4122
func newUUID() (string) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
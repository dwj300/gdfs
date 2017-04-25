package blob_server

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
	"github.com/golang/protobuf/protoc-gen-go/grpc"
)

var (
	dataDir = flag.String("dataDir", "/tmp/data/", "Data directory for blob server")
	port       = flag.Int("port", 9000, "The server port")
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func (s *blobsServer) GetBlob(ctx context.Context, filename *pb.Filename) (*pb.Data, error) {
	name := dataDir + filename.filename
	if _, err := os.Stat(name); err == nil {
		dat, err := ioutil.ReadFile(name)
		check(err)
		return &pb.Data{dat}, nil
	}
	var err = errors.New("blob server: file not found")
	return &pb.Data{}, nil
}

func main() {
	fmt.Println(newUUID())
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, new(BlobsServer))
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
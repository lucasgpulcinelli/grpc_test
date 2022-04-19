package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	f "github.com/lucasgpulcinelli/grpc_test/functions"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50501, "the port to listen")
)

type server struct {
	f.UnimplementedEchoServer
}

func (s *server) EchoStr(ctx context.Context, req *f.EchoRequest) (*f.EchoReply, error) {
	log.Printf("Received: %v", req.GetStr())
	return &f.EchoReply{Str: "Hello " + req.GetStr()}, nil
}

func main() {
	flag.Parse()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("Error in Listen: %v", err)
	}

	s := grpc.NewServer()
	f.RegisterEchoServer(s, &server{})
	log.Printf("listening at: %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

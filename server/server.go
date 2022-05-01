package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"

	f "github.com/lucasgpulcinelli/grpc_test/functions"
	"google.golang.org/grpc"
)

type server struct {
	f.UnimplementedEchoServer
}

func (s *server) EchoStr(ctx context.Context, req *f.EchoData) (*f.EchoData, error) {
	log.Printf("Received Echo: %v", req.GetStr())
	return &f.EchoData{Str: "Hello " + req.GetStr()}, nil
}

func (s *server) EchoCounter(req *f.EchoData, stream f.Echo_EchoCounterServer) error {
	req_str := req.GetStr()
	log.Printf("Received Counter: %v", req_str)

	var res_str string
	var n int
	_, err := fmt.Sscanf(req_str, "%d %s", &n, &res_str)

	if err != nil {
		log.Printf("[WARNING] counter: %v", err)
		return err
	}

	for i := 0; i < n; i++ {
		res := f.EchoData{Str: fmt.Sprintf("%d Hello %s", i, res_str)}

		if err := stream.Send(&res); err != nil {
			log.Printf("[WARNING] counter: %v", err)
			return err
		}
	}

	return nil
}

func (s *server) ConcatEchos(stream f.Echo_ConcatEchosServer) error {
	log.Printf("Received Concat")

	var res_str string
	for {
		req, err := stream.Recv()

		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("[WARNING] Concat: %v", err)
			return err
		}

		res_str += req.GetStr()
	}

	return stream.Send(&f.EchoData{Str: res_str})
}

func min(v1 int, v2 int) int {
	if v1 < v2 {
		return v1
	}
	return v2
}

func (s *server) PermuteEcho(stream f.Echo_PermuteEchoServer) error {
	log.Printf("Received Permute")

	var res_str string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Printf("[WARNING] Recv permute: %v", err)
			return err
		}

		req_str := req.GetStr()

		var new_res_str string

		minlen := min(len(req_str), len(res_str))
		for i := 0; i < minlen; i++ {
			new_res_str += string(res_str[i])
			new_res_str += string(req_str[i])
		}

		if minlen == len(req_str) {
			new_res_str += res_str[minlen:]
		} else {
			new_res_str += req_str[minlen:]
		}

		res_str = new_res_str

		stream.Send(&f.EchoData{Str: res_str})
	}
}

func main() {

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "50501"
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
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

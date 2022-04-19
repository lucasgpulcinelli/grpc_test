package main

import (
	"context"
	"flag"
	"log"
	"time"

	f "github.com/lucasgpulcinelli/grpc_test/functions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:50501", "the address to connect to")
	str  = flag.String("str", "world", "str to echo")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := f.NewEchoClient(conn)

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := c.EchoStr(ctx, &f.EchoRequest{Str: *str})
		if err != nil {
			log.Fatalf("could not echo: %v", err)
		}
		log.Printf("Echoed: %s", r.GetStr())
		time.Sleep(time.Second)
	}
}

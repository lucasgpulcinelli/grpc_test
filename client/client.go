package main

import (
	"context"
	"flag"
	"io"
	"log"
	"os"
	"time"

	f "github.com/lucasgpulcinelli/grpc_test/functions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	echo_type = flag.String("echotype", "simple", "type of echo")
)

func EchoStr(c f.EchoClient, str string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := c.EchoStr(ctx, &f.EchoData{Str: str})
	if err != nil {
		log.Fatalf("Could not echo: %v", err)
	}
	log.Printf("Echoed: %s", res.GetStr())
}

func EchoCounter(c f.EchoClient, str string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := c.EchoCounter(ctx, &f.EchoData{Str: str})
	if err != nil {
		log.Fatalf("Could not count: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Could not count all: %v", err)
		}

		log.Printf("Counted: %s", res.GetStr())
	}
}

func ConcatEchos(c f.EchoClient, str_arr []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := c.ConcatEchos(ctx)
	if err != nil {
		log.Fatalf("Could not start concat: %v", err)
	}

	for _, s := range str_arr {
		err = stream.Send(&f.EchoData{Str: s})

		if err != nil {
			log.Fatalf("Could not send concat: %v", err)
		}
	}
	stream.CloseSend()

	res, err := stream.Recv()
	if err != nil {
		log.Fatalf("Could not get concat: %v", err)
	}
	log.Printf("Concated: %s", res.GetStr())
}

func PermuteEcho(c f.EchoClient, str_arr []string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	stream, err := c.PermuteEcho(ctx)
	if err != nil {
		log.Fatalf("Could not start permute: %v", err)
	}

	go func() {
		for _, s := range str_arr {
			err = stream.Send(&f.EchoData{Str: s})

			if err != nil {
				log.Fatalf("Could not send permute: %v", err)
			}
		}
		stream.CloseSend()
	}()

	for range str_arr {
		res, err := stream.Recv()
		if err != nil {
			log.Fatalf("Could not get permute: %v", err)
		}
		log.Printf("Permuted: %s", res.GetStr())
	}

}

func main() {
	flag.Parse()
	addr, exists := os.LookupEnv("ADDR")
	if !exists {
		addr = "localhost:50501"
	}

	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()
	c := f.NewEchoClient(conn)

	switch *echo_type {
	case "simple":
		EchoStr(c, flag.Arg(0))
	case "count":
		EchoCounter(c, flag.Arg(0))
	case "concat":
		ConcatEchos(c, flag.Args())
	case "permute":
		PermuteEcho(c, flag.Args())
	case "all":
		for {
			time.Sleep(time.Second)
			go EchoStr(c, flag.Arg(0))
			go EchoCounter(c, flag.Arg(1))
			go ConcatEchos(c, flag.Args()[2:])
			go PermuteEcho(c, flag.Args()[2:])
		}
	default:
		log.Fatalf("Unknown type of echo: \"%s\"", *echo_type)
	}
}

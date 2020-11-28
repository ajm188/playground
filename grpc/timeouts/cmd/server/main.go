package main

import (
	"flag"
	"net"

	"github.com/ajm188/playground/grpc/timeouts/proto/example"
	"github.com/ajm188/playground/grpc/timeouts/server"

	"google.golang.org/grpc"
)

func main() {
	addr := flag.String("addr", ":8080", "address to listen on")
	delay := flag.Duration("delay", 0, "delay to inject in responses")

	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	serv := grpc.NewServer()
	example.RegisterHelloServer(serv, &server.ExampleServer{Delay: *delay})

	if err := serv.Serve(lis); err != nil {
		panic(err)
	}
}

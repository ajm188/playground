package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/ajm188/playground/grpc/timeouts/proto/example"
	"google.golang.org/grpc"
)

type Client struct {
	example.HelloClient
	timeout time.Duration
}

func main() {
	addr := flag.String("addr", ":8081", "address to listen for requests to forward to the server")
	serverAddr := flag.String("server-addr", ":8080", "address to talk to the server on")
	timeout := flag.Duration("timeout", time.Millisecond*50, "how long a timeout to set on gRPC calls")

	flag.Parse()

	gconn, err := grpc.Dial(*serverAddr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer gconn.Close()

	client := &Client{
		HelloClient: example.NewHelloClient(gconn),
		timeout:     *timeout,
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Printf("cannot open listener on %s: err = %s", *addr, err)
		return
	}

	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("error accepting new conn: %s", err)
			return
		}

		go handle(client, conn)
	}
}

func handle(client *Client, conn net.Conn) {
	defer conn.Close()

	var buf bytes.Buffer
	io.Copy(&buf, conn)

	ctx, cancel := context.WithTimeout(context.Background(), client.timeout)
	defer cancel()

	name := buf.String()
	buf.Reset()

	resp, err := client.Hello(ctx, &example.HelloRequest{Name: name})
	if err != nil {
		buf.WriteString(fmt.Sprintf("error forwarding hello request %q; err = %q\n", name, err))
		io.Copy(os.Stdout, &buf)

		return
	}

	buf.WriteString(fmt.Sprintf("Hello, %s\n", resp.Name))
	io.Copy(os.Stdout, &buf)
}

package server

import (
	"context"
	"log"
	"time"

	"github.com/ajm188/playground/grpc/timeouts/proto/example"
)

type ExampleServer struct {
	Delay time.Duration
}

var _ example.HelloServer = (*ExampleServer)(nil)

func (s *ExampleServer) Hello(ctx context.Context, req *example.HelloRequest) (*example.HelloResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	start := time.Now()

	log.Printf("received request for %s", req.Name)
	select {
	case <-ctx.Done():
		log.Printf("context timed out after %s", time.Since(start))

		return nil, ctx.Err()
	case <-time.After(s.Delay):
		log.Printf("sending response after delay %s", s.Delay)

		return &example.HelloResponse{Name: req.Name}, nil
	}
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hsmtkk/cuddly-waffle/count"
	"github.com/hsmtkk/cuddly-waffle/env"
	"google.golang.org/grpc"
)

type server struct {
	count.UnimplementedCounterServer
	countChan <-chan int64
}

func newServer(countChan <-chan int64) *server {
	return &server{countChan: countChan}
}

func (s *server) Count(ctx context.Context, req *count.CountRequest) (*count.CountResponse, error) {
	log.Printf("request id:%d", req.GetId())
	cnt := <-s.countChan
	log.Printf("response count:%d", cnt)
	resp := &count.CountResponse{Count: cnt}
	return resp, nil
}

func main() {
	grpcPort := env.OptionalInt("GRPC_PORT", 50051)

	countChan := make(chan int64)
	go incrementer(countChan)

	addr := fmt.Sprintf("0.0.0.0:%d", grpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen; %s; %s", addr, err)
	}

	srv := grpc.NewServer()
	count.RegisterCounterServer(srv, newServer(countChan))
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to start server; %s", err)
	}
}

func incrementer(countChan chan<- int64) {
	var count int64 = 0
	for {
		countChan <- count
		count++
	}
}

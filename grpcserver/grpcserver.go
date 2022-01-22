package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/hsmtkk/cuddly-waffle/count"
	"github.com/hsmtkk/cuddly-waffle/env"
	"github.com/hsmtkk/cuddly-waffle/msg"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type server struct {
	count.UnimplementedCounterServer
	sugar       *zap.SugaredLogger
	countChan   <-chan int64
	natsConn    *nats.Conn
	natsChannel string
}

func newServer(sugar *zap.SugaredLogger, countChan <-chan int64, natsConn *nats.Conn, natsChannel string) *server {
	return &server{sugar: sugar, countChan: countChan, natsConn: natsConn, natsChannel: natsChannel}
}

func (s *server) Count(ctx context.Context, req *count.CountRequest) (*count.CountResponse, error) {
	id := req.GetId()
	s.sugar.Debugw("request", "id", id)
	cnt := <-s.countChan
	s.sugar.Debugw("response", "count", cnt)
	m, err := msg.NewMessage(id, cnt).ToJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to format message; %w", err)
	}
	if err := s.natsConn.Publish(s.natsChannel, m); err != nil {
		return nil, fmt.Errorf("failed to publish; %w", err)
	}
	resp := &count.CountResponse{Count: cnt}
	return resp, nil
}

func main() {
	grpcPort := env.OptionalInt("GRPC_PORT", 50051)
	natsHost := env.RequiredString("NATS_HOST")
	natsPort := env.RequiredInt("NATS_PORT")
	natsChannel := env.RequiredString("NATS_CHANNEL")

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger; %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	countChan := make(chan int64)
	go incrementer(countChan)

	natsAddr := fmt.Sprintf("nats://%s:%d", natsHost, natsPort)
	natsConn, err := nats.Connect(natsAddr)
	if err != nil {
		sugar.Fatalw("failed to connect NATS", "address", natsAddr, "error", err)
	}
	defer natsConn.Close()

	addr := fmt.Sprintf("0.0.0.0:%d", grpcPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		sugar.Fatalw("failed to listen", "address", addr, "error", err)
	}

	srv := grpc.NewServer()
	count.RegisterCounterServer(srv, newServer(sugar, countChan, natsConn, natsChannel))
	if err := srv.Serve(lis); err != nil {
		sugar.Fatalw("failed to start server", "error", err)
	}
}

func incrementer(countChan chan<- int64) {
	var count int64 = 0
	for {
		countChan <- count
		count++
	}
}

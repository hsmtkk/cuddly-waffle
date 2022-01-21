package main

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/hsmtkk/cuddly-waffle/count"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	users int
	iters int
	host  string
	port  uint16
)

func main() {
	command := &cobra.Command{
		Run: run,
	}
	command.Flags().IntVar(&users, "users", 100, "number of users")
	command.Flags().IntVar(&iters, "iters", 100, "number of iterations")
	command.Flags().StringVar(&host, "host", "127.0.0.1", "gRPC host")
	command.Flags().Uint16Var(&port, "port", 50051, "gRPC port")
	if err := command.Execute(); err != nil {
		log.Fatalf("failed to execute command; %s", err)
	}
}

func run(cmd *cobra.Command, args []string) {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect gRPC; %s; %s", addr, err)
	}
	defer conn.Close()
	clt := count.NewCounterClient(conn)

	var wg sync.WaitGroup
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			run2(clt, id)
		}(i)
	}
	wg.Wait()
}

func run2(clt count.CounterClient, id int) {
	for i := 0; i < iters; i++ {
		req := count.CountRequest{
			Id: int64(id),
		}
		resp, err := clt.Count(context.Background(), &req)
		if err != nil {
			log.Printf("failed to call gRPC; %s", err)
		}
		log.Printf("id:%d count:%d", id, resp.Count)
	}
}

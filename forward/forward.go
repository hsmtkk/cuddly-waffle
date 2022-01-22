package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/hsmtkk/cuddly-waffle/env"
	"github.com/hsmtkk/cuddly-waffle/msg"
	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func main() {
	natsHost := env.RequiredString("NATS_HOST")
	natsPort := env.RequiredInt("NATS_PORT")
	natsChannel := env.RequiredString("NATS_CHANNEL")
	redisHost := env.RequiredString("REDIS_HOST")
	redisPort := env.RequiredInt("REDIS_PORT")

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("failed to init logger; %s", err)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	redisAddr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer redisClient.Close()

	natsAddr := fmt.Sprintf("nats://%s:%d", natsHost, natsPort)
	natsConn, err := nats.Connect(natsAddr)
	if err != nil {
		sugar.Fatalw("failed to connect NATS", "address", natsAddr, "error", err)
	}
	defer natsConn.Close()

	handler := newMessageHandler(sugar, redisClient)

	sub, err := natsConn.Subscribe(natsChannel, handler.Handle)
	if err != nil {
		sugar.Fatalw("failed to subscribe channel", "channel", natsChannel, "error", err)
	}
	defer sub.Unsubscribe()

	select {}
}

type messageHandler struct {
	sugar       *zap.SugaredLogger
	redisClient *redis.Client
}

func newMessageHandler(sugar *zap.SugaredLogger, redisClient *redis.Client) *messageHandler {
	return &messageHandler{sugar, redisClient}
}

func (hdl *messageHandler) Handle(natsMsg *nats.Msg) {
	natsMsgStr := string(natsMsg.Data)
	m, err := msg.FromJSON(natsMsg.Data)
	if err != nil {
		hdl.sugar.Errorw("failed to decode message", "message", natsMsgStr, "error", err)
		return
	}
	if err := hdl.redisClient.Set(strconv.FormatInt(m.ID, 10), natsMsgStr, 0).Err(); err != nil {
		hdl.sugar.Errorw("failed to set redis", "error", err)
		return
	}
}

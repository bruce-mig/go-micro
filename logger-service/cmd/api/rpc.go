package main

import (
	"context"
	"log"
	"time"

	"github.com/bruce-mig/go-micro/log-service/data"
)

type RPCServer struct{}

// RPCPayload is the type for data we receive from RPC
type RPCPayload struct {
	Name string
	Data string
}

// Loginfo writes our payload to mongo
func (r *RPCServer) LogInfo(payload RPCPayload, resp *string) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), data.LogEntry{
		Name:      payload.Name,
		Data:      payload.Data,
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Println("error writing to mongo:", err)
		return err
	}

	*resp = "Processed payload via RPC:" + payload.Name
	return nil
}

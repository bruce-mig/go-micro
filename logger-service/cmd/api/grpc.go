package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/bruce-mig/go-micro/log-service/data"
	"github.com/bruce-mig/go-micro/log-service/logs"
	"google.golang.org/grpc"
)

// LogServer is type used for writing to the log via gRPC. Note that we embed the
// data.Models type, so we have access to Mongo.
type LogServer struct {
	logs.UnimplementedLogServiceServer
	Models data.Models
}

// WriteLog writes the log after receiving a call from a gRPC client. This function
// must exist, and is defined in logs/logs.proto, in the "service LogService" bit
// at the end of the file.
func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	// write the log
	logEntry := data.LogEntry{
		Name: input.Name,
		Data: input.Data,
	}

	err := l.Models.LogEntry.Insert(logEntry)
	if err != nil {
		res := &logs.LogResponse{Result: "failed"}
		return res, err
	}

	// return response
	res := &logs.LogResponse{Result: "logged via GRPC!"}
	return res, nil
}

// gRPCListen starts the gRPC server
func (app *Config) gRPCListen() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRPCPort))
	if err != nil {
		log.Fatalf("Failed to listen for gRPC: %v", gRPCPort)
	}

	s := grpc.NewServer()
	srv := &LogServer{
		Models: app.Models,
	}

	// register the service, handing it models (so we can write to the database)
	logs.RegisterLogServiceServer(s, srv)

	log.Printf("gRPC Server started on port %s", gRPCPort)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}

}

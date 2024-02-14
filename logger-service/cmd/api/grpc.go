package main

import (
	"context"
	"fmt"
	"github.com/f1zm0n/logger-service/data"
	"github.com/f1zm0n/logger-service/logs"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

type LogServer struct {
	*logs.UnimplementedLogServiceServer
	Models data.Models
}

func (l *LogServer) WriteLog(ctx context.Context, req *logs.LogRequest) (*logs.LogResponse, error) {
	input := req.GetLogEntry()

	logsy := data.LogEntry{
		Name:      input.Name,
		Data:      input.Data,
		CreatedAt: time.Now(),
	}

	err := l.Models.LogEntry.Insert(logsy)
	if err != nil {
		return &logs.LogResponse{Result: "failed mongo"}, err
	}

	return &logs.LogResponse{Result: "logged"}, nil
}

func (c *Config) listenGRPC() error {
	log.Println("starting grpc server on port ", gRpcPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", gRpcPort))
	if err != nil {
		log.Panic(err)
	}

	s := grpc.NewServer()

	logs.RegisterLogServiceServer(s, &LogServer{Models: c.Models})

	if err = s.Serve(lis); err != nil {
		log.Panic(err)
	}

	log.Println("grpc server started on port ", gRpcPort)
	return nil
}

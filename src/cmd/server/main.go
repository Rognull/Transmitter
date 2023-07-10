package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"os"
	"os/signal"
	"syscall"
	"transmitter/internal/server"
)

func main() {
	signals := []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP}
	ctx, cancel := signal.NotifyContext(context.Background(), signals...)
	defer cancel()

	s := grpc.NewServer()
	go server.StartServer(ctx, s)

	<-ctx.Done()
	cancel()
	s.GracefulStop()
	log.Println("Server stopped greacefully")
}

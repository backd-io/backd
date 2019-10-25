package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/fernandezvara/backd/cmd/sessions/store"
	"github.com/fernandezvara/backd/internal/db"
	"github.com/fernandezvara/backd/internal/instrumentation"
	"github.com/fernandezvara/backd/internal/pbsessions"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	// TODO: Add configuration...

	var (
		server   sessionsServer
		mongoURL string
		err      error
	)

	mongoURL = os.Getenv("MONGO_URL")
	if mongoURL == "" {
		fmt.Println("MONGO_URL not found")
		os.Exit(1)
	}

	server.inst, err = instrumentation.New("0.0.0.0:8182", true)
	if err != nil {
		log.Fatal(err)
	}

	server.inst.RegisterGauge("sessions_in_use", "Sessions in use", []string{"hostname"})

	server.store = store.New(server.inst, 10*time.Second)

	err = server.store.Open(true, "1", 8282)
	if err != nil {
		log.Fatal(err)
	}

	server.mongo, err = db.NewMongo(mongoURL)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pbsessions.RegisterSessionsServer(s, server)
	// Register reflection service on gRPC server.
	reflection.Register(s)

	// graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := server.inst.Start(); err != nil {
			server.inst.Error(err.Error())
		}
	}()

	go func() {
		if err := s.Serve(lis); err != nil {
			server.inst.Error(err.Error())
		}
	}()

	<-stop

	server.inst.Info("Shutting down the server.")

	if err := server.inst.Shutdown(); err != nil {
		server.inst.Info(err.Error())
	}

	s.GracefulStop()

}

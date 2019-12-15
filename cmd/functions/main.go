package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/lua"
	"github.com/backd-io/backd/internal/rest"
	"google.golang.org/grpc"
)

// remove, this must be provided by the discovery
const (
	configURLObjects   = "http://objects:8081"
	configURLAuth      = "http://auth:8083"
	configURLAdmin     = "http://admin:8084"
	configURLFunctions = "http://functions:8084"
	configURLSessions  = "sessions:8082"
)

func main() {

	var (
		routes map[string]map[string]rest.APIEndpoint

		server   *rest.REST
		conn     *grpc.ClientConn
		inst     *instrumentation.Instrumentation
		mongo    *db.Mongo
		api      *apiStruct
		b        *backd.Backd
		mongoURL string
		err      error
	)

	mongoURL = os.Getenv("MONGO_URL")
	if mongoURL == "" {
		fmt.Println("MONGO_URL not found")
		os.Exit(1)
	}

	// TODO: REMOVE! AND CONFIGURE PROPERLY
	address := "sessions:8082"

	// Set up a connection to the sessions server.
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongo, err = db.NewMongo(ctx, mongoURL)
	er(err)

	inst, err = instrumentation.New("0.0.0.0:8185", true)
	er(err)

	api = &apiStruct{
		inst:     inst,
		mongo:    mongo,
		sessions: conn,
	}

	b = backd.NewClient(configURLAdmin, configURLObjects, configURLAdmin, configURLFunctions)
	api.lua = lua.New(b).PrepareFunctions()

	routes = map[string]map[string]rest.APIEndpoint{
		"POST": {
			"/functions/:name": {
				Handler: api.postFunction,
				Matcher: []string{"", "^[a-zA-Z0-9-]{2,32}$"},
			},
		},
	}

	server = rest.New("0.0.0.0:8085")
	server.SetupRouter(routes, inst)

	// graceful
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	go func() {
		if err := inst.Start(); err != nil {
			inst.Error(err.Error())
		}
	}()

	go func() {
		if err = server.Start(); err != nil {
			inst.Error(err.Error())
		}
	}()

	<-stop

	inst.Info("Shutting down the server.")

	if err = inst.Shutdown(); err != nil {
		inst.Info(err.Error())
	}

	if err = server.Shutdown(); err != nil {
		inst.Info(err.Error())
	}

}

func er(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

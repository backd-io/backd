package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/rest"
	"google.golang.org/grpc"
)

func main() {

	var (
		routes map[string]map[string]rest.APIEndpoint

		server *rest.REST
		conn   *grpc.ClientConn
		inst   *instrumentation.Instrumentation
		mongo  *db.Mongo
		api    *apiStruct

		err error
	)

	// TODO: REMOVE! AND CONFIGURE PROPERLY
	address := "localhost:8082"

	// Set up a connection to the sessions server.
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	mongo, err = db.NewMongo("mongodb://localhost:27017")
	er(err)

	inst, err = instrumentation.New("0.0.0.0:8181", true)
	er(err)

	api = &apiStruct{
		inst:     inst,
		mongo:    mongo,
		sessions: conn,
	}

	routes = map[string]map[string]rest.APIEndpoint{
		"GET": {
			"/:collection/:id": {
				Handler: api.getDataID,
				Matcher: []string{"^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"},
			},
			"/:collection": {
				Handler: api.getData,
				Matcher: []string{"^[a-zA-Z0-9-]{1,32}$"},
			},
		},
		"POST": {
			"/:collection": {
				Handler: api.postData,
				Matcher: []string{"^[a-zA-Z0-9-]{1,32}$"},
			},
		},
		"PUT": {
			"/:collection/:id": {
				Handler: api.putDataID,
				Matcher: []string{"^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"},
			},
		},
		"DELETE": {
			"/:collection/:id": {
				Handler: api.deleteDataID,
				Matcher: []string{"^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"},
			},
		},
	}

	server = rest.New("0.0.0.0:8081")
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

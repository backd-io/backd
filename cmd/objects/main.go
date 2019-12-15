package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/rest"
	"google.golang.org/grpc"
)

func main() {

	var (
		routes map[string]map[string]rest.APIEndpoint

		server   *rest.REST
		conn     *grpc.ClientConn
		inst     *instrumentation.Instrumentation
		mongo    *db.Mongo
		api      *apiStruct
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

	inst, err = instrumentation.New("0.0.0.0:8181", true)
	er(err)

	api = &apiStruct{
		inst:     inst,
		mongo:    mongo,
		sessions: conn,
	}

	routes = map[string]map[string]rest.APIEndpoint{
		"GET": {
			"/objects/:collection/:id": {
				Handler: api.getObjectID,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"},
			},
			"/objects/:collection": {
				Handler: api.getObject,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$"},
			},
			"/objects/:collection/:id/:relation/:direction": {
				Handler: api.getObjectIDRelations,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$", "^[a-zA-Z0-9-]{1,32}$", "^(in|out)$"},
			},
			"/related/:collection/:id/:direction": {
				Handler: api.getRelations,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$", "^(in|out)$"},
			},
			"/relations/:id": {
				Handler: api.getRelationID,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$"},
			},
		},
		"POST": {
			"/objects/:collection": {
				Handler: api.postObject,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$"},
			},
			"/relations": {
				Handler: api.postRelation,
				Matcher: []string{""},
			},
		},
		"PUT": {
			"/objects/:collection/:id": {
				Handler: api.putObjectID,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"},
			},
		},
		"DELETE": {
			"/objects/:collection/:id": {
				Handler: api.deleteObjectID,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"},
			},
			"/relations/:id": {
				Handler: api.deleteRelationID,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$"},
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

package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/backd-io/backd/internal/db"

	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/rest"
)

func main() {

	var (
		routes   map[string]map[string]rest.APIHandler
		matchers map[string]map[string]rest.APIMatcher

		server *rest.REST
		inst   *instrumentation.Instrumentation
		mongo  *db.Mongo
		api    *apiStruct

		err error
	)

	mongo, err = db.NewMongo("mongodb://localhost:27017")
	er(err)

	inst, err = instrumentation.New("0.0.0.0:8183", true)
	er(err)

	api = &apiStruct{
		inst:  inst,
		mongo: mongo,
	}

	routes = map[string]map[string]rest.APIHandler{
		"GET": {
			"/session": api.getSession,
		},
		"POST": {
			"/session": api.postSession,
		},
		"DELETE": {
			"/session": api.deleteSession,
		},
	}

	matchers = map[string]map[string]rest.APIMatcher{
		"GET": {
			"/session": []string{},
		},
		"POST": {
			"/session": []string{},
		},
		"DELETE": {
			"/session": []string{},
		},
	}

	server = rest.New("0.0.0.0:8083")
	server.SetupRouter(routes, matchers, inst)

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

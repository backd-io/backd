package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/backd-io/backd/backd"
	"github.com/backd-io/backd/pkg/lua"

	"github.com/backd-io/backd/internal/db"
	"github.com/backd-io/backd/internal/instrumentation"
	"github.com/backd-io/backd/internal/rest"
	"google.golang.org/grpc"
)

// remove, this must be provided by the discovery
const (
	configURLObjects  = "http://localhost:8081"
	configURLAuth     = "http://localhost:8083"
	configURLAdmin    = "http://localhost:8084"
	configURLSessions = "localhost:8082"
)

func main() {

	var (
		routes map[string]map[string]rest.APIEndpoint

		server *rest.REST
		conn   *grpc.ClientConn
		inst   *instrumentation.Instrumentation
		mongo  *db.Mongo
		api    *apiStruct
		b      *backd.Backd

		err error
	)

	// TODO: REMOVE! AND CONFIGURE PROPERLY
	// address := "localhost:8082"

	// Set up a connection to the sessions server.
	conn, err = grpc.Dial(configURLSessions, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	mongo, err = db.NewMongo("mongodb://localhost:27017")
	er(err)

	inst, err = instrumentation.New("0.0.0.0:8185", true)
	er(err)

	api = &apiStruct{
		inst:     inst,
		mongo:    mongo,
		sessions: conn,
	}

	b = backd.NewClient(configURLAdmin, configURLObjects, configURLAdmin)
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

// import (
// 	"log"
// 	"net"
// 	"os"
// 	"os/signal"

// 	"github.com/backd-io/backd/backd"
// 	"github.com/backd-io/backd/pkg/lua"
// 	"github.com/backd-io/backd/internal/db"
// 	"github.com/backd-io/backd/internal/instrumentation"
// 	"github.com/backd-io/backd/internal/pbfunctions"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )

// // func main() {
// // 	var (
// // 		server functionsServer
// // 		b      *backd.Backd
// // 		err    error
// // 	)

// // 	// initialize instrumentation
// // 	server.inst, err = instrumentation.New("0.0.0.0:8185", true)
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	// initialize mongo connection
// // 	server.mongo, err = db.NewMongo("mongodb://localhost:27017")
// // 	if err != nil {
// // 		log.Fatal(err)
// // 	}

// // 	// fu := structs.Function{
// // 	// 	ID:   db.NewXID().String(),
// // 	// 	Name: "test",
// // 	// 	API:  true,
// // 	// 	Code: "function code(input)\ninput.nuevo = 'nuevo'\nreturn input\nend",
// // 	// }

// // 	// fu.SetCreate("backd", "asdfasdfasdfasdfasdf")

// // 	// err = server.mongo.Insert("_backd", constants.ColFunctions, fu)
// // 	// if err != nil {
// // 	// 	panic(err)
// // 	// }

// // 	// initialize lang & backd
// // 	b = backd.NewClient(configURLAuth, configURLObjects, configURLAdmin)
// // 	server.lang = lua.New(b).PrepareFunctions()

// // 	lis, err := net.Listen("tcp", ":8085")
// // 	if err != nil {
// // 		log.Fatalf("failed to listen: %v", err)
// // 	}

// // 	s := grpc.NewServer()
// // 	pbfunctions.RegisterFunctionsServer(s, server)
// // 	// Register reflection service on gRPC server.
// // 	reflection.Register(s)

// // 	// graceful
// // 	stop := make(chan os.Signal, 1)
// // 	signal.Notify(stop, os.Interrupt)

// // 	go func() {
// // 		if err := server.inst.Start(); err != nil {
// // 			server.inst.Error(err.Error())
// // 		}
// // 	}()

// // 	go func() {
// // 		if err := s.Serve(lis); err != nil {
// // 			server.inst.Error(err.Error())
// // 		}
// // 	}()

// // 	<-stop

// // 	server.inst.Info("Shutting down the server.")

// // 	if err := server.inst.Shutdown(); err != nil {
// // 		server.inst.Info(err.Error())
// // 	}

// // 	s.GracefulStop()
// // }

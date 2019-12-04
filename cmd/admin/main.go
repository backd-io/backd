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
		routes   map[string]map[string]rest.APIEndpoint
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

	inst, err = instrumentation.New("0.0.0.0:8184", true)
	er(err)

	api = &apiStruct{
		inst:     inst,
		mongo:    mongo,
		sessions: conn,
	}

	err = api.isBootstrapped(ctx)
	er(err)

	// "^[a-zA-Z0-9-]{1,32}$", "^[a-zA-Z0-9]{20}$"

	routes = map[string]map[string]rest.APIEndpoint{
		"GET": {
			"/applications": {
				Handler: api.getApplications,
				Matcher: []string{""},
			},
			"/applications/:id": {
				Handler: api.getApplicationByID,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$"},
			},
			"/applications/:id/functions": {
				Handler: api.getFunctions,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$", ""},
			},
			"/applications/:id/functions/:name": {
				Handler: api.getFunctionByID,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$", "", "^[a-zA-Z0-9-]{2,32}$"},
			},
			"/domains": {
				Handler: api.getDomains,
				Matcher: []string{""},
			},
			"/domains/:domain": {
				Handler: api.getDomainByID,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$"},
			},
			"/domains/:domain/users": {
				Handler: api.getUsers,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", ""},
			},
			"/domains/:domain/users/:id": {
				Handler: api.getUserByID,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$"},
			},
			"/domains/:domain/users/:id/groups": {
				Handler: api.getUserGroups,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$", ""},
			},
			"/domains/:domain/groups": {
				Handler: api.getGroups,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", ""},
			},
			"/domains/:domain/groups/:id": {
				Handler: api.getGroupByID,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$"},
			},
			"/domains/:domain/groups/:id/members": {
				Handler: api.getMembers,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$", ""},
			},
		},
		"POST": {
			"/applications": {
				Handler: api.postApplication,
				Matcher: []string{""},
			},
			"/applications/:id/rbac": {
				Handler: api.rbacApplications,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$", ""},
			},
			"/applications/:id/functions": {
				Handler: api.postFunction,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$", ""},
			},
			"/bootstrap": {
				Handler: api.postBootstrap,
				Matcher: []string{""},
			},
			"/domains": {
				Handler: api.postDomain,
				Matcher: []string{""},
			},
			"/domains/:domain/users": {
				Handler: api.postUser,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", ""},
			},
			"/domains/:domain/groups": {
				Handler: api.postGroup,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", ""},
			},
			"/domains/:domain/rbac": {
				Handler: api.rbacDomains,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", ""},
			},
		},
		"PUT": {
			"/applications/:id": {
				Handler: api.putApplication,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$"},
			},
			"/applications/:id/functions/:name": {
				Handler: api.putFunction,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$", "", "^[a-zA-Z0-9-]{2,32}$"},
			},
			"/domains/:domain": {
				Handler: api.putDomain,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$"},
			},
			"/domains/:domain/users/:id": {
				Handler: api.putUser,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$"},
			},
			"/domains/:domain/groups/:id": {
				Handler: api.putGroup,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$"},
			},
			"/domains/:domain/groups/:id/members/:user_id": {
				Handler: api.putMembership,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$", "", "^[a-zA-Z0-9]{20}$"},
			},
		},
		"DELETE": {
			"/applications/:id": {
				Handler: api.deleteApplication,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$"},
			},
			"/applications/:id/functions/:name": {
				Handler: api.deleteFunction,
				Matcher: []string{"", "^[a-zA-Z0-9]{20}$", "", "^[a-zA-Z0-9-]{1,32}$"},
			},
			"/domains/:domain": {
				Handler: api.deleteDomain,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$"},
			},
			"/domains/:domain/users/:id": {
				Handler: api.deleteUser,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$"},
			},
			"/domains/:domain/groups/:id": {
				Handler: api.deleteGroup,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$"},
			},
			"/domains/:domain/groups/:id/members/:user_id": {
				Handler: api.deleteMembership,
				Matcher: []string{"", "^[a-zA-Z0-9-]{1,32}$", "", "^[a-zA-Z0-9]{20}$", "", "^[a-zA-Z0-9]{20}$"},
			},
		},
	}

	server = rest.New("0.0.0.0:8084")
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

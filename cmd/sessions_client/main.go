package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/backd-io/backd/internal/sessionspb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8082"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := sessionspb.NewSessionsClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 12*time.Minute)
	defer cancel()

	past := time.Now()
	for a := 1; a <= 10000; a++ {
		r, err := c.CreateSession(ctx, &sessionspb.CreateSessionRequest{
			UserId:          "userfake",
			DomainId:        "domainfake",
			DurationSeconds: 300,
		})

		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}

		if a%1000 == 0 {
			log.Printf("[%d][%v]session: %+v", a, err, r)

			rr, err := c.GetSession(ctx, &sessionspb.GetSessionRequest{
				UserId:   "userfake",
				DomainId: "domainfake",
				Id:       r.GetId(),
			})

			log.Printf("[%d][%v]get:session: %+v", a, err, rr)

			rrr, err := c.DeleteSession(ctx, &sessionspb.GetSessionRequest{
				UserId:   "userfake",
				DomainId: "domainfake",
				Id:       r.GetId(),
			})

			log.Printf("[%d][%v]delete:session: %+v", a, err, rrr)

		}
	}

	fmt.Println("Time:", time.Since(past).String())
}

package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"github.com/willschroeder/fingerprint/pb"
)

const (
	address     = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFingerprintServiceClient(conn)

	// Contact the server and print out its response.
	email := "test@test.com"
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.CreateUser(ctx, &pb.CreateUserRequest{Email: email})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetUser().Uuid)
}
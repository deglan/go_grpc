package main

import (
	"context"
	"fmt"
	"go_grpc/greet/greetpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	fmt.Println("Hello World")

	conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)

	//fmt.Printf("Created client: %v", c)
	doUnary(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	request := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	response, err := c.Greet(context.Background(), request)

	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", response)
}

package main

import (
	"context"
	"fmt"
	"go_grpc/greet/greetpb"
	"log"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	greetpb.GreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Received Greet RPC: %v", req)
	firstName := req.GetGreeting().GetFirstName()
	lastName := req.GetGreeting().GetLastName()

	result := "Hello " + firstName + " " + lastName

	return &greetpb.GreetResponse{
		Message: result,
	}, nil
}

func main() {
	fmt.Println("Hello World")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

	fmt.Println("Server started at port 50051")
}

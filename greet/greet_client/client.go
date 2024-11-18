package main

import (
	"context"
	"fmt"
	"go_grpc/greet/greetpb"
	"io"
	"log"
	"time"

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

	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)

	doBiDiStreaming(c)

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

func doServerStreaming(c greetpb.GreetServiceClient) {
	request := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "John",
			LastName:  "Doe",
		},
	}
	stream, err := c.GreetManyTimes(context.Background(), request)

	if err != nil {
		log.Fatalf("Error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", response.GetResult())
	}
}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Do CLient Streaming")

	stream, err := c.LongGreet(context.Background())

	req := []*greetpb.LongGreetRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
				LastName:  "Doe",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephan",
				LastName:  "Doe",
			},
		},
	}

	if err != nil {
		log.Fatalf("error while calling LongGreet RPC: %v", err)
	}

	for _, req := range req {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("error while receiving response: %v", err)
	}

	fmt.Printf("LongGreet Response: %v\n", res)
}

func doBiDiStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Do BiDi Streaming")

	req := []*greetpb.GreetEveryoneRequest{
		{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Jane",
				LastName:  "Doe",
			},
		},
		{
			Greeting: &greetpb.Greeting{
				FirstName: "Stephan",
				LastName:  "Doe",
			},
		},
	}

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while calling GreetEveryone RPC: %v", err)
		return
	}

	waitChannel := make(chan struct{})

	// send data to server
	go func() {
		for _, req := range req {
			fmt.Printf("Sending req: %v\n", req)
			stream.Send(req)
			time.Sleep(1 * time.Second)
		}
		stream.CloseSend()
	}()

	// recive data from client
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while receiving response: %v", err)
				break
			}
			fmt.Printf("Recived response from GreetEveryone: %v\n", res)
		}
		close(waitChannel)
	}()

	// block until everything is done
	<-waitChannel
}

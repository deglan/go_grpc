package main

import (
	"context"
	"fmt"
	"go_grpc/calculator/calculatorpb"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func main() {
	fmt.Println("Hello World")

	conn, err := grpc.NewClient("0.0.0.0:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}

	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	//fmt.Printf("Created client: %v", c)
	doUnary(c)

	doErrorUnary(c)

}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	request := &calculatorpb.SumRequest{
		Num1: 3,
		Num2: 10,
	}
	response, err := c.Sum(context.Background(), request)

	if err != nil {
		log.Fatalf("Error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", response)
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	doErrorCall(c, 10)
	doErrorCall(c, -10)

}

func doErrorCall(c calculatorpb.CalculatorServiceClient, number int32) {
	response, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: number})

	if err != nil {
		responseError, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Println("Error message from server: ", responseError.Message())
			fmt.Println("Error code from server: ", responseError.Code())
			if responseError.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent negative number!")
				return
			}
		} else {
			log.Fatalf("Big Error while calling Square Root RPC: %v", err)
		}
	}
	log.Printf("Response from Square Root: %v : %v", number, response.GetNumberRoot())
}

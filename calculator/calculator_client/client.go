package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ctc316/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := calculatorpb.NewCalculatorServiceClient(conn)

	doUnary(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.CalculatorRequest{
		Input: &calculatorpb.Input{
			Num1: 111.111,
			Num2: 222.222,
		},
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

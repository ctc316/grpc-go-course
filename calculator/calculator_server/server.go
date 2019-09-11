package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/ctc316/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	num1 := req.GetInput().GetNum1()
	num2 := req.GetInput().GetNum2()
	result := num1 + num2
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}

func main() {
	fmt.Println("Calculator Service is ON")

	lis, err := net.Listen("tcp", "0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve, %v", err)
	}
}

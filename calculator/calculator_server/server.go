package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/ctc316/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	num1 := req.GetInput2Num().GetNum1()
	num2 := req.GetInput2Num().GetNum2()
	result := num1 + num2
	res := &calculatorpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeDecompose(req *calculatorpb.PrimeDecomposeRequest, stream calculatorpb.CalculatorService_PrimeDecomposeServer) error {
	fmt.Printf("PrimeDecompose function was invoked with %v\n", stream)
	N := req.GetNum()
	var k int64 = 2
	for N > 1 {
		if N%k == 0 {
			res := &calculatorpb.PrimeDecomposeReponse{
				Result: k,
			}
			stream.Send(res)
			time.Sleep(500 * time.Millisecond)

			N = N / k
		} else {
			k++
		}

	}
	return nil
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

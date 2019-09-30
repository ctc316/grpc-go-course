package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/ctc316/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc/reflection"

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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Printf("ComputeAverage function was invoked with %v\n", stream)
	var result float32 = 0.0
	n := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: float32(result / float32(n)),
			})
		}
		if err != nil {
			log.Fatalf("Error while reading client stream")
		}
		result += req.GetNum()
		n++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked with a streaming request")

	maxNum := -math.MaxFloat64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}
		maxNum = math.Max(maxNum, float64(req.GetNum()))
		sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
			Result: float32(maxNum),
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SqurareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Calculator Service is ON")

	lis, err := net.Listen("tcp", "0.0.0:50052")
	if err != nil {
		log.Fatalf("Failed to listen %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC sever
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve, %v", err)
	}
}

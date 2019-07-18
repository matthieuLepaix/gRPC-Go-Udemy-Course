package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/matthieulepaix/gRPC-Udemy-Course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, in *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	sum := in.GetPart1() + in.GetPart2()
	res := &calculatorpb.SumResponse{
		Result: sum,
	}
	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	var k int32 = 2
	N := req.GetNumber()
	for N > 1 {
		mod := N % k
		if mod == 0 {
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				DecompositionPart: k,
			}
			stream.Send(res)
			N = N / k
		} else {
			k++
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {

	sum := int64(0)
	counter := int64(0)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			average := float64(sum) / float64(counter)

			return stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: average,
			})
		}
		if err != nil {
			log.Fatalf("Error while reading stream ComputerAverage: %v", err)
			return err
		}
		sum += msg.GetNumber()
		counter++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	currentMax := int32(math.MinInt32)
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while receiving message FindMaximum: %v", err)
			return err
		}
		result := msg.GetNumber()
		if result > currentMax {
			currentMax = result
			stream.Send(&calculatorpb.FindMaximumResponse{
				Maximum: currentMax,
			})
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
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
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

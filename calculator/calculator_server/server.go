package main

import (
	"context"
	"io"
	"log"
	"net"

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
		}
		sum += msg.GetNumber()
		counter++
	}
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

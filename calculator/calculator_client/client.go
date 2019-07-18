package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/matthieulepaix/gRPC-Udemy-Course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	doUnary(c)

	doServerStreaming(c)

	doClientStreaming(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.SumRequest{
		Part1: 3,
		Part2: 10,
	}
	res, err := c.Sum(context.Background(), req)

	if err != nil {
		log.Fatalf("Failed while calling Sum RPC: %v", err)
	}
	fmt.Printf("Result of %v + %v: %v", req.GetPart1(), req.GetPart2(), res)
}

func doServerStreaming(c calculatorpb.CalculatorServiceClient) {
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	res, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed while calling stream PrimeNumberDecomposition: %v", err)
	}
	for {
		msg, err := res.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed while reading stream Response from PrimeNumberDecomposition: %v", err)
		}
		log.Printf("Prime Decomposition of 120: %v", msg)
	}
}


func doClientStreaming(c calculatorpb.CalculatorServiceClient) {

	numbers := []int64{
		1,
		2,
		3,
		4,
	}

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage: %v", err)
	}

	log.Printf("Average of: ")
	for _, number := range numbers {
		
		log.Printf("%v ", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
	}
	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receive the response ComputeAverage: %v", err)
	}
	log.Printf("Result: %v", response.Result)
}

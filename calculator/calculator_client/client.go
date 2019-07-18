package main

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	//doSum(c) //Unary

	//doPrimeNumber(c) // Server Streaming

	//doAverage(c) // Client streaming

	//doFindMaximum(c) // Bidirecitonal streaming

	doSquareRoot(c) //Error handling
}

func doSum(c calculatorpb.CalculatorServiceClient) {
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

func doPrimeNumber(c calculatorpb.CalculatorServiceClient) {
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

func doAverage(c calculatorpb.CalculatorServiceClient) {

	numbers := []int64{1, 2, 3, 4}

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

func doFindMaximum(c calculatorpb.CalculatorServiceClient) {
	numbers := []int32{1, 5, 3, 6, 2, 20}

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while calling the client FindMaximum: %v", err)
	}

	waitc := make(chan struct{})

	// send data
	go func() {
		for _, number := range numbers {
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// Receive data
	go func() {
		for {
			res, recErr := stream.Recv()

			if recErr == io.EOF {
				break
			}

			if recErr != nil {
				log.Fatalf("Error while Receiving FindMaximum: %v", recErr)
			}

			log.Printf("New Maximum: %v", res.GetMaximum())
		}
		close(waitc)
	}()

	<-waitc
}

func doSquareRoot(c calculatorpb.CalculatorServiceClient) {
	doSquareRootCall(c, 10)
	doSquareRootCall(c, -10)
}

func doSquareRootCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{
		Number: n,
	})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// user error
			log.Println(respErr.Message())
			log.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				log.Println("We probably sent a negative number")
				return
			}
		}else {
			log.Fatalf("Big error calling SquareRoot: %v", err)
		}
	}
	fmt.Printf("Result of sqaure root of %v: %v\n", n, res.GetNumberRoot())
}

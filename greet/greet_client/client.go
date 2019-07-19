package main

import (
	"context"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/matthieulepaix/gRPC-Udemy-Course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	//doUnary(c)

	//doServerStreaming(c)

	//doBiDirectionalStreaming(c)

	doUnaryWithDeadLine(c, 1*time.Second)
	doUnaryWithDeadLine(c, 5*time.Second)
}

func doUnary(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Matthieu",
			LastName:  "Lepaix",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC : %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreaming(c greetpb.GreetServiceClient) {
	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Matthieu",
			LastName:  "Lepaix",
		},
	}
	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error streaming while calling GreetManyTimes RPC : %v", err)
	}

	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}

		log.Printf("Response from GreetManyTimes: %v", msg)
	}

}

func doBiDirectionalStreaming(c greetpb.GreetServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Matthieu",
				LastName:  "Lepaix",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Yoann",
				LastName:  "Pencalet",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Rami",
				LastName:  "Ben Hamed",
			},
		},
	}

	waitc := make(chan struct{})

	// Send messages
	go func() {
		for _, req := range requests {
			SendErr := stream.Send(req)
			if SendErr != nil {
				log.Fatalf("Error while sending BiDirectionalStreaming: %v", SendErr)
				return
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// Receive messages
	go func() {
		for {
			res, RecErr := stream.Recv()
			if RecErr == io.EOF {
				break
			}
			if RecErr != nil {
				log.Fatalf("Error while receiving BiDirectionalStreaming: %v", RecErr)
				close(waitc)
			}
			log.Printf("Received: %v", res.GetResult())
		}
		close(waitc)
	}()

	<-waitc
}

func doUnaryWithDeadLine(c greetpb.GreetServiceClient, t time.Duration) {
	req := &greetpb.GreetWithDeadLineRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Matthieu",
			LastName:  "Lepaix",
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	res, err := c.GreetWithDeadLine(ctx, req)
	if err != nil {

		statusErr, ok := status.FromError(err)
		if ok {
			if statusErr.Code() == codes.DeadlineExceeded {
				log.Printf("Timeout was hit! Deadline was exceeded")
			} else {
				log.Printf("Unexpected error: %v", statusErr)
			}
			return
		} else {
			log.Fatalf("error while calling Greet RPC : %v", err)
		}
	}
	log.Printf("Response from Greet: %v", res.Result)
}

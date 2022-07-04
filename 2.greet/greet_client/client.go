package main

import (
	"context"
	pb "github.com/manerajona/grpc/2.greet/greetpb"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/credentials"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	Insecure = false
	CertFile = ".ssl/ca.crt" // Certificate Authority Trust certificate
)

func doUnary(c pb.GreetServiceClient, req *pb.GreetRequest) {
	log.Printf("Starting Unary...")

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Server: %v", err)
	}
	log.Printf("Response from Server: %v", res.Result)
}

func doServerStreaming(c pb.GreetServiceClient, req *pb.GreetManyTimesRequest) {
	log.Printf("Starting Server Streaming...")

	stream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				break
			}
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from Server: %v", res.Result)
	}

}

func doClientStreaming(c pb.GreetServiceClient, names ...string) {
	log.Printf("Starting Client Streaming...")

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, name := range names {
		log.Printf("Sending: %v\n", name)
		stream.Send(&pb.LongGreetRequest{
			Greeting: &pb.Greeting{
				FirstName: name,
			},
		})
		time.Sleep(500 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while calling Server: %v", err)
	}
	log.Printf("Response from Server: %v", res.Result)

}

func doBiDiStreaming(c pb.GreetServiceClient, names ...string) {
	log.Printf("Starting Bi-directional Streaming...")

	// we create a stream by invoking the client
	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
		return
	}

	// we send a bunch of messages to the client (go routine)
	go func() {
		// function to send a bunch of messages
		for _, name := range names {
			log.Printf("Sending: %v\n", name)
			stream.Send(&pb.GreetEveryoneRequest{
				Greeting: &pb.Greeting{
					FirstName: name,
				},
			})
			time.Sleep(500 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	waitChan := make(chan struct{})

	go func() {
		// function to receive a bunch of messages
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					// we've reached the end of the stream
					break
				}
				log.Fatalf("error while reading stream: %v", err)
			}
			log.Printf("Received: %v\n", res.GetResult())
		}
		close(waitChan)
	}()

	// block until everything is done
	<-waitChan
}

func doUnaryWithDeadline(c pb.GreetServiceClient, req *pb.GreetWithDeadlineRequest, timeout time.Duration) {
	log.Printf("Starting UnaryWithDeadline RPC...")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	res, err := c.GreetWithDeadline(ctx, req)
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			log.Printf("%v: %v\n", respErr.Code(), respErr.Message())
		} else {
			log.Fatalf("Unknown error: %v", err)
		}
		return
	}
	log.Printf("Response from Server: %v", res.Result)
}

func main() {

	var opts grpc.DialOption
	if Insecure {
		opts = grpc.WithTransportCredentials(insecure.NewCredentials())
	} else {
		cred, sslErr := credentials.NewClientTLSFromFile(CertFile, "")
		if sslErr != nil {
			log.Fatalf("Error while loading CA trust certificate: %v", sslErr)
		}
		opts = grpc.WithTransportCredentials(cred)
	}

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("%v is unrechable", err)
	}
	defer cc.Close()

	c := pb.NewGreetServiceClient(cc)

	doUnary(c, &pb.GreetRequest{Greeting: &pb.Greeting{FirstName: "Jona", LastName: "Manera"}})

	//doServerStreaming(c, &pb.GreetManyTimesRequest{Greeting: &pb.Greeting{FirstName: "Jona", LastName: "Manera"}})

	//doClientStreaming(c, "Jona", "Stephane", "Joe", "Jane")

	//doBiDiStreaming(c, "Jona", "Stephane", "Joe", "Jane")

	//request := &pb.GreetWithDeadlineRequest{Greeting: &pb.Greeting{FirstName: "Jona", LastName: "Manera"}}
	//doUnaryWithDeadline(c, request, 5*time.Second) // should complete
	//doUnaryWithDeadline(c, request, time.Second)   // should time out
}

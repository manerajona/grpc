package main

import (
	"context"
	pb "github.com/manerajona/grpc/3.calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/status"
)

func doUnary(c pb.CalculatorServiceClient, arg0 rune, arg1 rune) {
	log.Println("Sending Sum Unary RPC...")
	req := &pb.SumRequest{
		FirstNumber:  arg0,
		SecondNumber: arg1,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumResult)
}

func doServerStreaming(c pb.CalculatorServiceClient, number int64) {
	log.Println("Starting PrimeDecomposition Server Streaming RPC...")

	req := &pb.PrimeNumberDecompositionRequest{
		Number: number,
	}
	stream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeDecomposition RPC: %v", err)
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
		log.Printf("Response from Server: %v", res.PrimeFactor)
	}
}

func doClientStreaming(c pb.CalculatorServiceClient, numbers ...int32) {
	log.Println("Starting ComputeAverage Client Streaming RPC...")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream: %v", err)
	}

	for _, number := range numbers {
		log.Printf("Sending number: %v\n", number)
		stream.Send(&pb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(300 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Avg: %v", res.Average)
}

func doBiDiStreaming(c pb.CalculatorServiceClient, numbers ...int32) {
	log.Println("Starting FindMaximum Bi-directional Streaming RPC...")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while opening stream and calling FindMaximum: %v", err)
	}

	// send go routine
	go func() {
		for _, number := range numbers {
			log.Printf("Sending number: %v\n", number)
			stream.Send(&pb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(300 * time.Millisecond)
		}
		stream.CloseSend()
	}()

	waitChan := make(chan struct{})

	// receive go routine
	go func() {
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					// we've reached the end of the stream
					break
				}
				log.Fatalf("error while reading stream: %v", err)
			}
			log.Printf("Received: %v\n", res.GetMaximum())
		}
		close(waitChan)
	}()
	<-waitChan
}

func doErrorCall(c pb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &pb.SquareRootRequest{Number: n})
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
	log.Printf("Square root of %v eq %v\n", n, res.GetNumberRoot())
}

func main() {

	cc, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("%v is unrechable", err)
	}
	defer cc.Close()

	c := pb.NewCalculatorServiceClient(cc)

	//doUnary(c, 2, 10)

	//doServerStreaming(c, 12390392840)

	//doClientStreaming(c,3, 5, 9, 54, 23)

	//doBiDiStreaming(c, 4, 7, 2, 19, 4, 6, 32)

	doErrorCall(c, 9)

	doErrorCall(c, -10)
}

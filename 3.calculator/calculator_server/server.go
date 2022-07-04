package main

import (
	"context"
	"fmt"
	pb "github.com/manerajona/grpc/3.calculator/calculatorpb"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(_ context.Context, req *pb.SumRequest) (*pb.SumResponse, error) {

	sum := req.FirstNumber + req.SecondNumber

	return &pb.SumResponse{
		SumResult: sum,
	}, nil
}

func (*server) PrimeNumberDecomposition(req *pb.PrimeNumberDecompositionRequest, stream pb.CalculatorService_PrimeNumberDecompositionServer) error {

	divisor := int64(2)

	for n := req.GetNumber(); n > 1; {
		if n%divisor == 0 {
			stream.Send(&pb.PrimeNumberDecompositionResponse{
				PrimeFactor: divisor,
			})
			n /= divisor
			continue
		}
		divisor++
	}
	return nil
}

func (*server) ComputeAverage(stream pb.CalculatorService_ComputeAverageServer) error {

	sum := int32(0)

	for count := 0; ; count++ {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we've reached the end of the stream
				average := float64(sum) / float64(count)
				return stream.SendAndClose(&pb.ComputeAverageResponse{
					Average: average,
				})
			}
			log.Fatalf("error while reading stream: %v", err)
		}
		sum += req.GetNumber()
	}

}

func (*server) FindMaximum(stream pb.CalculatorService_FindMaximumServer) error {

	for maximum := int32(0); ; {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we have finished reading the client stream
				return nil
			}
			log.Fatalf("error while reading stream: %v", err)
			return err
		}

		if number := req.GetNumber(); number > maximum {
			maximum = number
			sendErr := stream.Send(&pb.FindMaximumResponse{
				Maximum: maximum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *pb.SquareRootRequest) (*pb.SquareRootResponse, error) {

	if n := req.GetNumber(); n < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received negative number %v", n),
		)
	}

	sqrt := math.Sqrt(float64(req.GetNumber()))

	return &pb.SquareRootResponse{
		NumberRoot: sqrt,
	}, nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCalculatorServiceServer(s, &server{})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	log.Println("serving at 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

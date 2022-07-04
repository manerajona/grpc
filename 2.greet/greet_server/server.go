package main

import (
	"context"
	"fmt"
	pb "github.com/manerajona/grpc/2.greet/greetpb"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

const (
	StreamingLimit = 10
	TLSNeeded      = true
	CertFile       = ".ssl/server.crt"
	KeyFile        = ".ssl/server.pem"
)

type server struct{}

func (*server) Greet(_ context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {

	return &pb.GreetResponse{
		Result: fmt.Sprintf("Hello %v!", req.GetGreeting().GetFirstName()),
	}, nil
}

func (*server) GreetManyTimes(req *pb.GreetManyTimesRequest, stream pb.GreetService_GreetManyTimesServer) error {

	firstName := req.GetGreeting().GetFirstName()

	for index := 0; index < StreamingLimit; index++ {
		res := &pb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %v! [%v]", firstName, index),
		}
		stream.Send(res)
		time.Sleep(300 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream pb.GreetService_LongGreetServer) error {

	for result := ""; ; {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we have finished reading the client stream
				return stream.SendAndClose(&pb.LongGreetResponse{
					Result: result,
				})
			}
			log.Fatalf("error while reading stream: %v", err)
		}
		result += fmt.Sprintf("Hello %v! ", req.GetGreeting().GetFirstName())
	}
}

func (*server) GreetEveryone(stream pb.GreetService_GreetEveryoneServer) error {

	for {
		req, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				// we have finished reading the client stream
				return nil
			}
			log.Fatalf("error while reading stream: %v", err)
			return err
		}

		sendErr := stream.Send(&pb.GreetEveryoneResponse{
			Result: fmt.Sprintf("Hello %v! ", req.GetGreeting().GetFirstName()),
		})
		if sendErr != nil {
			log.Fatalf("Error while sending data to client: %v", sendErr)
			return sendErr
		}
	}

}

func (*server) GreetWithDeadline(ctx context.Context, req *pb.GreetWithDeadlineRequest) (*pb.GreetWithDeadlineResponse, error) {

	// delay 3 sec
	for sec := 3; sec > 0; sec-- {
		if ctx.Err() == context.DeadlineExceeded {
			// the client canceled the request
			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}
		time.Sleep(time.Second)
	}

	res := &pb.GreetWithDeadlineResponse{
		Result: fmt.Sprintf("Hello %v!", req.GetGreeting().GetFirstName()),
	}
	return res, nil
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	if TLSNeeded {
		cred, sslErr := credentials.NewServerTLSFromFile(CertFile, KeyFile)
		if sslErr != nil {
			log.Fatalf("Failed loading certificates: %v", sslErr)
		}
		opts = append(opts, grpc.Creds(cred))
	}

	s := grpc.NewServer(opts...)
	pb.RegisterGreetServiceServer(s, &server{})

	log.Println("serving at 0.0.0.0:50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

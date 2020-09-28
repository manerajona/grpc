package main

import (
	"context"
	"fmt"
	"net"

	echo "github.com/manerajona/grpc/02-server/echo"
	grpc "google.golang.org/grpc"
)

type EchoServer struct{}

func (e *EchoServer) Echo(ctx context.Context, req *echo.EchoRequest) (*echo.EchoResponse, error) {
	return &echo.EchoResponse{
		Response: "My Echo: " + req.Message,
	}, nil
}

func main() {
	lst, _ := net.Listen("tcp", ":8080")

	s := grpc.NewServer()
	srv := &EchoServer{}
	echo.RegisterEchoServerServer(s, srv)

	fmt.Println("Serving at port 8080")

	if err := s.Serve(lst); err != nil {
		panic(err)
	}
}

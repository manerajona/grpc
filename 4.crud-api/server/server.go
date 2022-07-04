package main

import pb "github.com/manerajona/grpc/4.crud-api/blogpb"

type Server struct {
	pb.BlogServiceServer
}

//go:build !test
// +build !test

package main

import (
	"context"
	pb "github.com/manerajona/grpc/4.crud-api/blogpb"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"log"
	"net"
)

const (
	address  = "0.0.0.0:50051"
	mongoUri = "mongodb://root:root@localhost:27017/"
	dbName   = "blogdb"
)

func connectToMongo() error {
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUri))
	if err != nil {
		return err
	}

	if err := client.Connect(context.Background()); err != nil {
		return err
	}
	log.Printf("Connected to %v\n", mongoUri)

	collection = client.Database(dbName).Collection("blog")
	log.Printf("Database %v created\n", dbName)

	return nil
}

func main() {

	if err := connectToMongo(); err != nil {
		log.Fatal(err)
	}

	// serve
	serv := grpc.NewServer()
	pb.RegisterBlogServiceServer(serv, &Server{})

	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}

	log.Printf("Listening at %s\n", address)
	if err := serv.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

package main

import (
	"context"
	"fmt"
	pb "github.com/manerajona/grpc/4.crud-api/blogpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"log"
	"strconv"
)

const target = "0.0.0.0:50051"

func createBlog(c pb.BlogServiceClient, blog *pb.Blog) (res *pb.BlogId) {
	res, err := c.CreateBlog(context.Background(), blog)
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	return
}

func readOne(c pb.BlogServiceClient, id string) (res *pb.Blog) {
	res, err := c.ReadBlog(context.Background(), &pb.BlogId{Id: id})
	if err != nil {
		log.Fatalf("Error happened while reading: %v\n", err)
	}
	return
}

func updateBlog(c pb.BlogServiceClient, id string, blog *pb.Blog) {
	_, err := c.UpdateBlog(context.Background(), blog)
	if err != nil {
		log.Printf("Error happened while updating: %v\n", err)
	}
}

func deleteBlog(c pb.BlogServiceClient, id string) {
	_, err := c.DeleteBlog(context.Background(), &pb.BlogId{Id: id})
	if err != nil {
		log.Fatalf("Error happened while deleting: %v\n", err)
	}
}

func listBlog(c pb.BlogServiceClient) (res []*pb.Blog) {
	stream, err := c.ListBlogs(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Fatalf("Error while calling ListBlogs: %v\n", err)
	}

	for {
		blog, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Something happened: %v\n", err)
		}
		res = append(res, blog)
	}
	return
}

func main() {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Couldn't connect to %v: %v\n", target, err)
	}
	defer conn.Close()

	c := pb.NewBlogServiceClient(conn)

	// create
	newBlog := &pb.Blog{
		AuthorId: "Original Author",
		Title:    "My First Blog",
		Content:  "Content of the first blog",
	}

	id := createBlog(c, newBlog).Id
	fmt.Println("Blog created with id: " + id)

	// read
	blog := readOne(c, id)
	fmt.Println("Blog found: " + blog.String())

	// update
	toUpdate := &pb.Blog{
		Id:       id,
		AuthorId: "Changed Author",
		Title:    "My First Blog (edited)",
		Content:  "Content of the first blog, with some awesome additions!",
	}

	updateBlog(c, id, toUpdate)
	fmt.Println("Blog with id " + id + " updated")

	// stream
	blogs := listBlog(c)
	fmt.Println("List of Blogs:")
	for i, blog := range blogs {
		fmt.Println("[" + strconv.Itoa(i) + "] " + blog.String())
	}

	// delete
	deleteBlog(c, id)
	fmt.Println("Blog with id " + id + " deleted")

	// NotFound
	readOne(c, id)
}

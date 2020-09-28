package main

import (
	"context"
	"fmt"

	echo "github.com/manerajona/grpc/03-client/echo"
	grpc "google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	conn, _ := grpc.Dial("localhost:8080", grpc.WithInsecure()) // no ssl

	defer conn.Close()
	{
		e := echo.NewEchoServerClient(conn)
		resp, err := e.Echo(ctx, &echo.EchoRequest{
			Message: "Hello World!",
		})
		if err != nil {
			panic(err)
		}

		fmt.Println("Got from server:", resp.Response)
	}
}

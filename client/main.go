package main

import (
	"context"
	"flag"
	"log"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	pb "github.com/okpalaChidiebere/go-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	addr = flag.String("addr", "localhost:9000", "the address to connect to")
	text = flag.String("t", "", "The todo text")
)

func main() {
	
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	//the client connects to the server. It is important to Note that.
	c := pb.NewTodoServiceClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second) //create a new context
	defer cancel()
	r, err := c.CreateTodo(ctx, &pb.CreateTodoRequest{ Text: *text })
	if err != nil {
		log.Fatalf("could not Create Todo: %v", err)
	}
	log.Printf("New Todo: %s", r.GetTodo())

	rtr, err := c.ReadTodos(ctx, &empty.Empty{})
	if err != nil {
		log.Fatalf("could not Create Todo: %v", err)
	}
	log.Printf("Read the todos from server: %s", rtr.String())
}
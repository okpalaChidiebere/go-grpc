package main

import (
	"bufio"
	"context"
	"flag"
	"io"
	"log"
	"os"
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

	// The client is reading a stream of data from the server
	stream, err := c.ReadTodosStream(ctx, &empty.Empty{})
    if err != nil {
        log.Fatalf("%v.Execute(ctx) = %v, %v: ", c, stream, err)
    }
	for {
		item, err := stream.Recv()
		if err == io.EOF {
			log.Printf("Server done!")
			break
		} else if err != nil {
			log.Fatalln("something went wrong with getting items: ", err)
		}
		log.Printf("Received item from server: %s", item.String())
	}

	ch := make(chan *pb.AddTodoPhotoResponse)
	//The Client Sending a stream to the server
	go func (ch chan *pb.AddTodoPhotoResponse)  {
		

		stream, err := c.AddPhoto(context.Background())
		if err != nil {
			log.Fatalf("%v.Execute(ctx) = %v, %v: ", c, stream, err)
		}

		stream.Send(&pb.AddTodoPhotoRequest{ Request: &pb.AddTodoPhotoRequest_Info{ Info: &pb.AddTodoPhotoRequest_PhotoInfo{ TodoId: "SomeRandomTodoId"} }})

		f, err := os.Open("Penguins.jpeg")

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		reader := bufio.NewReader(f)
		buf := make([]byte, 256)

		for {
			n, err := reader.Read(buf)
			if n == 0 {
				if err == nil {
					continue
				}
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}

			// process buf
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}
			stream.Send(&pb.AddTodoPhotoRequest{ Request: &pb.AddTodoPhotoRequest_Data{ Data: buf }})
		}

		res, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Error when closing the stream and receiving the response: %v", err)
		}
		ch <- res
	}(ch)

	mRes := <-ch
	log.Printf("AddPhoto Server Response to client: %s", mRes.String())
}
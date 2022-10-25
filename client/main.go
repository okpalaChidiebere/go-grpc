package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/okpalaChidiebere/go-grpc/data"
	"google.golang.org/grpc/credentials"

	pb "github.com/okpalaChidiebere/go-grpc/pb"
	"google.golang.org/grpc"
)

var (
	addr = flag.String("addr", "0.0.0.0:9000", "the address to connect to") //we listen for all IP address on local machine :). You can target only localhost and any external connection if you want. see https://www.howtogeek.com/225487/what-is-the-difference-between-127.0.0.1-and-0.0.0.0/
	text = flag.String("t", "", "The todo text")
)


func newClientTLSFromFile(certFile, serverNameOverride string) (credentials.TransportCredentials, error) {
	//Load the certificate of the CA who signed the server's certificate 
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	cp := x509.NewCertPool()
	//attempt to verify the authenticity of the certificate we loaded to make sure its the right server
	if !cp.AppendCertsFromPEM(b) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	//create credentials and return it
	return credentials.NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}

func main() {
	
	flag.Parse()

	// Create tls based credential.
	creds, err := newClientTLSFromFile(data.Path("cert/ca_cert.pem"), "x.example.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	/* FYI: You use grpc.WithBlock() when you want to block main thread until connection is made. By default the connection is made in the background thread
	 Eg Using DialContext() and use a context with timeout to avoid wait a lot on the main thread; then you will need grpc.WithBlock() */
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds)) // Set up a connection to the server.
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
package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/okpalaChidiebere/go-grpc/data"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	api "github.com/okpalaChidiebere/go-grpc/api"
	todosservice "github.com/okpalaChidiebere/go-grpc/businessLogic/todos"
	todosrepo "github.com/okpalaChidiebere/go-grpc/dataLayer/todos"
)

var (
	port = flag.Int("port", 9000, "The server port")
)

// type Server struct {
// 	Todo  pb.TodoServiceServer
// }

func loadTLSCredentials(certFile, keyFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	config := tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientAuth:   tls.NoClientCert, //Our server don't care about the client calling its API
    }
	return credentials.NewTLS(&config) , nil
}

func main(){
	flag.Parse()

	//NOTE: If your repo needs any vendor client as params like dynamoDB, 
	//awsSecretManager, etc, you initialize them here and inject them into 
	//the repo
	todosRepo := todosrepo.NewTodosListDataRepo()
	todosService := todosservice.New(todosRepo)

	/*
	NOTE: By default gRPC uses HTTP2 which needs to be secure. If you are working
	in demo you could use insecure credentials by calling insecure.NewCredentials()
	but it is good practice to try to encrypt your data over the wire even when working
	in demo when you can create certificates for free using openssl. 
	*/
	creds, err := loadTLSCredentials(data.Path("cert/server_cert.pem"), data.Path("cert/server_key.pem"))// Create tls based credential.
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}
	grpcServer := grpc.NewServer(grpc.Creds(creds))

	s := api.Servers{
		TodoServer: api.NewTodoServer(todosService),
	}
	s.RegisterAllService(grpcServer)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

//https://www.youtube.com/watch?v=-f4Gbk-U758

//https://www.youtube.com/watch?v=7YgaZIFn7mY
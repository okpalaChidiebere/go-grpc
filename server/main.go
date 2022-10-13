package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	api "github.com/okpalaChidiebere/go-grpc/api"
	todosservice "github.com/okpalaChidiebere/go-grpc/businessLogic/todos"
	todosrepo "github.com/okpalaChidiebere/go-grpc/dataLayer/todos"
	pb "github.com/okpalaChidiebere/go-grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port = flag.Int("port", 9000, "The server port")
)

type Server struct {
	Todo  pb.TodoServiceServer
}

func main(){
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	//NOTE: If your repo needs any vendor client as params like dynamoDB, 
	//awsSecretManager, etc, you initialize them here and inject them into 
	//the repo
	todosRepo := todosrepo.NewTodosListDataRepo()
	todosService := todosservice.New(todosRepo)

	/*
	NOTE: By default gRPC uses HTTP2 which needs to be secure. But 
	for this demo we will not pass any options to the NewServer() 
	method like certificates to encrypt our data over the wire
	
	@see https://letsencrypt.org/ to get a free certificate

	To see how to add the cert to ur server got to /examples/features/encryption
	in zip folder downloaded here https://grpc.io/docs/languages/go/quickstart/#get-the-example-code
	*/
	creds := insecure.NewCredentials()

	grpcServer := grpc.NewServer(grpc.Creds(creds))

	s := Server{
		Todo: api.NewTodoServer(todosService),
	}
	pb.RegisterTodoServiceServer(grpcServer, s.Todo)

	log.Printf("server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
package api

import (
	"google.golang.org/grpc"

	pb "github.com/okpalaChidiebere/go-grpc/pb"
)

type Servers struct {
	TodoServer  pb.TodoServiceServer
}

func (a Servers) RegisterAllService (s *grpc.Server){
	pb.RegisterTodoServiceServer(s, a.TodoServer)
}
package api

import (
	"context"
	"log"

	//NOTE: always good to alias your import for good code readability. It is not needed though because you can use the package 'name' on the imported file :)
	apiadapters "github.com/okpalaChidiebere/go-grpc/api/adapters"
	todos "github.com/okpalaChidiebere/go-grpc/businessLogic/todos"
	pb "github.com/okpalaChidiebere/go-grpc/pb"
)

// TodoServer is used to implement todos.TodoServiceServer interface. It's a MUST
type TodoServer struct {
	pb.UnimplementedTodoServiceServer
	TodoService todos.Service
}

//constructor for our TodoServer
func NewTodoServer(todoService todos.Service) pb.TodoServiceServer {
	return &TodoServer{ 
		TodoService: todoService,
	}
}

// CreateTodo implements todos.TodoServiceServer.
func (s *TodoServer) CreateTodo(ctx context.Context, in *pb.CreateTodoRequest) (*pb.CreateTodoResponse, error) {
	log.Printf("Received todo: %v", in.GetText())
	todo := s.TodoService.Create(in.GetText())
	return &pb.CreateTodoResponse{Todo: apiadapters.TodoToProto(&todo) } , nil
}
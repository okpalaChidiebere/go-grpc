package api

import (
	"context"
	"log"

	//NOTE: always good to alias your import for good code readability. It is not needed though because you can use the package 'name' on the imported file :)
	apiadapters "github.com/okpalaChidiebere/go-grpc/api/adapters"
	todos "github.com/okpalaChidiebere/go-grpc/businessLogic/todos"
	pb "github.com/okpalaChidiebere/go-grpc/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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

func (s *TodoServer) ReadTodos(ctx context.Context, in *emptypb.Empty) (*pb.ReadTodosResponse, error) {
	todos := apiadapters.TodosToProto(s.TodoService.ReadTodos())
	return &pb.ReadTodosResponse{ Items: todos } , nil
}

func (s *TodoServer) ReadTodosStream(req *emptypb.Empty, stream pb.TodoService_ReadTodosStreamServer) error{
	for _, todo := range s.TodoService.ReadTodos(){
		if err := stream.Send(apiadapters.TodoToProto(&todo)); err != nil {
            log.Fatalf("%v.Send(%v) = %v: ", stream, todo, err)
        }
	}
	return nil
}
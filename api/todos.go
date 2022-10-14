package api

import (
	"bytes"
	"context"
	"io"
	"log"

	//NOTE: always good to alias your import for good code readability. It is not needed though because you can use the package 'name' on the imported file :)
	apiadapters "github.com/okpalaChidiebere/go-grpc/api/adapters"
	todos "github.com/okpalaChidiebere/go-grpc/businessLogic/todos"
	pb "github.com/okpalaChidiebere/go-grpc/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const maxImageSie = 1 << 20 //max imageSize of 1MB. if you want max of 10MB fileSize, you will have 10 << 20

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

func (s *TodoServer) AddPhoto(stream pb.TodoService_AddPhotoServer) error {

	/*
	One way we wan go about receviing data from oneof in protobuf  is using switch statement 
	@see https://developers.google.com/protocol-buffers/docs/reference/go-generated
	This way we dont call the `Recv`` method twice like we are doing below :)
	*/

	req, err := stream.Recv()
	if err != nil {
		return logError(status.Errorf(codes.Unknown, "Cannot receive request: %v\n", err))
	}

	todoId := req.GetInfo().GetTodoId()
	log.Printf("Received request to upload an image for todo: %s", todoId)

	imageData := bytes.NewBuffer(make([] byte,0)) //new Buffer is more memory efficient than bytes.Buffer{} because of pointer
	imageSize := 0

	//start handling receiving image in chunks from the stream
	for {
		m, err := stream.Recv()
		if err == io.EOF {
			log.Printf("No more data from the stream\n")
			break
		} else if err != nil {
			return logError(status.Errorf(codes.Unknown, "Cannot receive Photo Data: %v\n", err))
		}

		//the pieces of the image
		var chunk []byte = m.GetData()
		size := len(chunk)
		imageSize += size
		if imageSize > maxImageSie { //check if the user is uploading an image that is too large
			return logError(status.Errorf(codes.InvalidArgument, "Image is too Large: %d > %d", imageSize, maxImageSie))
		}

		_, err = imageData.Write(chunk)
		if err != nil {
			//return error if an internal server error occurs
			return logError(status.Errorf(codes.Internal, "Cannot write Photo Data: %v", err))
		}
	}
	//end handling receiving image in chunks from the stream

	//At this point we have collected all Photo Data bytes from the stream buffer
	var res *pb.AddTodoPhotoResponse
	if err := s.TodoService.SaveTodoPhoto(todoId, imageData); err != nil {
		res = &pb.AddTodoPhotoResponse{
			Success: false,
			Size: uint32(imageSize),
		}
    }else {
		res = &pb.AddTodoPhotoResponse{
			Success: true,
			Size: uint32(imageSize),
		}
	}
	return stream.SendAndClose(res)
}

/*
Helper method to help log error on our server. It also will return the
 error incase you want to send the error logged back to the client
*/
func logError(err error) error {
	if err != nil {
		log.Print(err)
	}
	return err
}
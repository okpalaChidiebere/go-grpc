package api_adapters

import (
	"github.com/okpalaChidiebere/go-grpc/models"
	pb "github.com/okpalaChidiebere/go-grpc/pb"
)


func TodoToProto (t *models.Todo) *pb.TodoItem{
	return &pb.TodoItem{
		Id: t.Id,
		Text: t.Text,
	}
}
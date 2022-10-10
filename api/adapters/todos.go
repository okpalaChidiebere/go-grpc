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

func TodosToProto (ts []models.Todo) []*pb.TodoItem{
	var todos []*pb.TodoItem
	for _, t := range ts {
		todos = append(todos, &pb.TodoItem{
			Id: t.Id,
			Text: t.Text,
		})
	}
	return todos
}
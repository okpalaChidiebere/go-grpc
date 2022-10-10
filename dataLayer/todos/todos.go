package todoAccess

import (
	"crypto/rand"
	"fmt"
	"log"

	"github.com/okpalaChidiebere/go-grpc/models"
)

type Repository interface {
	Create(t string) (models.Todo)
}

type TodosListDataRepository struct {
	todos []models.Todo
}

func NewTodosListDataRepo() Repository {
	todos := make([]models.Todo, 0)

	return &TodosListDataRepository{todos}
}


func (r *TodosListDataRepository) Create(t string) (models.Todo){
	var todoId string = generateId()
	todoItem := models.Todo{
		Id: todoId,
		Text: t,
	} 
	r.todos = append(r.todos, todoItem)
	return todoItem
}


func generateId() string{
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:])

	return uuid
}
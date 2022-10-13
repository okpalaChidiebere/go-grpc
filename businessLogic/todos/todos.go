package todos

import (
	"bytes"
	"fmt"
	"os"

	todosRepo "github.com/okpalaChidiebere/go-grpc/dataLayer/todos"
	"github.com/okpalaChidiebere/go-grpc/models"
)

/*Other developers might call this Service*/
type Service interface {
	Create(t string) (models.Todo)
	ReadTodos() ([]models.Todo)
	SaveTodoPhoto(todoId string, imageData *bytes.Buffer) (error)
}

type ServiceImpl struct {
	todosRepo todosRepo.Repository
}

func New(repo todosRepo.Repository) Service {
	return &ServiceImpl{repo}
}

func (g *ServiceImpl) Create(t string) (models.Todo) {
	return g.todosRepo.Create(t)
}

func (g *ServiceImpl) ReadTodos() ([]models.Todo) {
	return g.todosRepo.GetMulti()
}

func (g *ServiceImpl) SaveTodoPhoto(todoId string, imageData *bytes.Buffer) error {
	path := fmt.Sprintf("%s/%s.png", "tmp" , todoId)
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("cannot create image file %w", err)
	}
	_, err = imageData.WriteTo(file)
	if err != nil {
		return fmt.Errorf("cannot write image file %w", err)
	}
	return nil
}

package todos

import (
	todosRepo "github.com/okpalaChidiebere/go-grpc/dataLayer/todos"
	"github.com/okpalaChidiebere/go-grpc/models"
)

/*Other developers might call this Service*/
type Service interface {
	Create(t string) (models.Todo)
	ReadTodos() ([]models.Todo)
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
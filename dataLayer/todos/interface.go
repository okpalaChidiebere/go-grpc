package todoAccess

import "github.com/okpalaChidiebere/go-grpc/models"

/*
We define the interface. You can further choose to divide this interface
into a Reader and Writer for more separation of concerns
*/

type Repository interface {
	Create(t string) (models.Todo)
	GetMulti() ([]models.Todo)
}

package interfaces

import "github.com/sunr3d/basic-marketplace/models"

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=UserService --output=../../mocks
type UserService interface {
	RegisterUser(login, password string) (*models.User, error)
}

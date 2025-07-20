package interfaces

import "github.com/sunr3d/basic-marketplace/models"

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=UserRepo --output=../../../mocks
type UserRepo interface {
	CreateUser(user *models.User) error
	GetUserByLogin(login string) (*models.User, error)
}

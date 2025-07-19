package interfaces

import "github.com/sunr3d/basic-marketplace/models"

type UserRepo interface {
	CreateUser(user *models.User) error
	GetUserByLogin(login string) (*models.User, error)
}

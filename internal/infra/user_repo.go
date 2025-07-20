package infra

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/user"
	"github.com/sunr3d/basic-marketplace/models"
)

var _ interfaces.UserRepo = (*UserRepoPG)(nil)

type UserRepoPG struct {
	db *gorm.DB
}

func NewUserRepoPG(db *gorm.DB) interfaces.UserRepo {
	return &UserRepoPG{db: db}
}

func (r *UserRepoPG) CreateUser(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepoPG) GetUserByLogin(login string) (*models.User, error) {
	var user models.User
	err := r.db.Where("login = ?", login).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("не удалось найти пользователя \"%s\": %w", login, err)
		}
		return nil, fmt.Errorf("ошибка при поиске пользователя \"%s\": %w", login, err)
	}

	return &user, nil
}

package logic

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/sunr3d/basic-marketplace/internal/interfaces"
	"github.com/sunr3d/basic-marketplace/models"
)

var _ interfaces.UserService = (*userService)(nil)

type userService struct {
	UserRepo interfaces.UserRepo
}

func NewUserService(userRepo interfaces.UserRepo) interfaces.UserService {
	return &userService{UserRepo: userRepo}
}

func (s *userService) RegisterUser(login, password string) error {
	if err := validateLogin(login); err != nil {
		return fmt.Errorf("validateLogin: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return fmt.Errorf("validatePassword: %w", err)
	}

	_, err := s.UserRepo.GetUserByLogin(login)
	if err == nil {
		return fmt.Errorf("пользователь \"%s\" уже существует", login)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("не удалось захэшировать пароль: %w", err)
	}

	user := &models.User{
		Login:        login,
		PasswordHash: string(hash),
	}
	return s.UserRepo.CreateUser(user)
}

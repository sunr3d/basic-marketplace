package logic

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/sunr3d/basic-marketplace/internal/interfaces"
	"github.com/sunr3d/basic-marketplace/internal/pkg/utils"
	"github.com/sunr3d/basic-marketplace/models"
)

var _ interfaces.UserService = (*userService)(nil)

type userService struct {
	UserRepo interfaces.UserRepo
	JWTSecret []byte
}

func NewUserService(userRepo interfaces.UserRepo, jwtSecret []byte) interfaces.UserService {
	return &userService{UserRepo: userRepo, JWTSecret: jwtSecret}
}

func (s *userService) RegisterUser(login, password string) (*models.User, error) {
	if err := validateLogin(login); err != nil {
		return nil, fmt.Errorf("validateLogin: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return nil, fmt.Errorf("validatePassword: %w", err)
	}

	user, err := s.UserRepo.GetUserByLogin(login)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("GetUserByLogin: %w", err)
	}
	if user != nil {
		return nil, fmt.Errorf("пользователь \"%s\" уже существует", login)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("не удалось захэшировать пароль: %w", err)
	}

	user = &models.User{
		Login:        login,
		PasswordHash: string(hash),
	}
	err = s.UserRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("CreateUser: %w", err)
	}
	
	return user, nil
}

func (s *userService) AuthUser(login, password string) (string, error) {
	user, err := s.UserRepo.GetUserByLogin(login)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", fmt.Errorf("неверный логин или пароль")
		}
		return "", fmt.Errorf("не удалось авторизоваться: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", fmt.Errorf("неверный логин или пароль")
	}

	token, err := utils.GenerateJWT(user.ID, user.Login, s.JWTSecret)
	if err != nil {
		return "", fmt.Errorf("не удалось сгенерировать токен: %w", err)
	}
	
	return token, nil
}

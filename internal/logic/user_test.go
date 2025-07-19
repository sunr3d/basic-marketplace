package logic

import (
	"fmt"
	"testing"

	"gorm.io/gorm"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/sunr3d/basic-marketplace/mocks"
	"github.com/sunr3d/basic-marketplace/models"
)

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserByLogin", "newuser").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(nil)

	err := service.RegisterUser("newuser", "ValidPass123!")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_AlreadyExists(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserByLogin", "existinguser").Return(&models.User{Login: "existinguser"}, nil)

	err := service.RegisterUser("existinguser", "ValidPass123!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "уже существует")
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_InvalidLogin(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	service := NewUserService(mockRepo)

	err := service.RegisterUser("ab", "ValidPass123!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "логин должен содержать от")
	mockRepo.AssertNotCalled(t, "GetUserByLogin", mock.Anything)
	mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestRegisterUser_InvalidPassword(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	service := NewUserService(mockRepo)

	err := service.RegisterUser("validuser", "short")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "пароль должен содержать от")
	mockRepo.AssertNotCalled(t, "GetUserByLogin", mock.Anything)
	mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestRegisterUser_RepoErrorOnGetUser(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserByLogin", "user1").Return(nil, fmt.Errorf("ошибка БД"))

	err := service.RegisterUser("user1", "ValidPass123!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка БД")
	mockRepo.AssertNotCalled(t, "CreateUser", mock.Anything)
}

func TestRegisterUser_RepoErrorOnCreateUser(t *testing.T) {
	mockRepo := new(mocks.UserRepo)
	service := NewUserService(mockRepo)

	mockRepo.On("GetUserByLogin", "user2").Return(nil, gorm.ErrRecordNotFound)
	mockRepo.On("CreateUser", mock.AnythingOfType("*models.User")).Return(fmt.Errorf("ошибка БД"))

	err := service.RegisterUser("user2", "ValidPass123!")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ошибка БД")
}
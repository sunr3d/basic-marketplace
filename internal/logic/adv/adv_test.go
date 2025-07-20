package logic

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	"github.com/sunr3d/basic-marketplace/mocks"
	"github.com/sunr3d/basic-marketplace/models"
)

func TestCreateAd_Success(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	input := interfaces.AdInput{
		AdvBase: interfaces.AdvBase{
			Title:       "Продам велосипед",
			Description: "Почти новый",
			ImageURL:    "https://example.com/image.jpg",
			Price:       10000,
		},
		OwnerID: 1,
	}
	expected := &models.Adv{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		OwnerID:     input.OwnerID,
	}
	mockRepo.On("CreateAdv", mock.AnythingOfType("*models.Adv")).Return(expected, nil)

	adv, err := service.CreateAd(input)
	assert.NoError(t, err)
	assert.NotNil(t, adv)
	assert.Equal(t, input.Title, adv.Title)
	mockRepo.AssertExpectations(t)
}

func TestCreateAd_ValidationError(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	input := interfaces.AdInput{
		AdvBase: interfaces.AdvBase{
			Title:       "П",
			Description: "Почти новый",
			ImageURL:    "https://example.com/image.jpg",
			Price:       10000,
		},
		OwnerID: 1,
	}

	adv, err := service.CreateAd(input)
	assert.Error(t, err)
	assert.Nil(t, adv)
	mockRepo.AssertNotCalled(t, "CreateAdv", mock.Anything)
}

func TestCreateAd_RepoError(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	input := interfaces.AdInput{
		AdvBase: interfaces.AdvBase{
			Title:       "Продам велосипед",
			Description: "Почти новый",
			ImageURL:    "https://example.com/image.jpg",
			Price:       10000,
		},
		OwnerID: 1,
	}
	mockRepo.On("CreateAdv", mock.AnythingOfType("*models.Adv")).Return(nil, errors.New("ошибка БД"))

	adv, err := service.CreateAd(input)
	assert.Error(t, err)
	assert.Nil(t, adv)
	mockRepo.AssertExpectations(t)
}

func TestCreateAd_EmptyImageURL(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	input := interfaces.AdInput{
		AdvBase: interfaces.AdvBase{
			Title:       "Продам велосипед",
			Description: "Почти новый",
			ImageURL:    "",
			Price:       10000,
		},
		OwnerID: 1,
	}
	expected := &models.Adv{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		OwnerID:     input.OwnerID,
	}
	mockRepo.On("CreateAdv", mock.AnythingOfType("*models.Adv")).Return(expected, nil)

	adv, err := service.CreateAd(input)
	assert.NoError(t, err)
	assert.NotNil(t, adv)
	assert.Equal(t, "", adv.ImageURL)
	mockRepo.AssertExpectations(t)
}

func TestCreateAd_MinPrice(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	input := interfaces.AdInput{
		AdvBase: interfaces.AdvBase{
			Title:       "Продам велосипед",
			Description: "Почти новый",
			ImageURL:    "https://example.com/image.jpg",
			Price:       1.0,
		},
		OwnerID: 1,
	}
	expected := &models.Adv{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		OwnerID:     input.OwnerID,
	}
	mockRepo.On("CreateAdv", mock.AnythingOfType("*models.Adv")).Return(expected, nil)

	adv, err := service.CreateAd(input)
	assert.NoError(t, err)
	assert.NotNil(t, adv)
	assert.Equal(t, 1.0, adv.Price)
	mockRepo.AssertExpectations(t)
}

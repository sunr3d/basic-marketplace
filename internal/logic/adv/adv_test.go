package logic

import (
	"errors"
	"testing"
	"time"

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

func TestShowAdsFeed_Success(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	filter := interfaces.AdvFilter{
		SortBy:   "created_at",
		Order:    "desc",
		MinPrice: 0,
		MaxPrice: 0,
	}
	currentUserID := uint(21)
	advs := []interfaces.AdvWithOwner{
		{
			Adv: models.Adv{
				ID: 1, Title: "Товар 1", Description: "Описание 1", ImageURL: "https://img.com/1.jpg", Price: 1000, OwnerID: 21,
				CreatedAt: time.Now(),
			},
			OwnerLogin: "user21",
		},
		{
			Adv: models.Adv{
				ID: 2, Title: "Товар 2", Description: "Описание 2", ImageURL: "https://img.com/2.jpg", Price: 2000, OwnerID: 7,
				CreatedAt: time.Now(),
			},
			OwnerLogin: "user7",
		},
	}
	mockRepo.On("FindMany", filter).Return(advs, nil)

	feed, err := service.ShowAdsFeed(filter, currentUserID)
	assert.NoError(t, err)
	assert.Len(t, feed, 2)
	assert.True(t, feed[0].IsOwner)
	assert.False(t, feed[1].IsOwner)
	assert.Equal(t, "user21", feed[0].OwnerLogin)
	assert.Equal(t, "user7", feed[1].OwnerLogin)
	mockRepo.AssertExpectations(t)
}

func TestShowAdsFeed_ValidationError(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	filter := interfaces.AdvFilter{
		SortBy: "unknown_field",
	}
	feed, err := service.ShowAdsFeed(filter, 1)
	assert.Error(t, err)
	assert.Nil(t, feed)
	mockRepo.AssertNotCalled(t, "FindMany", mock.Anything)
}

func TestShowAdsFeed_RepoError(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	filter := interfaces.AdvFilter{
		SortBy: "created_at",
	}
	mockRepo.On("FindMany", filter).Return(nil, errors.New("ошибка БД"))

	feed, err := service.ShowAdsFeed(filter, 1)
	assert.Error(t, err)
	assert.Nil(t, feed)
	mockRepo.AssertExpectations(t)
}

func TestShowAdsFeed_EmptyResult(t *testing.T) {
	mockRepo := new(mocks.AdvRepo)
	service := NewAdvService(mockRepo)

	filter := interfaces.AdvFilter{
		SortBy: "created_at",
	}
	mockRepo.On("FindMany", filter).Return([]interfaces.AdvWithOwner{}, nil)

	feed, err := service.ShowAdsFeed(filter, 1)
	assert.NoError(t, err)
	assert.Empty(t, feed)
	mockRepo.AssertExpectations(t)
}

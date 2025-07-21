package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	"github.com/sunr3d/basic-marketplace/mocks"
)

func setupFeedRouterWithAuth(svc interfaces.AdvService, log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ads", func(c *gin.Context) {
		c.Set("user_id", float64(21))
		FeedHandler(svc, log)(c)
	})
	return router
}

func TestFeedHandler_Success(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	filter := interfaces.AdvFilter{
		Limit:  20,
		Offset: 0,
	}
	expected := []interfaces.AdvFeedItem{
		{
			ID:         1,
			AdvBase:    interfaces.AdvBase{Title: "Товар", Description: "Описание", ImageURL: "https://img.com/1.jpg", Price: 1000},
			OwnerLogin: "user21",
			IsOwner:    true,
			CreatedAt:  "2025-07-21 12:00:00",
		},
	}
	mockService.On("ShowAdsFeed", filter, uint(21)).Return(expected, nil)

	router := setupFeedRouterWithAuth(mockService, log)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ads", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []interfaces.AdvFeedItem
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)
	assert.Equal(t, "user21", resp[0].OwnerLogin)
	mockService.AssertExpectations(t)
}

func TestFeedHandler_ValidationError(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	filter := interfaces.AdvFilter{
		Limit:  20,
		Offset: 0,
		Order:  "badorder",
	}
	mockService.On("ShowAdsFeed", filter, uint(21)).Return(nil, errors.New("может быть ошибка валидации"))

	router := setupFeedRouterWithAuth(mockService, log)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ads?order=badorder", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestFeedHandler_InternalError(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	filter := interfaces.AdvFilter{
		Limit:  20,
		Offset: 0,
	}
	mockService.On("ShowAdsFeed", filter, uint(21)).Return(nil, errors.New("ошибка БД"))

	router := setupFeedRouterWithAuth(mockService, log)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ads", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

func TestFeedHandler_EmptyFeed(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	filter := interfaces.AdvFilter{
		Limit:  20,
		Offset: 0,
	}
	mockService.On("ShowAdsFeed", filter, uint(21)).Return([]interfaces.AdvFeedItem{}, nil)

	router := setupFeedRouterWithAuth(mockService, log)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ads", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []interfaces.AdvFeedItem
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Empty(t, resp)
	mockService.AssertExpectations(t)
}

func TestFeedHandler_NoAuth(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/ads", FeedHandler(mockService, log))

	filter := interfaces.AdvFilter{
		Limit:  20,
		Offset: 0,
	}
	expected := []interfaces.AdvFeedItem{
		{
			ID:         1,
			AdvBase:    interfaces.AdvBase{Title: "Товар", Description: "Описание", ImageURL: "https://img.com/1.jpg", Price: 1000},
			OwnerLogin: "user21",
			IsOwner:    false,
			CreatedAt:  "2025-07-21 12:00:00",
		},
		{
			ID:         2,
			AdvBase:    interfaces.AdvBase{Title: "Товар2", Description: "Описание2", ImageURL: "https://img.com/2.jpg", Price: 2000},
			OwnerLogin: "user7",
			IsOwner:    false,
			CreatedAt:  "2025-07-21 12:01:00",
		},
	}

	mockService.On("ShowAdsFeed", filter, uint(0)).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ads", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []interfaces.AdvFeedItem
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	for _, item := range resp {
		assert.False(t, item.IsOwner)
	}
	mockService.AssertExpectations(t)
}

func TestFeedHandler_Authorized_IsOwnerCorrect(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()
	router := setupFeedRouterWithAuth(mockService, log)

	filter := interfaces.AdvFilter{
		Limit:  20,
		Offset: 0,
	}
	expected := []interfaces.AdvFeedItem{
		{
			ID:         1,
			AdvBase:    interfaces.AdvBase{Title: "Товар", Description: "Описание", ImageURL: "https://img.com/1.jpg", Price: 1000},
			OwnerLogin: "user21",
			IsOwner:    true,
			CreatedAt:  "2025-07-21 12:00:00",
		},
		{
			ID:         2,
			AdvBase:    interfaces.AdvBase{Title: "Товар2", Description: "Описание2", ImageURL: "https://img.com/2.jpg", Price: 2000},
			OwnerLogin: "user7",
			IsOwner:    false,
			CreatedAt:  "2025-07-21 12:01:00",
		},
	}

	mockService.On("ShowAdsFeed", filter, uint(21)).Return(expected, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ads", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp []interfaces.AdvFeedItem
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 2)
	assert.True(t, resp[0].IsOwner)
	assert.False(t, resp[1].IsOwner)
	mockService.AssertExpectations(t)
}

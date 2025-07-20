package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	"github.com/sunr3d/basic-marketplace/mocks"
	"github.com/sunr3d/basic-marketplace/models"
)

func setupRouterWithAuth(svc interfaces.AdvService, log *zap.Logger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/ad", func(c *gin.Context) {
		c.Set("user_id", float64(21))
		CreateAdvHandler(svc, log)(c)
	})
	return router
}

func TestCreateAdvHandler_Success(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	reqBody := map[string]any{
		"title":       "Продам велосипед",
		"description": "Почти новый",
		"image_url":   "https://example.com/image.jpg",
		"price":       10000,
	}
	expected := &models.Adv{
		ID:          1,
		Title:       "Продам велосипед",
		Description: "Почти новый",
		ImageURL:    "https://example.com/image.jpg",
		Price:       10000,
		OwnerID:     21,
		CreatedAt:   time.Now(),
	}
	mockService.On("CreateAd", mock.AnythingOfType("interfaces.AdInput")).Return(expected, nil)

	router := setupRouterWithAuth(mockService, log)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ad", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Продам велосипед", resp["title"])
	assert.Equal(t, float64(21), resp["owner_id"])
	mockService.AssertExpectations(t)
}

func TestCreateAdvHandler_BadRequest(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()
	router := setupRouterWithAuth(mockService, log)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ad", bytes.NewReader([]byte("{")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateAdvHandler_ValidationError(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	mockService.On("CreateAd", mock.AnythingOfType("interfaces.AdInput")).Return(nil, fmt.Errorf("validateAdInput: заголовок слишком короткий"))

	router := setupRouterWithAuth(mockService, log)
	reqBody := map[string]any{
		"title":       "П",
		"description": "Почти новый",
		"image_url":   "https://example.com/image.jpg",
		"price":       10000,
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ad", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

func TestCreateAdvHandler_NoAuth(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/ad", CreateAdvHandler(mockService, log))

	reqBody := map[string]any{
		"title":       "Продам велосипед",
		"description": "Почти новый",
		"image_url":   "https://example.com/image.jpg",
		"price":       10000,
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ad", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestCreateAdvHandler_EmptyImageURL(t *testing.T) {
	mockService := new(mocks.AdvService)
	log := zap.NewNop()

	reqBody := map[string]any{
		"title":       "Продам велосипед",
		"description": "Почти новый",
		"image_url":   "",
		"price":       10000,
	}
	expected := &models.Adv{
		ID:          2,
		Title:       "Продам велосипед",
		Description: "Почти новый",
		ImageURL:    "",
		Price:       10000,
		OwnerID:     21,
		CreatedAt:   time.Now(),
	}
	mockService.On("CreateAd", mock.AnythingOfType("interfaces.AdInput")).Return(expected, nil)

	router := setupRouterWithAuth(mockService, log)
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/ad", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "", resp["image_url"])
	mockService.AssertExpectations(t)
}

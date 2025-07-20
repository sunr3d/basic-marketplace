package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/mocks"
)

func TestLoginHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()

	mockService.On("AuthUser", "user1", "password123").Return("sometoken321", nil)

	router := gin.New()
	router.POST("/login", LoginHandler(mockService, log))

	reqBody := map[string]string{
		"login":    "user1",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "Bearer sometoken321", w.Header().Get("Authorization"))
	var resp map[string]string
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "sometoken321", resp["jwt"])
	mockService.AssertExpectations(t)
}

func TestLoginHandler_InvalidCreds(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()

	mockService.On("AuthUser", "baduser", "badpass").Return("", fmt.Errorf("неверный логин или пароль"))

	router := gin.New()
	router.POST("/login", LoginHandler(mockService, log))

	reqBody := map[string]string{
		"login":    "baduser",
		"password": "badpass",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	mockService.AssertExpectations(t)
}

func TestLoginHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()

	router := gin.New()
	router.POST("/login", LoginHandler(mockService, log))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader([]byte("{")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestLoginHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()

	mockService.On("AuthUser", "user1", "password123").Return("", fmt.Errorf("ошибка БД"))

	router := gin.New()
	router.POST("/login", LoginHandler(mockService, log))

	reqBody := map[string]string{
		"login":    "user1",
		"password": "password123",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	mockService.AssertExpectations(t)
}

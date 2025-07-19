package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/mocks"
	"github.com/sunr3d/basic-marketplace/models"
)

func TestRegisterHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()

	mockService.On("RegisterUser", "newuser", "ValidPass123!").Return(&models.User{
		ID:        1,
		Login:     "newuser",
		CreatedAt: time.Now(),
	}, nil)

	router := gin.New()
	router.POST("/register", RegisterHandler(mockService, log))

	reqBody := map[string]string{
		"login":    "newuser",
		"password": "ValidPass123!",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]any
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "newuser", resp["login"])
	mockService.AssertExpectations(t)
}

func TestRegisterHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()
	router := gin.New()
	router.POST("/register", RegisterHandler(mockService, log))

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader([]byte("{")))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegisterHandler_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(mocks.UserService)
	log := zap.NewNop()
	mockService.On("RegisterUser", "baduser", "badpass").Return(nil, assert.AnError)
	router := gin.New()
	router.POST("/register", RegisterHandler(mockService, log))

	reqBody := map[string]string{
		"login": "baduser",
		"password": "badpass",
	}
	body, _ := json.Marshal(reqBody)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockService.AssertExpectations(t)
}

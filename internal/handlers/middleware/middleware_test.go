package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"

	"github.com/sunr3d/basic-marketplace/internal/pkg/utils"
)

var testSecret = []byte("testsecrettoken")

func TestAuthMiddleware_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthMiddleware(testSecret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthMiddleware(testSecret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthMiddleware(testSecret), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		login, _ := c.Get("login")
		c.JSON(http.StatusOK, gin.H{
			"user_id": userID,
			"login":   login,
		})
	})

	token, _ := utils.GenerateJWT(21, "testuser", testSecret)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "testuser")
	assert.Contains(t, w.Body.String(), "21")
}

func TestAuthMiddleware_NoIDClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthMiddleware(testSecret), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	claims := jwt.MapClaims{
		"login": "testuser",
		"exp":   9999999999,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, _ := token.SignedString(testSecret)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenStr)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Внутренняя ошибка сервера")
}

func TestAuthMiddleware_PanicRec(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/protected", AuthMiddleware(testSecret), func(c *gin.Context) {
		panic("FLY YOU FOOLS")
	})

	token, _ := utils.GenerateJWT(21, "hobbitses", testSecret)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Внутренняя ошибка сервера")

}

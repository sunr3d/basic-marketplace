package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/user"
	"go.uber.org/zap"
)

type loginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginResp struct {
	JWT string `json:"jwt"`
}

func LoginHandler(userService interfaces.UserService, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req loginReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат запроса"})
			return
		}

		token, err := userService.AuthUser(req.Login, req.Password)
		if err != nil {
			if isAuthError(err) {
				log.Warn("Ошибка авторизации пользователя", zap.Error(err))
				c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
				return
			}
			log.Error("Внутреняя ошибка при авторизации пользователя", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутреняя ошибка сервера"})
			return
		}

		log.Info("Пользователь произвел успешную авторизацию", zap.String("login", req.Login))
		c.Header("Authorization", "Bearer "+token)
		c.JSON(http.StatusOK, loginResp{JWT: token})
	}
}

func isAuthError(err error) bool {
	return strings.Contains(err.Error(), "неверный логин или пароль")
}

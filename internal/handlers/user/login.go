package handlers

import (
	"net/http"

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
			log.Warn("Ошибка авторизации пользователя", zap.Error(err))
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.Header("Authorization", "Bearer "+token)
		c.JSON(http.StatusOK, loginResp{JWT: token})
	}
}

package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/user"
)

type registerReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type registerResp struct {
	ID        uint      `json:"id"`
	Login     string    `json:"login"`
	CreatedAt time.Time `json:"created_at"`
}

func RegisterHandler(userService interfaces.UserService, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req registerReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорретный формат запроса"})
			return
		}

		user, err := userService.RegisterUser(req.Login, req.Password)
		if err != nil {
			log.Warn("Ошибка регистрации пользователя", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		resp := registerResp{
			ID:        user.ID,
			Login:     user.Login,
			CreatedAt: user.CreatedAt,
		}
		c.JSON(http.StatusOK, resp)
	}
}

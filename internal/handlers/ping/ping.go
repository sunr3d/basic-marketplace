package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func PingHandler(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Info("Вызов Пинг эндпоинта")
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	}
}

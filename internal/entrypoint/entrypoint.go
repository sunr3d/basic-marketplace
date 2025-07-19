package entrypoint

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/internal/config"
)

func Run(cfg *config.Config, log *zap.Logger) error {
	router := gin.New()
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		log.Info("Получен пинг")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	addr := fmt.Sprintf("%s:%s", cfg.HTTPHost, cfg.HTTPPort)
	go func() {
		log.Info("Запуск сервера", zap.String("addr", addr))
		if err := router.Run(addr); err != nil {
			log.Fatal("Ошибка запуска сервера", zap.Error(err))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	<-done
	log.Info("Получен сигнал стоп")

	return nil
}

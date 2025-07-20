package entrypoint

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/internal/bootstrap"
	"github.com/sunr3d/basic-marketplace/internal/config"
	"github.com/sunr3d/basic-marketplace/internal/handlers"
	"github.com/sunr3d/basic-marketplace/models"
)

func Run(cfg *config.Config, log *zap.Logger) error {
	/// 1. Слой репозиториев (infra)
	// Создание контейнера зависимостей (Dependency Injection) + миграция
	container, err := bootstrap.NewContainer(cfg, log)
	if err != nil {
		return fmt.Errorf("bootstrap.NewContainer: %w", err)
	}
	log.Info("Контейнер зависимостей успешно создан")

	if err := container.DB.AutoMigrate(&models.User{}); err != nil {
		return fmt.Errorf("db.AutoMigrate: %w", err)
	}
	log.Info("Миграция модели User прошла успешно")

	// ЧТО-ТО
	router := gin.New()
	router.Use(gin.Recovery())

	// Регистрация ручек на роутере
	router.GET("/ping", handlers.PingHandler(log))
	router.POST("/register", handlers.RegisterHandler(container.UserService, log))
	router.POST("/login", handlers.LoginHandler(container.UserService, log))

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

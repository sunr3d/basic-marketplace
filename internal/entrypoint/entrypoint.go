package entrypoint

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/sunr3d/basic-marketplace/internal/bootstrap"
	"github.com/sunr3d/basic-marketplace/internal/config"
	adv_handlers "github.com/sunr3d/basic-marketplace/internal/handlers/adv"
	middleware "github.com/sunr3d/basic-marketplace/internal/handlers/middleware"
	ping_handlers "github.com/sunr3d/basic-marketplace/internal/handlers/ping"
	user_handlers "github.com/sunr3d/basic-marketplace/internal/handlers/user"
	"github.com/sunr3d/basic-marketplace/internal/server"
	"github.com/sunr3d/basic-marketplace/models"
)

func Run(cfg *config.Config, log *zap.Logger) error {
	// Создание контейнера зависимостей (Dependency Injection) + миграция
	container, err := bootstrap.NewContainer(cfg, log)
	if err != nil {
		return fmt.Errorf("bootstrap.NewContainer: %w", err)
	}
	log.Info("Контейнер зависимостей успешно создан")

	if err := container.DB.AutoMigrate(&models.User{}, &models.Adv{}); err != nil {
		return fmt.Errorf("db.AutoMigrate: %w", err)
	}
	log.Info("Миграция моделей прошла успешно")

	// Создание гин роутера
	router := gin.New()
	router.Use(gin.Recovery())

	// Регистрация ручек на роутере
	router.GET("/ping", ping_handlers.PingHandler(log))
	router.POST("/register", user_handlers.RegisterHandler(container.UserService, log))
	router.POST("/login", user_handlers.LoginHandler(container.UserService, log))
	router.GET("/ads", adv_handlers.FeedHandler(container.AdvService, log))

	auth := router.Group("/")
	auth.Use(middleware.AuthMiddleware([]byte(cfg.JWTSecret)))
	{
		auth.POST("/ads/create", adv_handlers.CreateAdvHandler(container.AdvService, log))
	}

	// 2. Создание сервера
	srv := server.New(router, cfg, log)
	return srv.Start()
}

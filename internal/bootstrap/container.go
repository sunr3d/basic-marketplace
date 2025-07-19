package bootstrap

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sunr3d/basic-marketplace/internal/config"
	"github.com/sunr3d/basic-marketplace/internal/infra"
	"github.com/sunr3d/basic-marketplace/internal/interfaces"
	"github.com/sunr3d/basic-marketplace/internal/logic"
)

type Container struct {
	DB       *gorm.DB
	UserRepo interfaces.UserRepo
	UserService interfaces.UserService
	// ... другие зависимости, добавим позже
}

func NewContainer(cfg *config.Config, log *zap.Logger) (*Container, error) {
	db, err := infra.InitPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("infra.InitPostgres: %w", err)
	}
	userRepo := infra.NewUserRepoPG(db)
	userService := logic.NewUserService(userRepo)
	// ... еще будут зависимости, типа Редиса для кеширования

	return &Container{
		DB:       db,
		UserRepo: userRepo,
		UserService: userService,
	}, nil
}

package bootstrap

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/sunr3d/basic-marketplace/internal/config"
	"github.com/sunr3d/basic-marketplace/internal/infra"
	adv_interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	user_interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/user"
	user_logic "github.com/sunr3d/basic-marketplace/internal/logic/user"
	adv_logic "github.com/sunr3d/basic-marketplace/internal/logic/adv"
)

type Container struct {
	DB          *gorm.DB
	UserRepo    user_interfaces.UserRepo
	UserService user_interfaces.UserService
	AdvRepo     adv_interfaces.AdvRepo
	AdvService  adv_interfaces.AdvService
	// ... другие зависимости, добавим позже
}

func NewContainer(cfg *config.Config, log *zap.Logger) (*Container, error) {
	db, err := infra.InitPostgres(cfg.Postgres)
	if err != nil {
		return nil, fmt.Errorf("infra.InitPostgres: %w", err)
	}
	userRepo := infra.NewUserRepoPG(db)
	userService := user_logic.NewUserService(userRepo, []byte(cfg.JWTSecret))
	advRepo := infra.NewAdvRepoPG(db)
	advService := adv_logic.NewAdvService(advRepo)
	// ... еще будут зависимости, типа Редиса для кеширования

	return &Container{
		DB:          db,
		UserRepo:    userRepo,
		UserService: userService,
		AdvRepo:     advRepo,
		AdvService:  advService,
	}, nil
}

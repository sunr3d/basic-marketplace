package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
)

type createAdvReq struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

type createAdvResp struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ImageURL    string    `json:"image_url"`
	Price       float64   `json:"price"`
	OwnerID     uint      `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func CreateAdvHandler(adService interfaces.AdvService, log *zap.Logger) gin.HandlerFunc {

}

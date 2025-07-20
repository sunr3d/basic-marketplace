package handlers

import (
	"net/http"
	"strings"
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
	return func(c *gin.Context) {
		var req createAdvReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный формат запроса"})
			return
		}

		userIDVal, exists := c.Get("user_id")
		userID, ok := userIDVal.(float64)
		if !exists || !ok || userID < 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			return
		}

		input := interfaces.AdInput{
			AdvBase: interfaces.AdvBase{
				Title:       req.Title,
				Description: req.Description,
				ImageURL:    req.ImageURL,
				Price:       req.Price,
			},
			OwnerID: uint(userID),
		}

		adv, err := adService.CreateAd(input)
		if err != nil {
			if strings.Contains(err.Error(), "validateAdInput") {
				log.Warn("Ошибка создания объявления", zap.Error(err))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Error("Внутренняя ошибка при создании объявления", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
			return
		}

		resp := createAdvResp{
			ID:          adv.ID,
			Title:       adv.Title,
			Description: adv.Description,
			ImageURL:    adv.ImageURL,
			Price:       adv.Price,
			OwnerID:     adv.OwnerID,
			CreatedAt:   adv.CreatedAt,
		}
		c.JSON(http.StatusOK, resp)
	}
}

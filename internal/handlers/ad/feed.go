package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
)

const defaultLimit = 20

func FeedHandler(adService interfaces.AdvService, log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		userID, ok := userIDVal.(float64)
		if !exists || !ok || userID < 1 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "требуется авторизация"})
			return
		}

		filter := interfaces.AdvFilter{
			Limit:  defaultLimit,
			Offset: 0,
		}

		if minPriceStr := c.Query("min_price"); minPriceStr != "" {
			if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
				filter.MinPrice = minPrice
			}
		}
		if maxPriceStr := c.Query("max_price"); maxPriceStr != "" {
			if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
				filter.MaxPrice = maxPrice
			}
		}
		if order := c.Query("order"); order != "" {
			filter.Order = order
		}

		feed, err := adService.ShowAdsFeed(filter, uint(userID))
		if err != nil {
			if isFeedValidationErr(err) {
				log.Warn("Ошибка валидации фильтра объявлений", zap.Error(err))
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			log.Error("Внутреняя ошибка при получении ленты объявлений", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Внутренняя ошибка сервера"})
			return
		}

		c.JSON(http.StatusOK, feed)
	}
}

func isFeedValidationErr(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "не может быть")
}

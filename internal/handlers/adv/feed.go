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
		var currentUserID uint = 0
		if userIDVal, exists := c.Get("user_id"); exists {
			if userID, ok := userIDVal.(float64); ok && userID > 0 {
				currentUserID = uint(userID)
			}
		}

		// Пагинация
		limit := defaultLimit
		if limitStr := c.Query("limit"); limitStr != "" {
			if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
				limit = l
			}
		}
		offset := 0
		if offsetStr := c.Query("offest"); offsetStr != "" {
			if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
				offset = 0
			}
		}
		if pageStr := c.Query("page"); pageStr != "" {
			if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
				offset = (p - 1) * limit
			}
		}

		filter := interfaces.AdvFilter{
			Limit:  limit,
			Offset: offset,
		}

		// Парсинг настроек фильтра
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

		feed, err := adService.ShowAdsFeed(filter, currentUserID)
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

		log.Info("Лента объявлений успешно получена", zap.Int("count", len(feed)))
		c.JSON(http.StatusOK, feed)
	}
}

func isFeedValidationErr(err error) bool {
	msg := err.Error()
	return strings.Contains(msg, "может быть")
}

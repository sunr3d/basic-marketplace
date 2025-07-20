package logic

import (
	"fmt"
	"regexp"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
)

const (
	minTitleLen  = 4
	maxTitleLen  = 35
	minDescLen   = 4
	maxDescLen   = 100
	minPrice     = 1.0
	maxPrice     = 1000000.0
	imagePattern = `^https?://.*\.(jpg|jpeg|png)$`
)

var imageRegexp = regexp.MustCompile(imagePattern)

func validateAdInput(input interfaces.AdInput) error {
	nTitle := len(input.Title)
	nDesc := len(input.Description)

	if nTitle < minTitleLen || nTitle > maxTitleLen {
		return fmt.Errorf("заголовок должен содержать от %d до %d символов", minTitleLen, maxTitleLen)
	}
	if nDesc < minDescLen || nDesc > maxDescLen {
		return fmt.Errorf("описание должно содержать от %d до %d символов", minDescLen, maxDescLen)
	}
	if input.Price < minPrice || input.Price > maxPrice {
		return fmt.Errorf("цена должен быть от %.2f до %.2f", minPrice, maxPrice)
	}
	if input.ImageURL != "" && !imageRegexp.MatchString(input.ImageURL) {
		return fmt.Errorf("некорректный формат ссылки на изображение")
	}

	return nil
}

func validateAdvFilter(filter interfaces.AdvFilter) error {
	if filter.MinPrice < 0 {
		return fmt.Errorf("минимальная цена не может быть отрицательной")
	}
	if filter.MaxPrice < 0 {
		return fmt.Errorf("максимальная цена не может быть отрицательной")
	}
	if filter.MinPrice > 0 && filter.MaxPrice > 0 && filter.MinPrice < filter.MaxPrice {
		return fmt.Errorf("максимальная цена не может быть меньше минимальной")
	}
	if filter.SortBy != "" && filter.SortBy != "created_at" && filter.SortBy != "price" {
		return fmt.Errorf("sort_by может быть только 'created_at' или 'price'")
	}
	if filter.Order != "" && filter.Order != "asc" && filter.Order != "desc" {
		return fmt.Errorf("order может быть только 'asc' или 'desc'")
	}
	return nil
}

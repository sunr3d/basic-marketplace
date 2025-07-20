package infra

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	"github.com/sunr3d/basic-marketplace/models"
)

var _ interfaces.AdvRepo = (*AdvRepoPG)(nil)

type AdvRepoPG struct {
	db *gorm.DB
}

func NewAdvRepoPG(db *gorm.DB) interfaces.AdvRepo {
	return &AdvRepoPG{db: db}
}

func (r *AdvRepoPG) CreateAdv(adv *models.Adv) (*models.Adv, error) {
	if err := r.db.Create(adv).Error; err != nil {
		return nil, errors.New("не удалось создать объявление")
	}
	return adv, nil
}

func (r *AdvRepoPG) FindMany(filter interfaces.AdvFilter) ([]interfaces.AdvWithOwner, error) {
	var ads []interfaces.AdvWithOwner
	// Добавляем LEFT JOIN к запросу
	db := r.db.Table("advs").
		Select("advs.*, users.login as owner_login").
		Joins("left join users on users.id = advs.owner_id")

	// Фильтрация по цене
	if filter.MinPrice > 0 {
		db = db.Where("price >= ?", filter.MinPrice)
	}
	if filter.MaxPrice > 0 {
		db = db.Where("price <= ?", filter.MaxPrice)
	}

	// Сортировка
	sortBy := "advs.created_at"
	if filter.SortBy == "price" {
		sortBy = "advs.price"
	}
	order := "desc"
	if filter.Order == "asc" {
		order = "asc"
	}
	db = db.Order(sortBy + " " + order)

	// Пагинация
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}

	// Отправка запроса в БД
	if err := db.Find(&ads).Error; err != nil {
		return nil, fmt.Errorf("не удалось получить список объявлений из БД: %w", err)
	}

	return ads, nil
}

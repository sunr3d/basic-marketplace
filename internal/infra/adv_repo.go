package infra

import (
	"errors"

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

func (r *AdvRepoPG) FindMany(filter interfaces.AdvFilter) ([]*models.Adv, error) {
	// TODO: ЗАГЛУШКА ПОКА
	return nil, nil
}

package interfaces

import "github.com/sunr3d/basic-marketplace/models"

type AdvFilter struct {
	Limit int
	Offset int
	SortBy string
	Order string
	MinPrice float64
	MaxPrice float64
}

type AdvRepo interface {
	CreateAdv(adv *models.Adv) (*models.Adv, error)
	FindMany(filter AdvFilter) ([]*models.Adv, error)
}

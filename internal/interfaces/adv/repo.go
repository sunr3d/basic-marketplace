package interfaces

import (
	"github.com/sunr3d/basic-marketplace/models"
)

type AdvFilter struct {
	Limit    int
	Offset   int
	SortBy   string
	Order    string
	MinPrice float64
	MaxPrice float64
}

type AdvWithOwner struct {
	models.Adv
	OwnerLogin string `gorm:"column:owner_login"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=AdvRepo --output=../../../mocks
type AdvRepo interface {
	CreateAdv(adv *models.Adv) (*models.Adv, error)
	FindMany(filter AdvFilter) ([]AdvWithOwner, error)
}

package interfaces

import "github.com/sunr3d/basic-marketplace/models"

type AdvBase struct {
	Title string
	Description string
	ImageURL string
	Price float64
}

type AdInput struct {
	AdvBase
	OwnerID uint
}

type AdvFeedItem struct {
	AdvBase
	ID uint
	OwnerLogin string
	IsOwner bool
	CreatedAt string
}
//go:generate go run github.com/vektra/mockery/v2@v2.53.2 --name=AdvService --output=../../../mocks
type AdvService interface {
	CreateAd(input AdInput) (*models.Adv, error)
	ShowAdsFeed(filter AdvFilter, currentUserID uint) ([]AdvFeedItem, error)
}

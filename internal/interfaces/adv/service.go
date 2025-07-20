package interfaces

import "github.com/sunr3d/basic-marketplace/models"

type AdInput struct {
	Title string
	Description string
	ImageURL string
	Price float64
	OwnerID uint
}

type AdvFeedItem struct {
	ID uint
	Title string
	Description string
	ImageURL string
	Price float64
	OwnerLogin string
	IsOwner bool
	CreatedAt string
}

type AdvService interface {
	CreateAd(input AdInput) (*models.Adv, error)
	ShowAdsFeed(filter AdvFilter, currentUserID uint) ([]AdvFeedItem, error)
}

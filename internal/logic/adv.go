package logic

import (
	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	"github.com/sunr3d/basic-marketplace/models"
)

type advService struct {
	AdvRepo interfaces.AdvRepo
}

func NewAdvService(repo interfaces.AdvRepo) interfaces.AdvService {
	return &advService{AdvRepo: repo}
}

func (s *advService) CreateAd(input interfaces.AdInput) (*models.Adv, error) {
	// TODO: Валидация

	adv := &models.Adv{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		Price:       input.Price,
		OwnerID:     input.OwnerID,
	}
	return s.AdvRepo.CreateAdv(adv)
}

func (s *advService) ShowAdsFeed(filter interfaces.AdvFilter, currentUserID uint) ([]interfaces.AdvFeedItem, error) {
	// TODO: Пока заглушка
	return nil, nil
}

package logic

import (
	"fmt"

	interfaces "github.com/sunr3d/basic-marketplace/internal/interfaces/adv"
	"github.com/sunr3d/basic-marketplace/models"
)

var _ interfaces.AdvService = (*advService)(nil)

type advService struct {
	AdvRepo interfaces.AdvRepo
}

func NewAdvService(repo interfaces.AdvRepo) interfaces.AdvService {
	return &advService{AdvRepo: repo}
}

func (s *advService) CreateAd(input interfaces.AdInput) (*models.Adv, error) {
	if err := validateAdInput(input); err != nil {
		return nil, fmt.Errorf("validateAdInput: %w", err)
	}

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

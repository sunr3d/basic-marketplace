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
	if err := validateAdvFilter(filter); err != nil {
		return nil, fmt.Errorf("validateAdvFilter: %w", err)
	}

	ads, err := s.AdvRepo.FindMany(filter)
	if err != nil {
		return nil, fmt.Errorf("FindMany: %w", err)
	}

	var feed []interfaces.AdvFeedItem
	for _, adv := range ads {
		feed = append(feed, interfaces.AdvFeedItem{
			AdvBase: interfaces.AdvBase{
				Title:       adv.Title,
				Description: adv.Description,
				ImageURL:    adv.ImageURL,
				Price:       adv.Price,
			},
			ID:         adv.ID,
			OwnerLogin: adv.OwnerLogin,
			IsOwner:    adv.OwnerID == currentUserID,
			CreatedAt:  adv.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return feed, nil
}

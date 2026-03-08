package wishlist

import (
	"errors"
	"kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
)

type Service interface {
	AddWishlist(w *entities.Wishlist) error
	GetWishlists(params common.ListParams, year string) ([]entities.Wishlist, int64, error)
	GetWishlistByID(id string) (*entities.Wishlist, error)
	UpdateWishlist(w *entities.Wishlist, updates map[string]interface{}) error
	DeleteWishlist(id string) error
}

type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{repository: r}
}

func (s *service) AddWishlist(w *entities.Wishlist) error {
	if w.Status == "" {
		w.Status = "planned"
	}

	if !isValidStatus(w.Status) {
		return errors.New("invalid status: must be 'planned' or 'acquired'")
	}

	if w.Position == 0 {
		maxPos, err := s.repository.GetMaxPosition()
		if err != nil {
			return err
		}
		w.Position = maxPos + 1
	}

	return s.repository.CreateWishlist(w)
}

func (s *service) GetWishlists(params common.ListParams, year string) ([]entities.Wishlist, int64, error) {
	return s.repository.ListWishlists(params, year)
}

func (s *service) GetWishlistByID(id string) (*entities.Wishlist, error) {
	return s.repository.GetWishlistByID(id)
}

func (s *service) UpdateWishlist(w *entities.Wishlist, updates map[string]interface{}) error {
	// If status is being updated, validate it
	if status, ok := updates["status"].(string); ok {
		if !isValidStatus(status) {
			return errors.New("invalid status: must be 'planned' or 'acquired'")
		}
	}

	// Ensure the item exists before updating
	if _, err := s.repository.GetWishlistByID(w.ID); err != nil {
		return err
	}
	return s.repository.UpdateWishlist(w, updates)
}

func isValidStatus(status string) bool {
	return status == "planned" || status == "acquired"
}

func (s *service) DeleteWishlist(id string) error {
	// Ensure the item exists before deleting
	if _, err := s.repository.GetWishlistByID(id); err != nil {
		return err
	}
	return s.repository.DeleteWishlist(id)
}

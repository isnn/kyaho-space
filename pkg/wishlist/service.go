package wishlist

import (
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
	return s.repository.CreateWishlist(w)
}

func (s *service) GetWishlists(params common.ListParams, year string) ([]entities.Wishlist, int64, error) {
	return s.repository.ListWishlists(params, year)
}

func (s *service) GetWishlistByID(id string) (*entities.Wishlist, error) {
	return s.repository.GetWishlistByID(id)
}

func (s *service) UpdateWishlist(w *entities.Wishlist, updates map[string]interface{}) error {
	return s.repository.UpdateWishlist(w, updates)
}

func (s *service) DeleteWishlist(id string) error {
	return s.repository.DeleteWishlist(id)
}

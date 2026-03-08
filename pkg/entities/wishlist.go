package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wishlist struct {
	ID           string `json:"id" gorm:"type:uuid;primaryKey"`
	TargetAction string `json:"target_action"`
	TargetYear   int    `json:"target_year"`
	Position     int    `json:"position"`
	Status       string `json:"status"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Wishlist) TableName() string {
	return "wishlists"
}

func (w *Wishlist) BeforeCreate(tx *gorm.DB) error {
	if w.ID == "" {
		w.ID = uuid.New().String()
	}
	return nil
}

type WishlistRequest struct {
	TargetAction string `json:"target_action" validate:"required,min=1"`
	TargetYear   int    `json:"target_year" validate:"omitempty,min=2000"`
	Position     int    `json:"position" validate:"omitempty,min=1"`
	Status       string `json:"status" validate:"omitempty,oneof=planned acquired"`
}

type WishlistUpdateRequest struct {
	TargetAction *string `json:"target_action" validate:"omitempty,min=1"`
	TargetYear   *int    `json:"target_year" validate:"omitempty,min=2000"`
	Position     *int    `json:"position" validate:"omitempty,min=1"`
	Status       *string `json:"status" validate:"omitempty,oneof=planned acquired"`
}

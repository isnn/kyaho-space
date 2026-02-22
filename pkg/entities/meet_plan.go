package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MeetPlan struct {
	ID           string `json:"id" gorm:"type:uuid;primaryKey"`
	Organization string `json:"organization"`
	Email        string `json:"email"`
	InquiryType  string `json:"inquiry_type"`
	Brief        string `json:"brief"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (MeetPlan) TableName() string {
	return "meet_plans"
}

func (m *MeetPlan) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

type MeetPlanRequest struct {
	Organization string `json:"organization" validate:"required,min=1,max=255"`
	Email        string `json:"email" validate:"required,email"`
	InquiryType  string `json:"inquiry_type" validate:"required"`
	Brief        string `json:"brief" validate:"required,min=10"`
}

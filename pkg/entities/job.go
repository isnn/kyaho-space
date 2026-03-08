package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	UserID      string    `json:"user_id"`
	CompanyName string    `json:"company_name"`
	JobTitle    string    `json:"job_title"`
	JobPostURL  string    `json:"job_post_url"`
	Status      string    `json:"status"`
	AppliedDate time.Time `json:"applied_date"`
	Notes       string    `json:"notes"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Job) TableName() string {
	return "jobs"
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}

type JobRequest struct {
	CompanyName string `json:"company_name" validate:"required,min=1"`
	JobTitle    string `json:"job_title" validate:"required,min=1"`
	JobPostURL  string `json:"job_post_url" validate:"omitempty"`
	Status      string `json:"status" validate:"omitempty,oneof=applied technical interview offer rejected ghosted"`
	Notes       string `json:"notes" validate:"omitempty"`
}

type JobUpdateRequest struct {
	CompanyName *string `json:"company_name" validate:"omitempty,min=1"`
	JobTitle    *string `json:"job_title" validate:"omitempty,min=1"`
	JobPostURL  *string `json:"job_post_url" validate:"omitempty"`
	Status      *string `json:"status" validate:"omitempty,oneof=applied technical interview offer rejected ghosted"`
	Notes       *string `json:"notes" validate:"omitempty"`
}

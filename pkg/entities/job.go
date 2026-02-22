package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Job struct {
	ID                   string    `json:"id" gorm:"type:uuid;primaryKey"`
	UserID               string    `json:"user_id"`
	CompanyName          string    `json:"company_name"`
	JobTitle             string    `json:"job_title"`
	JobPostURL           string    `json:"job_post_url"`
	RecruitmentPortalURL string    `json:"recruitment_portal_url"`
	Status               string    `json:"status"`
	AppliedDate          time.Time `json:"applied_date"`
	LastUpdated          time.Time `json:"last_updated"`
	Notes                string    `json:"notes"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName sets the table name for GORM
func (Job) TableName() string {
	return "jobs"
}

// BeforeCreate generates a UUID before inserting into the database
func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}

type JobRequest struct {
	CompanyName string `json:"company_name"`
	JobTitle    string `json:"job_title"`
}

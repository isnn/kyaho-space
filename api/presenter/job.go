package presenter

import "time"

type JobResponse struct {
	ID                   string    `json:"id"`
	UserID               string    `json:"user_id"`
	CompanyName          string    `json:"company_name"`
	JobPostURL           string    `json:"job_post_url"`
	RecruitmentPortalURL string    `json:"recruitment_portal_url"`
	Status               string    `json:"status"`
	AppliedDate          time.Time `json:"applied_date"`
	LastUpdated          time.Time `json:"last_updated"`
	Notes                string    `json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

package job

import (
	"kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"time"
)

type Service interface {
	AddJob(job *entities.Job) error
	GetJobs(params common.ListParams, userID string) ([]entities.Job, int64, error)
	GetJobByID(id string, userID string) (*entities.Job, error)
	UpdateJob(job *entities.Job, updates map[string]interface{}, userID string) error
	DeleteJob(id string, userID string) error
	GetStatistics(userID string) (map[string]interface{}, error)
}

type service struct {
	repository Repository
}

// NewService creates a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) AddJob(job *entities.Job) error {
	if job.Status == "" {
		job.Status = "applied"
	}
	if job.AppliedDate.IsZero() {
		job.AppliedDate = time.Now()
	}
	return s.repository.CreateJob(job)
}

func (s *service) GetJobs(params common.ListParams, userID string) ([]entities.Job, int64, error) {
	return s.repository.ListJobs(params, userID)
}

func (s *service) GetJobByID(id string, userID string) (*entities.Job, error) {
	return s.repository.GetJobByID(id, userID)
}

func (s *service) UpdateJob(job *entities.Job, updates map[string]interface{}, userID string) error {
	return s.repository.UpdateJob(job, updates, userID)
}

func (s *service) DeleteJob(id string, userID string) error {
	return s.repository.DeleteJob(id, userID)
}

func (s *service) GetStatistics(userID string) (map[string]interface{}, error) {
	return s.repository.GetStatistics(userID)
}

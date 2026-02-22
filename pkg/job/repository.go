package job

import (
	"kyaho-space/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	CreateJob(job *entities.Job) error
	ReadJob() ([]entities.Job, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepo creates a new job repository
// Each package gets its own NewRepo — pass the same *gorm.DB to all of them
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateJob inserts a new job record into the database
func (r *repository) CreateJob(job *entities.Job) error {
	return r.db.Create(job).Error
}

// ReadJob fetches all job records from the database
func (r *repository) ReadJob() ([]entities.Job, error) {
	var jobs []entities.Job
	err := r.db.Find(&jobs).Error
	return jobs, err
}

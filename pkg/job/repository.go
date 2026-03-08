package job

import (
	"kyaho-space/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	CreateJob(job *entities.Job) error
	ListJobs() ([]entities.Job, error)
	GetJobByID(id string) (*entities.Job, error)
	UpdateJob(job *entities.Job, updates map[string]interface{}) error
	DeleteJob(id string) error
	GetStatistics() (map[string]interface{}, error)
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

// ListJobs fetches all job records from the database
func (r *repository) ListJobs() ([]entities.Job, error) {
	var jobs []entities.Job
	err := r.db.Order("applied_date desc, created_at desc").Find(&jobs).Error
	return jobs, err
}

func (r *repository) GetJobByID(id string) (*entities.Job, error) {
	var job entities.Job
	err := r.db.Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *repository) UpdateJob(job *entities.Job, updates map[string]interface{}) error {
	result := r.db.Model(job).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) DeleteJob(id string) error {
	result := r.db.Delete(&entities.Job{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) GetStatistics() (map[string]interface{}, error) {
	var total int64
	var interview int64
	var ghosted int64
	var offer int64

	model := r.db.Model(&entities.Job{})

	model.Count(&total)

	r.db.Model(&entities.Job{}).Where("status = ?", "Interview").Count(&interview)
	r.db.Model(&entities.Job{}).Where("status = ?", "Ghosted").Count(&ghosted)
	r.db.Model(&entities.Job{}).Where("status = ?", "Offer").Count(&offer)

	conversionRate := 0.0
	if total > 0 {
		conversionRate = float64(offer) / float64(total) * 100
	}

	return map[string]interface{}{
		"total_applied": total,
		"interview":     interview,
		"ghosted":       ghosted,
		"offer":         offer,
		"conversion":    conversionRate,
	}, nil
}

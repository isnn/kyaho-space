package job

import (
	"kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"

	"gorm.io/gorm"
)

// allowedSortColumns defines which columns can be sorted for jobs
var allowedSortColumns = map[string]bool{
	"company_name": true,
	"job_title":    true,
	"status":       true,
	"applied_date": true,
	"created_at":   true,
}

type Repository interface {
	CreateJob(job *entities.Job) error
	ListJobs(params common.ListParams, userID string) ([]entities.Job, int64, error)
	GetJobByID(id string, userID string) (*entities.Job, error)
	UpdateJob(job *entities.Job, updates map[string]interface{}, userID string) error
	DeleteJob(id string, userID string) error
	GetStatistics(userID string) (map[string]interface{}, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateJob(job *entities.Job) error {
	return r.db.Create(job).Error
}

func (r *repository) ListJobs(params common.ListParams, userID string) ([]entities.Job, int64, error) {
	params.Sanitize(allowedSortColumns, "applied_date")

	query := r.db.Model(&entities.Job{}).Where("user_id = ?", userID)

	// Search: ILIKE on company_name or job_title
	if params.Search != "" {
		query = query.Where("company_name ILIKE ? OR job_title ILIKE ?", "%"+params.Search+"%", "%"+params.Search+"%")
	}

	// Count total (before pagination)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort + Paginate
	var jobs []entities.Job
	err := query.Order(params.OrderClause()).Limit(params.Limit).Offset(params.Offset()).Find(&jobs).Error
	return jobs, total, err
}

func (r *repository) GetJobByID(id string, userID string) (*entities.Job, error) {
	var job entities.Job
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *repository) UpdateJob(job *entities.Job, updates map[string]interface{}, userID string) error {
	result := r.db.Model(job).Where("user_id = ?", userID).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) DeleteJob(id string, userID string) error {
	result := r.db.Delete(&entities.Job{}, "id = ? AND user_id = ?", id, userID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) GetStatistics(userID string) (map[string]interface{}, error) {
	var total int64
	var interview int64
	var ghosted int64
	var offer int64

	model := r.db.Model(&entities.Job{}).Where("user_id = ?", userID)

	model.Count(&total)

	r.db.Model(&entities.Job{}).Where("user_id = ? AND status IN ?", userID, []string{"interview", "Interview"}).Count(&interview)
	r.db.Model(&entities.Job{}).Where("user_id = ? AND status IN ?", userID, []string{"ghosted", "Ghosted"}).Count(&ghosted)
	r.db.Model(&entities.Job{}).Where("user_id = ? AND status IN ?", userID, []string{"offer", "Offer"}).Count(&offer)

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

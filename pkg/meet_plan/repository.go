package meet_plan

import (
	"kyaho-space/pkg/entities"

	"gorm.io/gorm"
)

type Repository interface {
	CreateMeetPlan(meetPlan *entities.MeetPlan) error
	ListMeetPlans() ([]entities.MeetPlan, error)
	DeleteMeetPlan(id string) error
}

type repository struct {
	db *gorm.DB
}

// NewRepo creates a new job repository
// Each package gets its own NewRepo — pass the same *gorm.DB to all of them
func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

// CreateMeetPlan inserts a new meet plan record into the database
func (r *repository) CreateMeetPlan(meetPlan *entities.MeetPlan) error {
	return r.db.Create(meetPlan).Error
}

// ListMeetPlans fetches meet plan records from the database
func (r *repository) ListMeetPlans() ([]entities.MeetPlan, error) {
	var meetPlans []entities.MeetPlan
	err := r.db.Order("created_at desc").Limit(100).Find(&meetPlans).Error
	return meetPlans, err
}

func (r *repository) DeleteMeetPlan(id string) error {
	result := r.db.Delete(&entities.MeetPlan{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

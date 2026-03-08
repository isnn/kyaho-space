package wishlist

import (
	"kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"strconv"

	"gorm.io/gorm"
)

// allowedSortColumns defines which columns can be sorted for wishlists
var allowedSortColumns = map[string]bool{
	"target_action": true,
	"target_year":   true,
	"created_at":    true,
}

type Repository interface {
	CreateWishlist(w *entities.Wishlist) error
	ListWishlists(params common.ListParams, year string) ([]entities.Wishlist, int64, error)
	GetWishlistByID(id string) (*entities.Wishlist, error)
	UpdateWishlist(w *entities.Wishlist, updates map[string]interface{}) error
	DeleteWishlist(id string) error
	GetMaxPosition() (int, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepo(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateWishlist(w *entities.Wishlist) error {
	return r.db.Create(w).Error
}

func (r *repository) GetWishlistByID(id string) (*entities.Wishlist, error) {
	var w entities.Wishlist
	err := r.db.Where("id = ?", id).First(&w).Error
	if err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *repository) ListWishlists(params common.ListParams, year string) ([]entities.Wishlist, int64, error) {
	params.Sanitize(allowedSortColumns, "created_at")

	query := r.db.Model(&entities.Wishlist{})

	// Search: ILIKE on "target_action"
	if params.Search != "" {
		query = query.Where("target_action ILIKE ?", "%"+params.Search+"%")
	}

	// Filter: by target_year (only when non-empty)
	if year != "" {
		y, err := strconv.Atoi(year)
		if err == nil && y > 0 {
			query = query.Where("target_year = ?", y)
		}
	}

	// Count total (before pagination)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort + Paginate
	var wishlists []entities.Wishlist
	err := query.Order(params.OrderClause()).Limit(params.Limit).Offset(params.Offset()).Find(&wishlists).Error
	return wishlists, total, err
}

func (r *repository) UpdateWishlist(w *entities.Wishlist, updates map[string]interface{}) error {
	result := r.db.Model(w).Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) DeleteWishlist(id string) error {
	result := r.db.Delete(&entities.Wishlist{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *repository) GetMaxPosition() (int, error) {
	var max int
	err := r.db.Model(&entities.Wishlist{}).Select("COALESCE(MAX(position), 0)").Scan(&max).Error
	return max, err
}

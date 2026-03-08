package wishlist

import (
	"errors"
	"kyaho-space/pkg/common"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/wishlist/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddWishlist(t *testing.T) {
	tests := []struct {
		name          string
		input         *entities.Wishlist
		mockSetup     func(*mocks.MockWishlistRepository)
		expectedError error
	}{
		{
			name: "Success - With Default Values Setup",
			input: &entities.Wishlist{
				TargetAction: "Buy Laptop",
				TargetYear:   2027,
			},
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetMaxPosition").Return(5, nil).Once()
				
				// Expected to hit create with Status="planned" and Position=6
				m.On("CreateWishlist", mock.MatchedBy(func(w *entities.Wishlist) bool {
					return w.Status == "planned" && w.Position == 6
				})).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Error - Failing to Fetch Max Position",
			input: &entities.Wishlist{
				TargetAction: "Buy Phone",
			},
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetMaxPosition").Return(0, errors.New("db error on aggregate")).Once()
			},
			expectedError: errors.New("db error on aggregate"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockWishlistRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			err := service.AddWishlist(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetWishlists(t *testing.T) {
	mockData := []entities.Wishlist{
		{TargetAction: "Car", TargetYear: 2030},
	}
	mockParams := common.ListParams{Page: 1, Limit: 10}

	tests := []struct {
		name          string
		mockSetup     func(*mocks.MockWishlistRepository)
		expectedData  []entities.Wishlist
		expectedTotal int64
		expectedError error
	}{
		{
			name: "Success - Retrieve Wishlist Items",
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("ListWishlists", mockParams, "2030").Return(mockData, int64(1), nil).Once()
			},
			expectedData:  mockData,
			expectedTotal: 1,
			expectedError: nil,
		},
		{
			name: "Error - Database Failure",
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("ListWishlists", mockParams, "2030").Return(nil, int64(0), errors.New("db fetch failed")).Once()
			},
			expectedData:  nil,
			expectedTotal: 0,
			expectedError: errors.New("db fetch failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockWishlistRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, total, err := service.GetWishlists(mockParams, "2030")

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
				assert.Equal(t, int64(0), total)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, result)
				assert.Equal(t, tt.expectedTotal, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetWishlistByID(t *testing.T) {
	mockData := &entities.Wishlist{ID: "w-1", TargetAction: "Car"}

	tests := []struct {
		name          string
		input         string
		mockSetup     func(*mocks.MockWishlistRepository)
		expectedData  *entities.Wishlist
		expectedError error
	}{
		{
			name:  "Success - Find Wishlist",
			input: "w-1",
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-1").Return(mockData, nil).Once()
			},
			expectedData:  mockData,
			expectedError: nil,
		},
		{
			name:  "Error - Not Found",
			input: "w-99",
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-99").Return(nil, errors.New("not found")).Once()
			},
			expectedData:  nil,
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockWishlistRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.GetWishlistByID(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedData, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateWishlist(t *testing.T) {
	tests := []struct {
		name          string
		inputWishlist *entities.Wishlist
		inputUpdates  map[string]interface{}
		mockSetup     func(*mocks.MockWishlistRepository)
		expectedError error
	}{
		{
			name:          "Success - Update Wishlist",
			inputWishlist: &entities.Wishlist{ID: "w-1"},
			inputUpdates:  map[string]interface{}{"status": "acquired"},
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-1").Return(&entities.Wishlist{ID: "w-1"}, nil).Once()
				m.On("UpdateWishlist", &entities.Wishlist{ID: "w-1"}, map[string]interface{}{"status": "acquired"}).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:          "Error - Update Failed",
			inputWishlist: &entities.Wishlist{ID: "w-2"},
			inputUpdates:  map[string]interface{}{"status": "acquired"},
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-2").Return(&entities.Wishlist{ID: "w-2"}, nil).Once()
				m.On("UpdateWishlist", &entities.Wishlist{ID: "w-2"}, map[string]interface{}{"status": "acquired"}).Return(errors.New("db error")).Once()
			},
			expectedError: errors.New("db error"),
		},
		{
			name:          "Error - Item Not Found",
			inputWishlist: &entities.Wishlist{ID: "w-unknown"},
			inputUpdates:  map[string]interface{}{"status": "acquired"},
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-unknown").Return(nil, errors.New("record not found")).Once()
			},
			expectedError: errors.New("record not found"),
		},
		{
			name:          "Error - Invalid Status",
			inputWishlist: &entities.Wishlist{ID: "w-1"},
			inputUpdates:  map[string]interface{}{"status": "invalid-status"},
			mockSetup:     func(m *mocks.MockWishlistRepository) {}, // Should fail before calling DB
			expectedError: errors.New("invalid status: must be 'planned' or 'acquired'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockWishlistRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			err := service.UpdateWishlist(tt.inputWishlist, tt.inputUpdates)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestDeleteWishlist(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		mockSetup     func(*mocks.MockWishlistRepository)
		expectedError error
	}{
		{
			name:  "Success - Delete Wishlist",
			input: "w-1",
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-1").Return(&entities.Wishlist{ID: "w-1"}, nil).Once()
				m.On("DeleteWishlist", "w-1").Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name:  "Error - Delete Failed",
			input: "w-99",
			mockSetup: func(m *mocks.MockWishlistRepository) {
				m.On("GetWishlistByID", "w-99").Return(nil, errors.New("delete error")).Once()
			},
			expectedError: errors.New("delete error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockWishlistRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			err := service.DeleteWishlist(tt.input)

			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

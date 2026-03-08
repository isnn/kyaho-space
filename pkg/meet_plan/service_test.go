package meet_plan

import (
	"errors"
	"kyaho-space/pkg/entities"
	"kyaho-space/pkg/meet_plan/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAddMeetPlan(t *testing.T) {
	tests := []struct {
		name          string
		input         *entities.MeetPlan
		mockSetup     func(*mocks.MockMeetPlanRepository)
		expectedError error
	}{
		{
			name: "Success - MeetPlan Added",
			input: &entities.MeetPlan{
				Organization: "Tech Corp",
				Email:        "contact@techcorp.com",
				InquiryType:  "business",
				Brief:        "Discuss partnership",
			},
			mockSetup: func(m *mocks.MockMeetPlanRepository) {
				m.On("CreateMeetPlan", mock.AnythingOfType("*entities.MeetPlan")).Return(nil).Once()
			},
			expectedError: nil,
		},
		{
			name: "Error - Database Failure",
			input: &entities.MeetPlan{
				Organization: "Tech Corp",
				Email:        "contact@techcorp.com",
			},
			mockSetup: func(m *mocks.MockMeetPlanRepository) {
				m.On("CreateMeetPlan", mock.AnythingOfType("*entities.MeetPlan")).Return(errors.New("db insert failed")).Once()
			},
			expectedError: errors.New("db insert failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockMeetPlanRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			err := service.AddMeetPlan(tt.input)

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

func TestGetMeetPlans(t *testing.T) {
	mockData := []entities.MeetPlan{
		{Organization: "Org A", Email: "orga@test.com"},
		{Organization: "Org B", Email: "orgb@test.com"},
	}

	tests := []struct {
		name          string
		mockSetup     func(*mocks.MockMeetPlanRepository)
		expectedData  []entities.MeetPlan
		expectedError error
	}{
		{
			name: "Success - Retrieve MeetPlans",
			mockSetup: func(m *mocks.MockMeetPlanRepository) {
				m.On("ListMeetPlans").Return(mockData, nil).Once()
			},
			expectedData:  mockData,
			expectedError: nil,
		},
		{
			name: "Error - Database Failure",
			mockSetup: func(m *mocks.MockMeetPlanRepository) {
				m.On("ListMeetPlans").Return(nil, errors.New("db fetch failed")).Once()
			},
			expectedData:  nil,
			expectedError: errors.New("db fetch failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mocks.MockMeetPlanRepository)
			tt.mockSetup(mockRepo)

			service := NewService(mockRepo)
			result, err := service.GetMeetPlans()

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

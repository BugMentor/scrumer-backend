package graph_test

import (
	"fmt"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	. "scrumer-backend/graph"
	"scrumer-backend/graph/mocks"
	"scrumer-backend/models"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB := mocks.NewMockGormDB(ctrl)
	resolver := &Resolver{DB: mockDB}

	// Test case 1: Successful user creation
	t.Run("success", func(t *testing.T) {
		args := graphql.ResolveParams{
			Args: map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "password123",
			},
		}

		expectedUser := models.User{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "password123",
		}

		// Expect the Create method to be called with a pointer to a models.User
		// and return a mock *gorm.DB with no error.
		mockDB.EXPECT().Create(gomock.Any()).DoAndReturn(func(user *models.User) *gorm.DB {
			*user = expectedUser // Simulate GORM populating the user object with ID, etc.
			return &gorm.DB{}    // Return a dummy *gorm.DB
		}).Times(1)

		result, err := resolver.CreateUser(args)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		createdUser, ok := result.(models.User)
		assert.True(t, ok)
		assert.Equal(t, expectedUser.Username, createdUser.Username)
		assert.Equal(t, expectedUser.Email, createdUser.Email)
		assert.Equal(t, expectedUser.Password, createdUser.Password) // In a real app, you'd test the hashed password

	})

	// Test case 2: Failed user creation (e.g., database error)
	t.Run("failure_db_error", func(t *testing.T) {
		args := graphql.ResolveParams{
			Args: map[string]interface{}{
				"username": "failuser",
				"email":    "fail@example.com",
				"password": "password123",
			},
		}

		// Simulate a database error (gorm.DB with Error set; do not use AddError on zero value)
		mockDB.EXPECT().Create(gomock.Any()).Return(&gorm.DB{Error: fmt.Errorf("database error")}).Times(1)

		result, err := resolver.CreateUser(args)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user")
		assert.Nil(t, result)
	})
}
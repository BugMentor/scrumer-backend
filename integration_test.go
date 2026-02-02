package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"scrumer-backend/models"
)

func TestIntegrationPing(t *testing.T) {
	db := connectDB(t)
	router, err := setupRouter(db)
	require.NoError(t, err)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))
	assert.Equal(t, "pong", body["message"])
}

func TestIntegrationGraphQLHello(t *testing.T) {
	db := connectDB(t)
	router, err := setupRouter(db)
	require.NoError(t, err)

	body := map[string]interface{}{
		"query": `query { hello }`,
	}
	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp struct {
		Data struct {
			Hello string `json:"hello"`
		} `json:"data"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, "world", resp.Data.Hello)
}

func TestIntegrationGraphQLCreateUser(t *testing.T) {
	db := connectDB(t)
	router, err := setupRouter(db)
	require.NoError(t, err)

	body := map[string]interface{}{
		"query": `mutation CreateUser($username: String!, $email: String!, $password: String!) {
			createUser(username: $username, email: $email, password: $password) {
				id
				username
				email
			}
		}`,
		"variables": map[string]string{
			"username": "inttestuser",
			"email":    "inttest@example.com",
			"password": "password123",
		},
	}
	bodyBytes, _ := json.Marshal(body)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/graphql", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	respBytes := w.Body.Bytes()
	var resp struct {
		Data struct {
			CreateUser struct {
				ID       interface{} `json:"id"` // ID can be string or number from GraphQL
				Username string      `json:"username"`
				Email    string      `json:"email"`
			} `json:"createUser"`
		} `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	require.NoError(t, json.Unmarshal(respBytes, &resp))
	if len(resp.Errors) > 0 {
		t.Fatalf("GraphQL errors: %v", resp.Errors)
	}
	require.NotNil(t, resp.Data.CreateUser.ID, "createUser.id should be set (response: %s)", string(respBytes))
	assert.Equal(t, "inttestuser", resp.Data.CreateUser.Username)
	assert.Equal(t, "inttest@example.com", resp.Data.CreateUser.Email)
}

func connectDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("integration test DB: %v", err)
	}
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.Project{}, &models.Sprint{}, &models.Task{}))
	return db
}

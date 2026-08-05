package routes_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"task-management-api/internal/bootstrap"
	"task-management-api/internal/database"
	"task-management-api/internal/middleware"
	"task-management-api/internal/routes"
	"task-management-api/internal/util"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	require.NoError(t, err)

	err = database.AutoMigrate(db)
	require.NoError(t, err)

	jwtSecret := "e2e-test-secret"
	app := bootstrap.New(db, jwtSecret)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.RequestIDMiddleware())
	routes.Register(router, app, jwtSecret)

	return router
}

func TestAPI_FullUserLifecycle_E2E(t *testing.T) {
	router := setupTestRouter(t)

	// Step 1: Register User A
	regBody, _ := json.Marshal(map[string]string{
		"name":     "User A",
		"email":    "userA@e2e.com",
		"password": "password123",
	})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 2: Login User A
	loginBody, _ := json.Marshal(map[string]string{
		"email":    "userA@e2e.com",
		"password": "password123",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(loginBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResp util.Response
	err := json.Unmarshal(w.Body.Bytes(), &loginResp)
	require.NoError(t, err)

	loginDataMap := loginResp.Data.(map[string]interface{})
	tokenA := loginDataMap["access_token"].(string)
	require.NotEmpty(t, tokenA)

	// Step 3: Register User B & Get Token B
	regBodyB, _ := json.Marshal(map[string]string{
		"name":     "User B",
		"email":    "userB@e2e.com",
		"password": "password123",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(regBodyB))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var regRespB util.Response
	_ = json.Unmarshal(w.Body.Bytes(), &regRespB)
	regDataB := regRespB.Data.(map[string]interface{})
	tokenB := regDataB["access_token"].(string)
	require.NotEmpty(t, tokenB)

	// Step 4: User A creates Project A
	projBody, _ := json.Marshal(map[string]string{
		"name":  "Project A",
		"color": "#00FF00",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/projects", bytes.NewBuffer(projBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenA)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var projResp util.Response
	_ = json.Unmarshal(w.Body.Bytes(), &projResp)
	projData := projResp.Data.(map[string]interface{})
	projID := uint(projData["id"].(float64))

	// Step 5: User A creates Task A under Project A
	taskBody, _ := json.Marshal(map[string]interface{}{
		"title":      "Task A Title",
		"project_id": projID,
		"priority":   "high",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(taskBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+tokenA)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var taskResp util.Response
	_ = json.Unmarshal(w.Body.Bytes(), &taskResp)
	taskData := taskResp.Data.(map[string]interface{})
	taskID := uint(taskData["id"].(float64))

	// Step 6: User B attempts to access Task A using Token B -> MUST BE 404 NOT FOUND (IDOR Guard Verified)
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%d", taskID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenB)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	// Step 7: User A accesses Task A -> 200 OK
	w = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/tasks/%d", taskID), nil)
	req.Header.Set("Authorization", "Bearer "+tokenA)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

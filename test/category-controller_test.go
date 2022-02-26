package test

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"golang-restfulapi/app"
	"golang-restfulapi/controller"
	"golang-restfulapi/helper"
	"golang-restfulapi/middleware"
	"golang-restfulapi/model/domain"
	"golang-restfulapi/repository"
	"golang-restfulapi/service"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func setUpTestDB() *sql.DB {
	db, err := sql.Open("mysql", "root:@tcp(localhost:3306)/golangrestfulapi_test")
	helper.PanicIfError(err)

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func setUpRouter(db *sql.DB) http.Handler {
	validate := validator.New()
	categoryRepository := repository.NewCategoryRepository()
	categoryService := service.NewCategoryService(categoryRepository, db, validate)
	categoryController := controller.NewCategoryController(categoryService)

	router := app.NewRouter(categoryController)

	return middleware.NewAuthMiddleware(router)
}

func truncateCategory(db *sql.DB) {
	_, err := db.Exec("TRUNCATE category")
	helper.PanicIfError(err)
}

func TestCreateCategorySuccess(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name": "iPhone 12"}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusOK, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, `{"name": "Gadget"}`, responseBody["data"].(map[string]interface{})["name"])
}

func TestCreateCategoryFailed(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)
	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/api/categories", requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	response := recorder.Result()
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestUpdateCategorySuccess(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	repo := repository.NewCategoryRepository()
	category := repo.Save(context.Background(), tx, domain.Category{Name: "iPhone 13"})
	err := tx.Commit()
	helper.PanicIfError(err)

	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name": "iPhone 12"}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, "iPhone 12", responseBody["data"].(map[string]interface{})["name"])
}

func TestUpdateCategoryFailed(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	repo := repository.NewCategoryRepository()
	category := repo.Save(context.Background(), tx, domain.Category{Name: "iPhone 13"})
	err := tx.Commit()
	helper.PanicIfError(err)

	router := setUpRouter(db)

	requestBody := strings.NewReader(`{"name": ""}`)
	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), requestBody)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusBadRequest, int(responseBody["code"].(float64)))
	assert.Equal(t, "BAD REQUEST", responseBody["status"])
}

func TestGetCategorySuccess(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	repo := repository.NewCategoryRepository()
	category := repo.Save(context.Background(), tx, domain.Category{Name: "iPhone 13"})
	err := tx.Commit()
	helper.PanicIfError(err)

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
	assert.Equal(t, category.Id, int(responseBody["data"].(map[string]interface{})["id"].(float64)))
	assert.Equal(t, category.Name, responseBody["data"].(map[string]interface{})["name"])
}

func TestGetCategoryFailed(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestDeleteCategorySuccess(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	repo := repository.NewCategoryRepository()
	category := repo.Save(context.Background(), tx, domain.Category{Name: "iPhone 13"})
	err := tx.Commit()
	helper.PanicIfError(err)

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/"+strconv.Itoa(category.Id), nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])
}

func TestDeleteCategoryFailed(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodDelete, "http://localhost:3000/api/categories/404", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusNotFound, int(responseBody["code"].(float64)))
	assert.Equal(t, "NOT FOUND", responseBody["status"])
}

func TestGetListCategoriesSuccess(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	tx, _ := db.Begin()
	repo := repository.NewCategoryRepository()

	category := repo.Save(context.Background(), tx, domain.Category{Name: "iPhone 13 Pro"})
	category2 := repo.Save(context.Background(), tx, domain.Category{Name: "MacBook Pro"})

	err := tx.Commit()
	helper.PanicIfError(err)

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "RAHASIA")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err = json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusOK, int(responseBody["code"].(float64)))
	assert.Equal(t, "OK", responseBody["status"])

	var categories = responseBody["data"].([]map[string]interface{})

	assert.Equal(t, category.Id, categories[0]["id"])
	assert.Equal(t, category.Name, categories[0]["name"])

	assert.Equal(t, category2.Id, categories[1]["id"])
	assert.Equal(t, category2.Name, categories[1]["name"])
}

func TestUnauthorized(t *testing.T) {
	db := setUpTestDB()
	truncateCategory(db)

	router := setUpRouter(db)

	request := httptest.NewRequest(http.MethodPut, "http://localhost:3000/api/categories", nil)
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("X-API-KEY", "SALAH")

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	response := recorder.Result()

	body, _ := io.ReadAll(response.Body)
	var responseBody map[string]interface{}
	err := json.Unmarshal(body, &responseBody)
	helper.PanicIfError(err)

	assert.Equal(t, http.StatusUnauthorized, int(responseBody["code"].(float64)))
	assert.Equal(t, "UNAUTHORIZED", responseBody["status"])
}

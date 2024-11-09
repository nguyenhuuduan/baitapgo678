package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"go_library/database"
	"go_library/models"
)

func setupTestDB() {
	// Sử dụng SQLite in-memory database mà không cần CGO
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.Book{})
	database.DB = db
}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/books", CreateBook)
	r.PUT("/books/:id", UpdateBook)
	r.DELETE("/books/:id", DeleteBook)
	r.GET("/books", GetBooks)
	r.GET("/books/search", SearchBooksByTitle)
	return r
}
func TestCreateBook(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Tạo dữ liệu sách
	book := models.Book{
		Title:       "Test Book",
		Author:      "Test Author",
		Description: "Test Description",
	}
	bookJSON, _ := json.Marshal(book)

	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(bookJSON))
	req.Header.Set("Content-Type", "application/json")

	// Gửi yêu cầu POST để thêm sách
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra phản hồi
	assert.Equal(t, http.StatusCreated, w.Code)
	var responseBook models.Book
	json.Unmarshal(w.Body.Bytes(), &responseBook)
	assert.Equal(t, book.Title, responseBook.Title)
	assert.Equal(t, book.Author, responseBook.Author)
}

func TestUpdateBook(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Thêm sách ban đầu vào database
	book := models.Book{Title: "Old Title", Author: "Old Author"}
	database.DB.Create(&book)

	// Tạo dữ liệu cập nhật
	updatedBook := models.Book{Title: "New Title", Author: "New Author"}
	updatedBookJSON, _ := json.Marshal(updatedBook)

	// Gửi yêu cầu PUT
	req, _ := http.NewRequest("PUT", "/books/"+strconv.Itoa(int(book.ID)), bytes.NewBuffer(updatedBookJSON))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra phản hồi
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBook models.Book
	json.Unmarshal(w.Body.Bytes(), &responseBook)
	assert.Equal(t, "New Title", responseBook.Title)
	assert.Equal(t, "New Author", responseBook.Author)
}

func TestDeleteBook(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Thêm sách ban đầu vào database
	book := models.Book{Title: "Delete Me"}
	database.DB.Create(&book)

	// Gửi yêu cầu DELETE
	req, _ := http.NewRequest("DELETE", "/books/"+strconv.Itoa(int(book.ID)), nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra phản hồi
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]string
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Equal(t, "Book deleted", response["message"])

	// Kiểm tra sách có thực sự bị xóa không
	var deletedBook models.Book
	err := database.DB.First(&deletedBook, book.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
}

func TestSearchBooksByTitle_Found(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	// Thêm sách vào database
	book := models.Book{Title: "Find Me"}
	database.DB.Create(&book)

	req, _ := http.NewRequest("GET", "/books/search?title=Find", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra phản hồi
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBooks []models.Book
	json.Unmarshal(w.Body.Bytes(), &responseBooks)
	assert.Len(t, responseBooks, 1)
	assert.Equal(t, "Find Me", responseBooks[0].Title)
}

func TestSearchBooksByTitle_NotFound(t *testing.T) {
	setupTestDB()
	router := setupRouter()

	req, _ := http.NewRequest("GET", "/books/search?title=Unknown", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Kiểm tra phản hồi
	assert.Equal(t, http.StatusOK, w.Code)
	var responseBooks []models.Book
	json.Unmarshal(w.Body.Bytes(), &responseBooks)
	assert.Len(t, responseBooks, 0)
}

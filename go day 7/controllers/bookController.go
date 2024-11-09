package controllers

import (
	"go_library/database"
	"go_library/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Thêm sách
func CreateBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Create(&book)
	c.JSON(http.StatusCreated, book)
}

// Cập nhật thông tin sách
func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching book"})
		return
	}

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	database.DB.Save(&book)
	c.JSON(http.StatusOK, book)
}

// Xóa sách
func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	if err := database.DB.First(&book, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		return
	}
	database.DB.Delete(&book)
	c.JSON(http.StatusOK, gin.H{"message": "Book deleted"})
}

// Lấy danh sách sách
func GetBooks(c *gin.Context) {
	var books []models.Book
	database.DB.Find(&books)
	c.JSON(http.StatusOK, books)
}

// Hàm tìm kiếm sách theo tiêu đề, sửa lỗi ILIKE cho SQLite
func SearchBooksByTitle(c *gin.Context) {
	title := c.Query("title")
	var books []models.Book

	// Sử dụng LOWER và LIKE để hỗ trợ tìm kiếm không phân biệt hoa/thường
	database.DB.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(title)+"%").Find(&books)

	c.JSON(http.StatusOK, books)
}

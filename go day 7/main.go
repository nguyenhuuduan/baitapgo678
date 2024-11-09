package main

import (
	"go_library/controllers"
	"go_library/database"
	"go_library/models"

	"github.com/gin-gonic/gin"
)

func main() {
	// Khởi động Gin router
	r := gin.Default()

	// Kết nối database
	database.ConnectDB()

	// Tự động migrate model Book để tạo bảng nếu chưa có
	database.DB.AutoMigrate(&models.Book{})

	// Định nghĩa các route
	r.POST("/books", controllers.CreateBook)
	r.PUT("/books/:id", controllers.UpdateBook)
	r.DELETE("/books/:id", controllers.DeleteBook)
	r.GET("/books", controllers.GetBooks)
	r.GET("/books/search", controllers.SearchBooksByTitle)

	// Chạy server trên cổng 8080
	r.Run(":8080")
}

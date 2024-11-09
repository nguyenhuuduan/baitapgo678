package main

import (
	"student_management/config"
	"student_management/controllers"
	"student_management/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// Thêm cấu hình CORS để cho phép các yêu cầu từ mọi nguồn
	r.Use(cors.Default())

	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Student{})

	r.POST("/students", controllers.CreateStudent)
	r.GET("/students", controllers.GetStudents)
	r.GET("/students/:id", controllers.GetStudentByID)
	r.PUT("/students/:id", controllers.UpdateStudent)
	r.DELETE("/students/:id", controllers.DeleteStudent)

	r.Run("0.0.0.0:8080") // Chạy server tại cổng 8080
}

package routes

import (
	"student_service/controllers"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"     // Import swaggerFiles để sử dụng với Gin Swagger
	ginSwagger "github.com/swaggo/gin-swagger" // Import Gin Swagger
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Cấu hình route Swagger để hiển thị tài liệu API
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Các routes CRUD cho sinh viên
	studentRoutes := r.Group("/students")
	{
		studentRoutes.POST("", controllers.CreateStudent)
		studentRoutes.GET("/:id", controllers.GetStudent)
		studentRoutes.PUT("/:id", controllers.UpdateStudent)
		studentRoutes.DELETE("/:id", controllers.DeleteStudent)
	}
	return r
}

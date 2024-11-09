package main

import (
	_ "student_service/docs"
	"student_service/routes"
)

// @title Student Service API
// @version 1.0
// @description API cho hệ thống quản lý sinh viên với Redis cache
// @host localhost:8080
// @BasePath /

func main() {
	r := routes.SetupRouter()
	r.Run(":8080")
}

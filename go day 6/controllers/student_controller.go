package controllers

import (
	"net/http"
	"student_service/cache"
	"student_service/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var students = make(map[string]models.Student)
var redisCache = cache.NewRedisCache("localhost:6379", "", 0)

// @Summary Create a new student
// @Description Create a new student in the system
// @Accept  json
// @Produce  json
// @Param student body models.Student true "Student Information"
// @Success 200 {object} models.Student
// @Router /students [post]
func CreateStudent(c *gin.Context) {
	var student models.Student
	if err := c.ShouldBindJSON(&student); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	student.ID = uuid.New().String()
	students[student.ID] = student

	// Cache the student
	redisCache.Set(student.ID, student.Name, 5*time.Minute)
	c.JSON(http.StatusOK, student)
}

// @Summary Get a student by ID
// @Description Get a student by ID
// @Param id path string true "Student ID"
// @Success 200 {object} models.Student
// @Router /students/{id} [get]
func GetStudent(c *gin.Context) {
	id := c.Param("id")

	if studentName, err := redisCache.Get(id); err == nil {
		c.JSON(http.StatusOK, gin.H{"id": id, "name": studentName})
		return
	}

	if student, exists := students[id]; exists {
		redisCache.Set(id, student.Name, 5*time.Minute)
		c.JSON(http.StatusOK, student)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
}

// @Summary Update a student by ID
// @Description Update a student in the system
// @Accept  json
// @Produce  json
// @Param id path string true "Student ID"
// @Param student body models.Student true "Updated Student Information"
// @Success 200 {object} models.Student
// @Router /students/{id} [put]
func UpdateStudent(c *gin.Context) {
	id := c.Param("id")

	var updatedStudent models.Student
	if err := c.ShouldBindJSON(&updatedStudent); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, exists := students[id]; exists {
		updatedStudent.ID = id
		students[id] = updatedStudent
		redisCache.Set(id, updatedStudent.Name, 5*time.Minute)
		c.JSON(http.StatusOK, updatedStudent)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
}

// @Summary Delete a student by ID
// @Description Delete a student in the system
// @Param id path string true "Student ID"
// @Success 200 {object} string "Student deleted successfully"
// @Router /students/{id} [delete]
func DeleteStudent(c *gin.Context) {
	id := c.Param("id")

	if _, exists := students[id]; exists {
		delete(students, id)
		redisCache.Delete(id)
		c.JSON(http.StatusOK, gin.H{"message": "Student deleted successfully"})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Student not found"})
}

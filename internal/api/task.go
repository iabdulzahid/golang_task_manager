package api

import (
	"net/http"
	"strconv"

	// _ "github.com/iabdulzahid/golang_task_manager/docs" // Import Swagger docs

	"github.com/gin-gonic/gin"
	"github.com/iabdulzahid/golang_task_manager/internal/db"
	"github.com/iabdulzahid/golang_task_manager/internal/models"
	"github.com/iabdulzahid/golang_task_manager/pkg/globals"
)

// CreateTask godoc
// @Summary Create a new task
// @Description Create a new task with title, description, priority, and due date
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task data"
// @Success 201 {object} models.SuccessMessage "Task Created Successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	logger := globals.Logger
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		logger.Info("CreateTask", "err", err.Error())
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Validate task
	if task.Title == "" || task.Priority == "" || task.DueDate == "" {
		// Respond(c, http.StatusBadRequest, "Missing required fields", nil)
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: "Missing required fields"})
		return
	}

	// Set default value for empty labels
	if task.Labels == nil {
		task.Labels = []string{}
	}

	// Create task
	err := db.CreateTask(&task)
	if err != nil {
		// Respond(c, http.StatusBadRequest, "Missing required fields", nil)
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Return the created task
	c.JSON(http.StatusCreated, models.SuccessMessage{"Task Created Succesfully"})
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Description Get a list of all tasks in the system
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks [get]
func GetAllTasks(c *gin.Context) {
	// Call the function to get tasks
	tasks, err := db.GetTasks()
	if err != nil {
		// If there's an error, return 500 with the error message
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: "Failed to fetch tasks: " + err.Error()})
		return
	}

	// Return the list of tasks as a JSON response with status 200
	c.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Description Get task details by task ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	task, err := db.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Update an existing task
// @Description Update task details by task ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task data"
// @Success 200 {object} models.Task
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// Update task
	updatedTask, err := db.UpdateTask(id, &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}

// DeleteTask godoc
// @Summary Delete a task
// @Description Delete a task by its ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {object} models.SuccessMessage
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := db.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, models.SuccessMessage{Message: "Task deleted"})
}

// SendErrorResponse sends an error response with a custom key and error message
func SendResponse(c *gin.Context, statusCode int, messageKey string, message string) {
	// Create a map with dynamic key and message
	c.JSON(statusCode, map[string]string{
		messageKey: message, // Use the dynamic key passed in the function
	})
}

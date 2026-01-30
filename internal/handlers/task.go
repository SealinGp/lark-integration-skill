package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	larktask "github.com/larksuite/oapi-sdk-go/v3/service/task/v1"
	"lark-integration-skill/internal/models"
	"lark-integration-skill/pkg/larkclient"
)

type TaskHandler struct {
	Client *larkclient.ClientWrapper
}

func NewTaskHandler(client *larkclient.ClientWrapper) *TaskHandler {
	return &TaskHandler{Client: client}
}

// CreateTask creates a new task in Lark (using Task V1)
func (h *TaskHandler) CreateTask(c *gin.Context) {
	var req models.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}

	// Build the input for Create Task V1
	taskBody := larktask.NewTaskBuilder().
		Summary(req.Summary).
		Description(req.Description).
		Build()

	// If due time is provided
	if req.DueTime > 0 {
		taskBody.Due = larktask.NewDueBuilder().
			Time(fmt.Sprintf("%d", req.DueTime)).
			Build()
	}

	// NOTE: Task V1 Create does not support adding members directly in the Create call easily
	// via the builder in the same way V2 does. We will just create the task first.

	input := larktask.NewCreateTaskReqBuilder().
		Task(taskBody).
		UserIdType("open_id").
		Build()

	// Call Lark API V1
	resp, err := h.Client.Client.Task.Task.Create(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: fmt.Sprintf("SDK Error: %v", err)})
		return
	}

	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Lark API Error: %d - %s", resp.Code, resp.Msg),
		})
		return
	}

	// Success Response
	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.TaskResponse{
			TaskID:  *resp.Data.Task.Id,
			Summary: *resp.Data.Task.Summary,
			// URL is often not returned in V1 Create response, so we might omit it or construct it
			URL: "",
		},
	})
}

// GetTask retrieves a task
func (h *TaskHandler) GetTask(c *gin.Context) {
	taskID := c.Param("task_id")
	if taskID == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{Status: "error", Message: "Task ID is required"})
		return
	}

	input := larktask.NewGetTaskReqBuilder().
		TaskId(taskID).
		UserIdType("open_id").
		Build()

	resp, err := h.Client.Client.Task.Task.Get(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusNotFound, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Status: "success",
		Data: models.TaskResponse{
			TaskID:  *resp.Data.Task.Id,
			Summary: *resp.Data.Task.Summary,
		},
	})
}

// DeleteTask deletes a task
func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("task_id")
	input := larktask.NewDeleteTaskReqBuilder().TaskId(taskID).Build()

	resp, err := h.Client.Client.Task.Task.Delete(context.Background(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: err.Error()})
		return
	}
	if !resp.Success() {
		c.JSON(http.StatusInternalServerError, models.APIResponse{Status: "error", Message: resp.Msg})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{Status: "success", Message: "Task deleted"})
}

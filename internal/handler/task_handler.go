package handler

import (
	"net/http"
	"strconv"
	"task-manager/internal/domain"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var input struct {
		Title    string `json:"title" binding:"required"`
		Status   string `json:"status" binding:"required"`
		Assignee string `json:"assignee"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &domain.Task{
		Title:    input.Title,
		Status:   input.Status,
		Assignee: input.Assignee,
	}

	if err := h.service.CreateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *TaskHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت شناسه ارسالی نامعتبر است"})
		return
	}

	task, err := h.service.GetTask(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "تسک مورد نظر یافت نشد"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) List(c *gin.Context) {
	status := c.Query("status")
	assignee := c.Query("assignee")
	
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	tasks, err := h.service.ListTasks(c.Request.Context(), status, assignee, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت شناسه ارسالی نامعتبر است"})
		return
	}

	var input struct {
		Title    string `json:"title" binding:"required"`
		Status   string `json:"status" binding:"required"`
		Assignee string `json:"assignee"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := &domain.Task{
		ID:       id,
		Title:    input.Title,
		Status:   input.Status,
		Assignee: input.Assignee,
	}

	if err := h.service.UpdateTask(c.Request.Context(), task); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "تسک یافت نشد یا عملیات با خطا مواجه شد"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "فرمت شناسه ارسالی نامعتبر است"})
		return
	}

	if err := h.service.DeleteTask(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "تسک یافت نشد"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "تسک با موفقیت از سیستم حذف گردید"})
}

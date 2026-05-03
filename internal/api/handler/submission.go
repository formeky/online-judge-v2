package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-judge/internal/model"
	"online-judge/internal/mq"
	"online-judge/internal/repository"
)

type SubmissionHandler struct {
	subRepo     repository.SubmissionRepository
	problemRepo repository.ProblemRepository
	producer    *mq.Producer
}

func NewSubmissionHandler(
	subRepo repository.SubmissionRepository,
	problemRepo repository.ProblemRepository,
	producer *mq.Producer,
) *SubmissionHandler {
	return &SubmissionHandler{subRepo: subRepo, problemRepo: problemRepo, producer: producer}
}

func (h *SubmissionHandler) Submit(c *gin.Context) {
	var req struct {
		ProblemID uint           `json:"problem_id" binding:"required"`
		Language  model.Language `json:"language"   binding:"required"`
		Code      string         `json:"code"       binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := uint(c.GetFloat64("user_id"))

	problem, err := h.problemRepo.GetByID(c.Request.Context(), req.ProblemID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
		return
	}

	sub := &model.Submission{
		ProblemID: req.ProblemID,
		UserID:    userID,
		Language:  req.Language,
		Code:      req.Code,
		Status:    model.StatusPending,
	}
	if err := h.subRepo.Create(c.Request.Context(), sub); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	msg := &mq.JudgeMessage{
		SubmissionID: sub.ID,
		ProblemID:    problem.ID,
		Language:     string(req.Language),
		TimeLimit:    problem.TimeLimit,
		MemoryLimit:  problem.MemoryLimit,
	}
	if err := h.producer.SendJudgeMessage(c.Request.Context(), msg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "enqueue failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
}

func (h *SubmissionHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	problemID, _ := strconv.ParseUint(c.Query("problem_id"), 10, 64)
	userID, _ := strconv.ParseUint(c.Query("user_id"), 10, 64)
	if page < 1 {
		page = 1
	}
	subs, total, err := h.subRepo.List(c.Request.Context(), uint(problemID), uint(userID), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "page": page, "size": size, "data": subs})
}

func (h *SubmissionHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	sub, err := h.subRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, sub)
}

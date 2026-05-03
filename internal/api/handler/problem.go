package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"online-judge/internal/model"
	"online-judge/internal/repository"
)

type ProblemHandler struct {
	problemRepo  repository.ProblemRepository
	testCaseRepo repository.TestCaseRepository
}

func NewProblemHandler(pr repository.ProblemRepository, tr repository.TestCaseRepository) *ProblemHandler {
	return &ProblemHandler{problemRepo: pr, testCaseRepo: tr}
}

func (h *ProblemHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 20
	}
	problems, total, err := h.problemRepo.List(c.Request.Context(), page, size)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total": total, "page": page, "size": size, "data": problems})
}

func (h *ProblemHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	p, err := h.problemRepo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	tcs, _ := h.testCaseRepo.GetByProblemID(c.Request.Context(), p.ID)
	p.TestCases = tcs
	c.JSON(http.StatusOK, p)
}

func (h *ProblemHandler) Create(c *gin.Context) {
	var req struct {
		Title        string            `json:"title"         binding:"required"`
		Description  string            `json:"description"   binding:"required"`
		Difficulty   model.Difficulty  `json:"difficulty"`
		TimeLimit    int64             `json:"time_limit"`
		MemoryLimit  int64             `json:"memory_limit"`
		AllowedLangs string            `json:"allowed_langs"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p := &model.Problem{
		Title:        req.Title,
		Description:  req.Description,
		Difficulty:   req.Difficulty,
		TimeLimit:    req.TimeLimit,
		MemoryLimit:  req.MemoryLimit,
		AllowedLangs: req.AllowedLangs,
	}
	if p.TimeLimit == 0 {
		p.TimeLimit = 2000
	}
	if p.MemoryLimit == 0 {
		p.MemoryLimit = 262144
	}
	if p.AllowedLangs == "" {
		p.AllowedLangs = "c,cpp,java,python,go"
	}
	if err := h.problemRepo.Create(c.Request.Context(), p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func (h *ProblemHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.problemRepo.Update(c.Request.Context(), uint(id), updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "updated"})
}

func (h *ProblemHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := h.problemRepo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *ProblemHandler) AddTestCase(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var req struct {
		Seq        int     `json:"seq"          binding:"required"`
		InputData  string  `json:"input_data"`
		OutputData string  `json:"output_data"`
		InputFile  string  `json:"input_file"`
		OutputFile string  `json:"output_file"`
		Score      int     `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tc := &model.TestCase{
		ProblemID:  uint(id),
		Seq:        req.Seq,
		InputFile:  req.InputFile,
		OutputFile: req.OutputFile,
		Score:      req.Score,
	}
	if req.InputData != "" {
		tc.InputData = &req.InputData
	}
	if req.OutputData != "" {
		tc.OutputData = &req.OutputData
	}
	if tc.Score == 0 {
		tc.Score = 10
	}
	if err := h.testCaseRepo.Create(c.Request.Context(), tc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, tc)
}

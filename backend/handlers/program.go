package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"personal-coach/database"
	"personal-coach/models"
	"personal-coach/services"
)

// ProgramHandler handles workout program endpoints
type ProgramHandler struct {
	claude *services.ClaudeService
	db     *database.DB
}

// NewProgramHandler creates a new program handler
func NewProgramHandler(claude *services.ClaudeService, db *database.DB) *ProgramHandler {
	return &ProgramHandler{claude: claude, db: db}
}

// GenerateProgram handles POST /api/programs/generate
func (h *ProgramHandler) GenerateProgram(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	if req.Person.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person name is required"})
		return
	}
	if req.DaysPerWeek < 1 || req.DaysPerWeek > 7 {
		req.DaysPerWeek = 3
	}
	if req.Weeks < 1 {
		req.Weeks = 4
	}
	if req.Person.ID == "" {
		req.Person.ID = uuid.New().String()
	}

	program, err := h.claude.GenerateProgram(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate program: %v", err)})
		return
	}

	program.ID = uuid.New().String()
	program.PersonID = req.Person.ID
	program.ProfileID = c.GetString("profile_id")
	program.GeneratedAt = time.Now()

	if err := h.db.SaveProgram(*program); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save program: %v", err)})
		return
	}

	c.JSON(http.StatusOK, models.GenerateResponse{
		Program: *program,
		Message: "Programme généré avec succès",
	})
}

// GetProgram handles GET /api/programs/:id
func (h *ProgramHandler) GetProgram(c *gin.Context) {
	id := c.Param("id")

	program, err := h.db.GetProgram(id)
	if errors.Is(err, database.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Programme non trouvé"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}

	c.JSON(http.StatusOK, program)
}

// DownloadPDF handles GET /api/programs/:id/pdf
func (h *ProgramHandler) DownloadPDF(c *gin.Context) {
	id := c.Param("id")

	program, err := h.db.GetProgram(id)
	if errors.Is(err, database.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Programme non trouvé"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}

	pdfBytes, err := services.GeneratePDF(*program)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Erreur PDF: %v", err)})
		return
	}

	filename := fmt.Sprintf("programme-%s-semaine%d.pdf", program.PersonName, program.WeekNumber)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	c.Data(http.StatusOK, "application/pdf", pdfBytes)
}

// GetTimer handles GET /api/programs/:id/timer/:day
func (h *ProgramHandler) GetTimer(c *gin.Context) {
	id := c.Param("id")
	dayStr := c.Param("day")

	dayIndex, err := strconv.Atoi(dayStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid day index"})
		return
	}

	program, err := h.db.GetProgram(id)
	if errors.Is(err, database.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Programme non trouvé"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}

	timer, err := services.BuildTimer(*program, dayIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timer)
}

// ListPrograms handles GET /api/programs
func (h *ProgramHandler) ListPrograms(c *gin.Context) {
	profileID := c.GetString("profile_id")
	programs, err := h.db.ListPrograms(profileID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Database error: %v", err)})
		return
	}
	if programs == nil {
		programs = []models.Program{}
	}

	c.JSON(http.StatusOK, gin.H{"programs": programs, "total": len(programs)})
}

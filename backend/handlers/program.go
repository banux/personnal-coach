package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"personal-coach/models"
	"personal-coach/services"
)

// In-memory store for programs (replace with DB in production)
var (
	programStore = make(map[string]models.Program)
	storeMu      sync.RWMutex
)

// ProgramHandler handles workout program endpoints
type ProgramHandler struct {
	claude *services.ClaudeService
}

// NewProgramHandler creates a new program handler
func NewProgramHandler(claude *services.ClaudeService) *ProgramHandler {
	return &ProgramHandler{claude: claude}
}

// GenerateProgram handles POST /api/programs/generate
func (h *ProgramHandler) GenerateProgram(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request: %v", err)})
		return
	}

	// Validate required fields
	if req.Person.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Person name is required"})
		return
	}
	if req.DaysPerWeek < 1 || req.DaysPerWeek > 7 {
		req.DaysPerWeek = 3 // default
	}
	if req.Weeks < 1 {
		req.Weeks = 4 // default
	}

	// Assign IDs if not provided
	if req.Person.ID == "" {
		req.Person.ID = uuid.New().String()
	}

	// Generate program via Claude
	program, err := h.claude.GenerateProgram(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to generate program: %v", err)})
		return
	}

	// Assign program ID and metadata
	program.ID = uuid.New().String()
	program.PersonID = req.Person.ID
	program.GeneratedAt = time.Now()

	// Store program
	storeMu.Lock()
	programStore[program.ID] = *program
	storeMu.Unlock()

	c.JSON(http.StatusOK, models.GenerateResponse{
		Program: *program,
		Message: "Programme généré avec succès",
	})
}

// GetProgram handles GET /api/programs/:id
func (h *ProgramHandler) GetProgram(c *gin.Context) {
	id := c.Param("id")

	storeMu.RLock()
	program, exists := programStore[id]
	storeMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Programme non trouvé"})
		return
	}

	c.JSON(http.StatusOK, program)
}

// DownloadPDF handles GET /api/programs/:id/pdf
func (h *ProgramHandler) DownloadPDF(c *gin.Context) {
	id := c.Param("id")

	storeMu.RLock()
	program, exists := programStore[id]
	storeMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Programme non trouvé"})
		return
	}

	pdfBytes, err := services.GeneratePDF(program)
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

	storeMu.RLock()
	program, exists := programStore[id]
	storeMu.RUnlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Programme non trouvé"})
		return
	}

	timer, err := services.BuildTimer(program, dayIndex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, timer)
}

// ListPrograms handles GET /api/programs
func (h *ProgramHandler) ListPrograms(c *gin.Context) {
	storeMu.RLock()
	defer storeMu.RUnlock()

	programs := make([]models.Program, 0, len(programStore))
	for _, p := range programStore {
		programs = append(programs, p)
	}

	c.JSON(http.StatusOK, gin.H{"programs": programs, "total": len(programs)})
}

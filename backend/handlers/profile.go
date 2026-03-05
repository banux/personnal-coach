package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"personal-coach/database"
	"personal-coach/models"
)

// ProfileHandler handles profile CRUD endpoints
type ProfileHandler struct {
	db *database.DB
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(db *database.DB) *ProfileHandler {
	return &ProfileHandler{db: db}
}

// List handles GET /api/profiles
func (h *ProfileHandler) List(c *gin.Context) {
	profiles, err := h.db.ListProfiles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur base de données"})
		return
	}
	if profiles == nil {
		profiles = []models.Profile{}
	}
	c.JSON(http.StatusOK, gin.H{"profiles": profiles})
}

// Create handles POST /api/profiles
func (h *ProfileHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Le nom est requis"})
		return
	}

	p := &models.Profile{Name: req.Name}
	if err := h.db.SaveProfile(p); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Un profil avec ce nom existe déjà"})
		return
	}

	c.JSON(http.StatusCreated, p)
}

// Select handles POST /api/profiles/select — stores chosen profile in session
func (h *ProfileHandler) Select(c *gin.Context) {
	var req struct {
		ID string `json:"id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "L'identifiant du profil est requis"})
		return
	}

	profile, err := h.db.GetProfile(req.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profil non trouvé"})
		return
	}

	token := c.GetString("session_token")
	globalSessionStore.SetProfile(token, profile.ID, profile.Name)

	c.JSON(http.StatusOK, gin.H{
		"profile_id":   profile.ID,
		"profile_name": profile.Name,
	})
}

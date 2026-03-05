package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const sessionCookieName = "coach_session"
const sessionDuration = 24 * time.Hour

// sessionData holds per-session state
type sessionData struct {
	Expiry      time.Time
	ProfileID   string
	ProfileName string
}

// SessionStore manages active sessions in memory
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]*sessionData
}

var globalSessionStore = &SessionStore{
	sessions: make(map[string]*sessionData),
}

// generateToken creates a random 32-byte hex token
func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Add stores a new session
func (s *SessionStore) Add(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[token] = &sessionData{Expiry: time.Now().Add(sessionDuration)}
}

// Valid checks if a session token is valid and not expired
func (s *SessionStore) Valid(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	d, ok := s.sessions[token]
	if !ok {
		return false
	}
	return time.Now().Before(d.Expiry)
}

// Delete removes a session
func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
}

// SetProfile associates a profile with a session
func (s *SessionStore) SetProfile(token, profileID, profileName string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if d, ok := s.sessions[token]; ok {
		d.ProfileID = profileID
		d.ProfileName = profileName
	}
}

// GetProfile returns the profile_id and profile_name for a session
func (s *SessionStore) GetProfile(token string) (id, name string) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if d, ok := s.sessions[token]; ok {
		return d.ProfileID, d.ProfileName
	}
	return "", ""
}

// AuthHandler handles authentication endpoints
type AuthHandler struct{}

// NewAuthHandler creates a new auth handler
func NewAuthHandler() *AuthHandler {
	return &AuthHandler{}
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Requête invalide"})
		return
	}

	appPassword := os.Getenv("APP_PASSWORD")
	if appPassword == "" {
		appPassword = "coach2024" // default password if not set
	}

	if req.Password != appPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Mot de passe incorrect"})
		return
	}

	token, err := generateToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur interne"})
		return
	}

	globalSessionStore.Add(token)

	// Set secure cookie
	c.SetCookie(
		sessionCookieName,
		token,
		int(sessionDuration.Seconds()),
		"/",
		"",
		false, // secure: set to true behind HTTPS
		true,  // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{"message": "Connexion réussie"})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	token, err := c.Cookie(sessionCookieName)
	if err == nil {
		globalSessionStore.Delete(token)
	}
	c.SetCookie(sessionCookieName, "", -1, "/", "", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Déconnexion réussie"})
}

// Status handles GET /auth/status - returns auth state + active profile for frontend
func (h *AuthHandler) Status(c *gin.Context) {
	token, err := c.Cookie(sessionCookieName)
	if err != nil || !globalSessionStore.Valid(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}
	profileID, profileName := globalSessionStore.GetProfile(token)
	c.JSON(http.StatusOK, gin.H{
		"authenticated": true,
		"profile_id":    profileID,
		"profile_name":  profileName,
	})
}

// AuthRequired is a Gin middleware that enforces authentication and injects session info into context
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(sessionCookieName)
		if err != nil || !globalSessionStore.Valid(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non authentifié"})
			c.Abort()
			return
		}
		profileID, profileName := globalSessionStore.GetProfile(token)
		c.Set("session_token", token)
		c.Set("profile_id", profileID)
		c.Set("profile_name", profileName)
		c.Next()
	}
}

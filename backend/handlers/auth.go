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

// SessionStore manages active sessions in memory
type SessionStore struct {
	mu       sync.RWMutex
	sessions map[string]time.Time
}

var globalSessionStore = &SessionStore{
	sessions: make(map[string]time.Time),
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
	s.sessions[token] = time.Now().Add(sessionDuration)
}

// Valid checks if a session token is valid and not expired
func (s *SessionStore) Valid(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	exp, ok := s.sessions[token]
	if !ok {
		return false
	}
	return time.Now().Before(exp)
}

// Delete removes a session
func (s *SessionStore) Delete(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.sessions, token)
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

// Status handles GET /auth/status - returns auth state for frontend
func (h *AuthHandler) Status(c *gin.Context) {
	token, err := c.Cookie(sessionCookieName)
	if err != nil || !globalSessionStore.Valid(token) {
		c.JSON(http.StatusUnauthorized, gin.H{"authenticated": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"authenticated": true})
}

// AuthRequired is a Gin middleware that enforces authentication
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(sessionCookieName)
		if err != nil || !globalSessionStore.Valid(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Non authentifié"})
			c.Abort()
			return
		}
		c.Next()
	}
}

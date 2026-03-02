package main

import (
	"embed"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"personal-coach/handlers"
	"personal-coach/mcp"
	"personal-coach/services"
)

//go:embed dist
var frontendDist embed.FS

func main() {
	// Check if running as MCP server
	if len(os.Args) > 1 && os.Args[1] == "mcp" {
		runMCPServer()
		return
	}

	runHTTPServer()
}

func runHTTPServer() {
	// Initialize services
	claudeService := services.NewClaudeService()
	programHandler := handlers.NewProgramHandler(claudeService)
	authHandler := handlers.NewAuthHandler()

	// Set up Gin router
	r := gin.Default()

	// Auth endpoints (no auth required)
	auth := r.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/logout", authHandler.Logout)
		auth.GET("/status", authHandler.Status)
	}

	// Health check (no auth required)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "personal-coach", "version": "1.1.0"})
	})

	// Protected API routes
	api := r.Group("/api", handlers.AuthRequired())
	{
		programs := api.Group("/programs")
		{
			programs.POST("/generate", programHandler.GenerateProgram)
			programs.GET("", programHandler.ListPrograms)
			programs.GET("/:id", programHandler.GetProgram)
			programs.GET("/:id/pdf", programHandler.DownloadPDF)
			programs.GET("/:id/timer/:day", programHandler.GetTimer)
		}
	}

	// Serve embedded Vue.js frontend (SPA)
	setupFrontend(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Personal Coach server starting on port %s", port)
	log.Printf("Frontend: http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func setupFrontend(r *gin.Engine) {
	// Extract the dist subdirectory
	distFS, err := fs.Sub(frontendDist, "dist")
	if err != nil {
		log.Printf("Warning: could not load embedded frontend: %v", err)
		return
	}

	r.NoRoute(func(c *gin.Context) {
		urlPath := c.Request.URL.Path

		// Don't serve frontend for API/auth routes
		if strings.HasPrefix(urlPath, "/api/") || strings.HasPrefix(urlPath, "/auth/") {
			c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
			return
		}

		// Try to serve the file directly
		cleanPath := path.Clean(strings.TrimPrefix(urlPath, "/"))
		if cleanPath == "." {
			cleanPath = "index.html"
		}

		f, err := distFS.Open(cleanPath)
		if err != nil {
			// SPA fallback: serve index.html for all unmatched routes
			f, err = distFS.Open("index.html")
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Not found"})
				return
			}
			cleanPath = "index.html"
		}
		f.Close()

		// Set content type based on extension
		ext := path.Ext(cleanPath)
		contentType := mime.TypeByExtension(ext)
		if contentType == "" {
			contentType = "application/octet-stream"
		}
		c.Header("Content-Type", contentType)

		// Serve the file
		http.FileServer(http.FS(distFS)).ServeHTTP(c.Writer, c.Request)
	})
}

func runMCPServer() {
	claudeService := services.NewClaudeService()
	mcpServer := mcp.NewMCPServer(claudeService)
	mcpServer.Run()
}

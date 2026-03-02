package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"personal-coach/handlers"
	"personal-coach/mcp"
	"personal-coach/services"
)

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

	// Set up Gin router
	r := gin.Default()

	// CORS configuration for Vue.js frontend
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Disposition"},
		AllowCredentials: true,
	}))

	// API routes
	api := r.Group("/api")
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

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "personal-coach"})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Personal Coach API server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMCPServer() {
	claudeService := services.NewClaudeService()
	mcpServer := mcp.NewMCPServer(claudeService)
	mcpServer.Run()
}

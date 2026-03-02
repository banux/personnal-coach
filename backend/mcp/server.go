package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"personal-coach/models"
	"personal-coach/services"
)

// MCPServer implements the Model Context Protocol server
// It exposes tools for generating workout programs
type MCPServer struct {
	claude *services.ClaudeService
	in     io.Reader
	out    io.Writer
}

// MCPRequest represents an incoming MCP request
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents an MCP response
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Tool represents an MCP tool definition
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// NewMCPServer creates a new MCP server
func NewMCPServer(claude *services.ClaudeService) *MCPServer {
	return &MCPServer{
		claude: claude,
		in:     os.Stdin,
		out:    os.Stdout,
	}
}

// Run starts the MCP server and handles requests
func (s *MCPServer) Run() {
	log.SetOutput(os.Stderr)
	log.Println("MCP Personal Coach Server starting...")

	decoder := json.NewDecoder(s.in)
	encoder := json.NewEncoder(s.out)

	for {
		var req MCPRequest
		if err := decoder.Decode(&req); err != nil {
			if err == io.EOF {
				log.Println("MCP Server: EOF received, shutting down")
				return
			}
			log.Printf("MCP Server: decode error: %v", err)
			continue
		}

		resp := s.handleRequest(req)
		if err := encoder.Encode(resp); err != nil {
			log.Printf("MCP Server: encode error: %v", err)
		}
	}
}

func (s *MCPServer) handleRequest(req MCPRequest) MCPResponse {
	switch req.Method {
	case "initialize":
		return s.handleInitialize(req)
	case "tools/list":
		return s.handleListTools(req)
	case "tools/call":
		return s.handleCallTool(req)
	default:
		return MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Error: &MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("Method not found: %s", req.Method),
			},
		}
	}
}

func (s *MCPServer) handleInitialize(req MCPRequest) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result: map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
			"serverInfo": map[string]interface{}{
				"name":    "personal-coach",
				"version": "1.0.0",
			},
		},
	}
}

func (s *MCPServer) handleListTools(req MCPRequest) MCPResponse {
	tools := []Tool{
		{
			Name:        "generate_workout_program",
			Description: "Génère un programme d'entraînement personnalisé pour une personne selon ses objectifs, niveau et équipement disponible. Le programme inclut des exercices avec séries, répétitions, intensité et conseils techniques.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Nom de la personne",
					},
					"age": map[string]interface{}{
						"type":        "integer",
						"description": "Âge en années",
					},
					"weight": map[string]interface{}{
						"type":        "number",
						"description": "Poids en kg",
					},
					"height": map[string]interface{}{
						"type":        "number",
						"description": "Taille en cm",
					},
					"level": map[string]interface{}{
						"type":        "string",
						"enum":        []string{"beginner", "intermediate", "advanced"},
						"description": "Niveau d'entraînement",
					},
					"goals": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Objectifs: weight_loss, muscle_gain, endurance, strength, flexibility",
					},
					"equipment": map[string]interface{}{
						"type":        "array",
						"items":       map[string]interface{}{"type": "string"},
						"description": "Équipement disponible: barbell, dumbbell, machine, bodyweight, kettlebell, bands",
					},
					"days_per_week": map[string]interface{}{
						"type":        "integer",
						"description": "Nombre de jours d'entraînement par semaine (1-7)",
						"minimum":     1,
						"maximum":     7,
					},
					"weeks": map[string]interface{}{
						"type":        "integer",
						"description": "Durée du programme en semaines",
						"minimum":     1,
						"maximum":     12,
					},
				},
				"required": []string{"name", "days_per_week"},
			},
		},
		{
			Name:        "get_workout_timer",
			Description: "Génère la séquence de timer pour un programme d'entraînement spécifique (exercices avec durées de travail et de repos).",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"program_json": map[string]interface{}{
						"type":        "string",
						"description": "Programme JSON généré par generate_workout_program",
					},
					"day_index": map[string]interface{}{
						"type":        "integer",
						"description": "Index du jour (0-based)",
						"minimum":     0,
					},
				},
				"required": []string{"program_json", "day_index"},
			},
		},
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  map[string]interface{}{"tools": tools},
	}
}

func (s *MCPServer) handleCallTool(req MCPRequest) MCPResponse {
	var params struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal(req.Params, &params); err != nil {
		return s.errorResponse(req.ID, -32602, "Invalid params")
	}

	switch params.Name {
	case "generate_workout_program":
		return s.callGenerateProgram(req.ID, params.Arguments)
	case "get_workout_timer":
		return s.callGetTimer(req.ID, params.Arguments)
	default:
		return s.errorResponse(req.ID, -32601, fmt.Sprintf("Unknown tool: %s", params.Name))
	}
}

func (s *MCPServer) callGenerateProgram(id interface{}, args json.RawMessage) MCPResponse {
	var input struct {
		Name        string   `json:"name"`
		Age         int      `json:"age"`
		Weight      float64  `json:"weight"`
		Height      float64  `json:"height"`
		Level       string   `json:"level"`
		Goals       []string `json:"goals"`
		Equipment   []string `json:"equipment"`
		DaysPerWeek int      `json:"days_per_week"`
		Weeks       int      `json:"weeks"`
	}

	if err := json.Unmarshal(args, &input); err != nil {
		return s.errorResponse(id, -32602, fmt.Sprintf("Invalid arguments: %v", err))
	}

	if input.Level == "" {
		input.Level = "intermediate"
	}
	if input.Weeks == 0 {
		input.Weeks = 4
	}
	if len(input.Equipment) == 0 {
		input.Equipment = []string{"bodyweight"}
	}

	req := models.GenerateRequest{
		Person: models.Person{
			Name:      input.Name,
			Age:       input.Age,
			Weight:    input.Weight,
			Height:    input.Height,
			Level:     input.Level,
			Goals:     input.Goals,
			Equipment: input.Equipment,
		},
		DaysPerWeek: input.DaysPerWeek,
		Weeks:       input.Weeks,
	}

	program, err := s.claude.GenerateProgram(context.Background(), req)
	if err != nil {
		return s.errorResponse(id, -32603, fmt.Sprintf("Generation failed: %v", err))
	}

	programJSON, err := json.MarshalIndent(program, "", "  ")
	if err != nil {
		return s.errorResponse(id, -32603, "Failed to serialize program")
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("Programme généré pour %s:\n\n```json\n%s\n```", input.Name, string(programJSON)),
				},
			},
		},
	}
}

func (s *MCPServer) callGetTimer(id interface{}, args json.RawMessage) MCPResponse {
	var input struct {
		ProgramJSON string `json:"program_json"`
		DayIndex    int    `json:"day_index"`
	}

	if err := json.Unmarshal(args, &input); err != nil {
		return s.errorResponse(id, -32602, fmt.Sprintf("Invalid arguments: %v", err))
	}

	// Clean up the program JSON (might have markdown code blocks)
	programJSON := strings.TrimSpace(input.ProgramJSON)
	if strings.HasPrefix(programJSON, "```json") {
		programJSON = strings.TrimPrefix(programJSON, "```json")
		programJSON = strings.TrimSuffix(programJSON, "```")
		programJSON = strings.TrimSpace(programJSON)
	}

	var program models.Program
	if err := json.Unmarshal([]byte(programJSON), &program); err != nil {
		return s.errorResponse(id, -32602, fmt.Sprintf("Invalid program JSON: %v", err))
	}

	timer, err := services.BuildTimer(program, input.DayIndex)
	if err != nil {
		return s.errorResponse(id, -32603, err.Error())
	}

	timerJSON, err := json.MarshalIndent(timer, "", "  ")
	if err != nil {
		return s.errorResponse(id, -32603, "Failed to serialize timer")
	}

	// Build human-readable timer description
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Timer pour %s (Jour %d):\n", timer.DayName, input.DayIndex+1))
	sb.WriteString(fmt.Sprintf("Durée totale estimée: %d minutes\n\n", timer.TotalTime/60))

	currentExercise := ""
	for _, set := range timer.Sets {
		if set.ExerciseName != currentExercise {
			currentExercise = set.ExerciseName
			sb.WriteString(fmt.Sprintf("\n🏋️ %s\n", currentExercise))
		}
		sb.WriteString(fmt.Sprintf("  Série %d: %s reps (%ds travail) → %ds repos\n",
			set.SetNumber, set.Reps, set.WorkSeconds, set.RestSeconds))
	}

	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result: map[string]interface{}{
			"content": []map[string]interface{}{
				{
					"type": "text",
					"text": fmt.Sprintf("%s\n\n```json\n%s\n```", sb.String(), string(timerJSON)),
				},
			},
		},
	}
}

func (s *MCPServer) errorResponse(id interface{}, code int, message string) MCPResponse {
	return MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
}

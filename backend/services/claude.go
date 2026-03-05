package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"personal-coach/models"
)

// ClaudeService handles all Claude AI interactions
type ClaudeService struct {
	client *anthropic.Client
}

// NewClaudeService creates a new Claude service instance
func NewClaudeService() *ClaudeService {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		panic("ANTHROPIC_API_KEY environment variable is required")
	}

	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &ClaudeService{client: &client}
}

// GenerateProgram calls Claude to generate a personalized workout program
func (s *ClaudeService) GenerateProgram(ctx context.Context, req models.GenerateRequest) (*models.Program, error) {
	prompt := buildProgramPrompt(req)

	message, err := s.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaude3_7SonnetLatest,
		MaxTokens: 8096,
		System: []anthropic.TextBlockParam{
			{
				Text: systemPrompt(),
			},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})
	if err != nil {
		return nil, fmt.Errorf("claude API error: %w", err)
	}

	// Extract text content from response
	var responseText string
	for _, block := range message.Content {
		if block.Type == "text" {
			responseText += block.Text
		}
	}

	// Parse the JSON response from Claude
	program, err := parseProgramResponse(responseText, req)
	if err != nil {
		return nil, fmt.Errorf("failed to parse program response: %w", err)
	}

	return program, nil
}

func systemPrompt() string {
	return `Tu es un coach sportif expert et nutritionniste certifié. Tu génères des programmes d'entraînement personnalisés, structurés et scientifiquement fondés.

RÈGLES IMPORTANTES:
1. Réponds UNIQUEMENT avec du JSON valide, sans texte avant ou après
2. Le JSON doit respecter exactement la structure demandée
3. Adapte le programme selon le niveau, les objectifs, l'équipement disponible, le sexe (si fourni) et le ressenti de la semaine précédente
4. Pour chaque exercice, indique:
   - Le nom précis du mouvement
   - Le nombre de séries et de répétitions
   - L'intensité en %RM ou RPE (Rating of Perceived Exertion)
   - Le tempo si applicable (format: excentrique-pause-concentrique-pause, ex: "3-1-2-0")
   - Les consignes techniques importantes
   - Le temps de repos en secondes
5. Organise le programme en blocs logiques (échauffement, bloc principal, travail accessoire)
6. Prévois une progression logique entre les semaines
7. Si un feedback est fourni, ajuste le volume/intensité en conséquence`
}

func buildProgramPrompt(req models.GenerateRequest) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf(`Génère un programme d'entraînement personnalisé pour:

PROFIL:
- Nom: %s
- Âge: %d ans
- Poids: %.1f kg
- Taille: %.0f cm
- Niveau: %s
- Objectifs: %s`,
		req.Person.Name,
		req.Person.Age,
		req.Person.Weight,
		req.Person.Height,
		req.Person.Level,
		strings.Join(req.Person.Goals, ", "),
	))

	if req.Person.Sex != "" {
		sb.WriteString(fmt.Sprintf("\n- Sexe: %s", req.Person.Sex))
	}

	// Equipment: use detailed items if available, else fall back to simple list
	if len(req.Person.EquipmentItems) > 0 {
		var equipParts []string
		for _, item := range req.Person.EquipmentItems {
			if len(item.Weights) > 0 {
				wStrs := make([]string, len(item.Weights))
				for i, w := range item.Weights {
					wStrs[i] = fmt.Sprintf("%.4gkg", w)
				}
				equipParts = append(equipParts, fmt.Sprintf("%s (%s)", item.Type, strings.Join(wStrs, ", ")))
			} else {
				equipParts = append(equipParts, item.Type)
			}
		}
		sb.WriteString(fmt.Sprintf("\n- Équipement disponible: %s", strings.Join(equipParts, ", ")))
	} else {
		sb.WriteString(fmt.Sprintf("\n- Équipement disponible: %s", strings.Join(req.Person.Equipment, ", ")))
	}

	if req.Person.Description != "" {
		sb.WriteString(fmt.Sprintf("\n- Contexte supplémentaire: %s", req.Person.Description))
	}

	sb.WriteString(fmt.Sprintf(`

PROGRAMME DEMANDÉ:
- Durée: %d semaine(s)
- Fréquence: %d jours par semaine`,
		req.Weeks,
		req.DaysPerWeek,
	))

	if req.Feedback != nil {
		sb.WriteString(fmt.Sprintf(`

RESSENTI DE LA SEMAINE PRÉCÉDENTE:
- Niveau d'énergie: %d/10
- Courbatures: %d/10
- Motivation: %d/10
- Jours complétés: %d
- Notes: %s

AJUSTEMENTS REQUIS: Adapte le volume et l'intensité selon ce ressenti.`,
			req.Feedback.EnergyLevel,
			req.Feedback.SorenessLevel,
			req.Feedback.MotivationLevel,
			req.Feedback.CompletedDays,
			req.Feedback.Notes,
		))
	}

	sb.WriteString(`

Réponds avec un JSON valide correspondant EXACTEMENT à cette structure:
{
  "id": "uuid-string",
  "person_id": "uuid-string",
  "person_name": "string",
  "week_number": 1,
  "total_weeks": number,
  "objective": "string",
  "days": [
    {
      "day": 1,
      "name": "string",
      "focus": "string",
      "duration": number,
      "warmup_notes": "string",
      "cooldown_notes": "string",
      "blocks": [
        {
          "name": "string",
          "exercises": [
            {
              "name": "string",
              "sets": number,
              "reps": "string",
              "intensity": "string",
              "rest_seconds": number,
              "tempo": "string",
              "notes": "string",
              "duration_secs": number,
              "muscle_groups": ["string"]
            }
          ]
        }
      ]
    }
  ],
  "notes": "string"
}`)

	return sb.String()
}

func parseProgramResponse(responseText string, req models.GenerateRequest) (*models.Program, error) {
	// Clean the response - Claude sometimes adds markdown code blocks
	responseText = strings.TrimSpace(responseText)
	if strings.HasPrefix(responseText, "```json") {
		responseText = strings.TrimPrefix(responseText, "```json")
		responseText = strings.TrimSuffix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	} else if strings.HasPrefix(responseText, "```") {
		responseText = strings.TrimPrefix(responseText, "```")
		responseText = strings.TrimSuffix(responseText, "```")
		responseText = strings.TrimSpace(responseText)
	}

	var program models.Program
	if err := json.Unmarshal([]byte(responseText), &program); err != nil {
		return nil, fmt.Errorf("invalid JSON from Claude: %w\nResponse: %s", err, responseText[:min(len(responseText), 500)])
	}

	// Ensure required fields are set
	if program.PersonName == "" {
		program.PersonName = req.Person.Name
	}
	if program.TotalWeeks == 0 {
		program.TotalWeeks = req.Weeks
	}
	if req.Feedback != nil {
		program.Feedback = req.Feedback
	}

	return &program, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

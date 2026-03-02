package services

import (
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"personal-coach/models"
)

// GeneratePDF creates a PDF document for a workout program
func GeneratePDF(program models.Program) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)

	for i, day := range program.Days {
		pdf.AddPage()

		// Header
		pdf.SetFillColor(240, 248, 255) // Light blue background
		pdf.Rect(0, 0, 210, 30, "F")

		pdf.SetFont("Helvetica", "B", 18)
		pdf.SetTextColor(30, 60, 120)
		pdf.CellFormat(180, 10, fmt.Sprintf("Programme - %s", program.PersonName), "", 1, "C", false, 0, "")

		pdf.SetFont("Helvetica", "", 11)
		pdf.SetTextColor(80, 80, 80)
		pdf.CellFormat(180, 8, fmt.Sprintf("Semaine %d/%d | %s", program.WeekNumber, program.TotalWeeks, program.Objective), "", 1, "C", false, 0, "")

		pdf.Ln(5)

		// Day header
		pdf.SetFillColor(30, 60, 120)
		pdf.SetTextColor(255, 255, 255)
		pdf.SetFont("Helvetica", "B", 14)
		pdf.CellFormat(180, 10, fmt.Sprintf("Jour %d: %s - %s", day.Day, day.Name, day.Focus), "", 1, "C", true, 0, "")
		pdf.SetTextColor(0, 0, 0)

		pdf.SetFont("Helvetica", "I", 9)
		pdf.SetTextColor(100, 100, 100)
		pdf.CellFormat(180, 6, fmt.Sprintf("Durée estimée: %d min", day.Duration), "", 1, "L", false, 0, "")
		pdf.Ln(2)

		// Warm-up notes
		if day.WarmupNotes != "" {
			pdf.SetFillColor(255, 250, 230)
			pdf.SetFont("Helvetica", "B", 10)
			pdf.SetTextColor(180, 100, 0)
			pdf.CellFormat(180, 7, "Échauffement", "", 1, "L", true, 0, "")
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(60, 60, 60)
			pdf.MultiCell(180, 5, day.WarmupNotes, "LRB", "L", false)
			pdf.Ln(3)
		}

		// Exercise blocks
		for _, block := range day.Blocks {
			// Block header
			pdf.SetFillColor(220, 235, 255)
			pdf.SetFont("Helvetica", "B", 11)
			pdf.SetTextColor(30, 60, 120)
			pdf.CellFormat(180, 8, block.Name, "1", 1, "L", true, 0, "")

			// Exercise table header
			pdf.SetFillColor(200, 215, 240)
			pdf.SetFont("Helvetica", "B", 9)
			pdf.SetTextColor(30, 60, 120)
			pdf.CellFormat(50, 7, "Exercice", "1", 0, "C", true, 0, "")
			pdf.CellFormat(15, 7, "Séries", "1", 0, "C", true, 0, "")
			pdf.CellFormat(25, 7, "Reps", "1", 0, "C", true, 0, "")
			pdf.CellFormat(30, 7, "Intensité", "1", 0, "C", true, 0, "")
			pdf.CellFormat(20, 7, "Repos", "1", 0, "C", true, 0, "")
			pdf.CellFormat(40, 7, "Tempo/Notes", "1", 1, "C", true, 0, "")

			// Exercises
			for j, ex := range block.Exercises {
				fillColor := [3]int{248, 252, 255}
				if j%2 == 0 {
					fillColor = [3]int{240, 245, 255}
				}
				pdf.SetFillColor(fillColor[0], fillColor[1], fillColor[2])
				pdf.SetFont("Helvetica", "", 9)
				pdf.SetTextColor(30, 30, 30)

				// Calculate max height needed
				notesText := ex.Tempo
				if ex.Notes != "" {
					if notesText != "" {
						notesText += " | "
					}
					notesText += ex.Notes
				}

				restText := fmt.Sprintf("%ds", ex.RestSeconds)

				pdf.CellFormat(50, 7, ex.Name, "1", 0, "L", true, 0, "")
				pdf.CellFormat(15, 7, fmt.Sprintf("%d", ex.Sets), "1", 0, "C", true, 0, "")
				pdf.CellFormat(25, 7, ex.Reps, "1", 0, "C", true, 0, "")
				pdf.CellFormat(30, 7, ex.Intensity, "1", 0, "C", true, 0, "")
				pdf.CellFormat(20, 7, restText, "1", 0, "C", true, 0, "")
				pdf.CellFormat(40, 7, truncate(notesText, 25), "1", 1, "L", true, 0, "")

				// If there are notes, add them below
				if len(notesText) > 25 {
					pdf.SetFont("Helvetica", "I", 8)
					pdf.SetTextColor(100, 100, 100)
					pdf.CellFormat(50, 5, "", "", 0, "", false, 0, "")
					pdf.MultiCell(130, 4, notesText, "", "L", false)
				}
			}
			pdf.Ln(3)
		}

		// Cooldown notes
		if day.CooldownNotes != "" {
			pdf.SetFillColor(230, 255, 240)
			pdf.SetFont("Helvetica", "B", 10)
			pdf.SetTextColor(0, 120, 60)
			pdf.CellFormat(180, 7, "Retour au calme", "", 1, "L", true, 0, "")
			pdf.SetFont("Helvetica", "", 9)
			pdf.SetTextColor(60, 60, 60)
			pdf.MultiCell(180, 5, day.CooldownNotes, "LRB", "L", false)
		}

		// Page number
		pdf.SetY(-15)
		pdf.SetFont("Helvetica", "I", 8)
		pdf.SetTextColor(150, 150, 150)
		pdf.CellFormat(180, 10, fmt.Sprintf("Page %d | Personnel Coach AI", i+1), "", 0, "C", false, 0, "")
	}

	// Notes page if there are program notes
	if program.Notes != "" {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 14)
		pdf.SetTextColor(30, 60, 120)
		pdf.CellFormat(180, 10, "Notes du programme", "", 1, "L", false, 0, "")
		pdf.Ln(3)
		pdf.SetFont("Helvetica", "", 10)
		pdf.SetTextColor(60, 60, 60)
		pdf.MultiCell(180, 6, program.Notes, "", "L", false)
	}

	// Output to bytes
	var buf strings.Builder
	err := pdf.Output(&stringWriter{&buf})
	if err != nil {
		return nil, fmt.Errorf("PDF generation error: %w", err)
	}

	return []byte(buf.String()), nil
}

type stringWriter struct {
	sb *strings.Builder
}

func (sw *stringWriter) Write(p []byte) (n int, err error) {
	return sw.sb.Write([]byte(string(p)))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

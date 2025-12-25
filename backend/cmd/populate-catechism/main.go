package main

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// OnlineCatechismItem estrutura para parsing do JSON online
type OnlineCatechismItem struct {
	Number int    `json:"number"`
	Q      string `json:"q"` // Question (pode ser "question" também)
	A      string `json:"a"` // Answer (pode ser "answer" também)
	Question string `json:"question"` // Formato alternativo
	Answer   string `json:"answer"`   // Formato alternativo
}

func main() {
	var clearFlag = flag.Bool("clear", false, "Clear existing catechism questions before populating")
	var urlFlag = flag.String("url", "", "Custom URL to fetch catechism from (optional)")
	flag.Parse()

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	repo := repository.NewCatechismRepository()

	// Clear existing questions if flag is set
	if *clearFlag {
		log.Println("Clearing existing catechism questions...")
		_, err := database.DB.Exec("DELETE FROM westminster_catechism")
		if err != nil {
			log.Fatalf("Failed to clear catechism: %v", err)
		}
		log.Println("Existing questions cleared.")
	}

	// Tentar primeiro arquivo local, depois URL online
	var body []byte

	// Tentar múltiplos caminhos possíveis (dependendo de onde o comando é executado)
	// Prioriza Catecismo Maior, depois Menor
	possiblePaths := []string{
		"catechism_maior.json",                             // Catecismo Maior - quando executado de dentro de cmd/populate-catechism
		"cmd/populate-catechism/catechism_maior.json",      // Catecismo Maior - quando executado da raiz do backend
		"/app/cmd/populate-catechism/catechism_maior.json", // Catecismo Maior - caminho absoluto no container
		"catechism.json",                                    // Catecismo Menor (fallback)
		"cmd/populate-catechism/catechism.json",             // Catecismo Menor (fallback)
		"/app/cmd/populate-catechism/catechism.json",        // Catecismo Menor (fallback)
	}

	var localFile string
	var found bool
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			localFile = path
			found = true
			break
		}
	}

	// Debug: mostrar diretório atual
	if wd, err := os.Getwd(); err == nil {
		log.Printf("Current working directory: %s", wd)
	}

	if found {
		log.Printf("Reading catechism from local file: %s", localFile)
		var err error
		body, err = os.ReadFile(localFile)
		if err != nil {
			log.Fatalf("Failed to read local file: %v", err)
		}
	} else {
		// Tentar buscar online
		catechismURL := *urlFlag
		if catechismURL == "" {
			catechismURL = "https://raw.githubusercontent.com/ReformedWiki/westminster-shorter-catechism/master/data/catechism.json"
		}

		log.Printf("Fetching catechism from: %s", catechismURL)
		resp, err := http.Get(catechismURL)
		if err != nil {
			log.Fatalf("Failed to fetch catechism: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			log.Fatalf("Failed to fetch catechism (HTTP %d): %s", resp.StatusCode, string(body))
		}

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("Failed to read response: %v", err)
		}
	}

	log.Println("Populating Westminster Catechism...")

	// Parse JSON
	var items []OnlineCatechismItem
	if err := json.Unmarshal(body, &items); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(items) == 0 {
		log.Fatalf("No questions found in the response")
	}

	// Converter para nosso modelo e salvar
	questions := make([]*models.CatechismQuestion, 0, len(items))
	validCount := 0
	maxQuestion := 0

	for _, item := range items {
		if item.Number >= 1 {
			// Suporta ambos os formatos: "q"/"a" ou "question"/"answer"
			questionText := item.Q
			if questionText == "" {
				questionText = item.Question
			}
			answerText := item.A
			if answerText == "" {
				answerText = item.Answer
			}
			
			questions = append(questions, &models.CatechismQuestion{
				QuestionNumber: item.Number,
				QuestionText:   strings.TrimSpace(questionText),
				AnswerText:     strings.TrimSpace(answerText),
			})
			validCount++
			if item.Number > maxQuestion {
				maxQuestion = item.Number
			}
		}
	}

	if validCount == 0 {
		log.Fatalf("No valid questions found")
	}

	log.Printf("Found %d valid questions (1-%d), saving to database...", validCount, maxQuestion)

	// Salvar no banco de dados
	for i, question := range questions {
		if err := repo.Create(question); err != nil {
			log.Printf("Failed to save question %d: %v", question.QuestionNumber, err)
			continue
		}

		if (i+1)%20 == 0 {
			log.Printf("Saved %d/%d questions...", i+1, len(questions))
		}
	}

	log.Printf("✅ Successfully populated %d questions!", validCount)
	log.Printf("✅ Questions range from 1 to %d", maxQuestion)
}

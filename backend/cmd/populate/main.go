package main

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	var clearFlag = flag.Bool("clear", false, "Clear existing reading plans before populating")
	flag.Parse()

	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()

	repo := repository.NewReadingPlanRepository()

	// Clear existing plans if flag is set
	if *clearFlag {
		log.Println("Clearing existing reading plans...")
		_, err := database.DB.Exec("DELETE FROM reading_plans")
		if err != nil {
			log.Fatalf("Failed to clear reading plans: %v", err)
		}
	}

	// Generate reading plan for 365 days
	// This is a simplified version - you can expand this with actual Bible reading plan data
	log.Println("Populating reading plans for 365 days...")

	// Example structure: This is a placeholder - you'll need to populate with actual Bible references
	// For now, creating a basic structure that cycles through books
	oldTestamentBooks := []string{
		"Gênesis", "Êxodo", "Levítico", "Números", "Deuteronômio",
		"Josué", "Juízes", "Rute", "1 Samuel", "2 Samuel",
		"1 Reis", "2 Reis", "1 Crônicas", "2 Crônicas", "Esdras",
		"Neemias", "Ester", "Jó", "Salmos", "Provérbios",
		"Eclesiastes", "Cantares", "Isaías", "Jeremias", "Lamentações",
		"Ezequiel", "Daniel", "Oséias", "Joel", "Amós",
		"Obadias", "Jonas", "Miquéias", "Naum", "Habacuque",
		"Sofonias", "Ageu", "Zacarias", "Malaquias",
	}

	newTestamentBooks := []string{
		"Mateus", "Marcos", "Lucas", "João", "Atos",
		"Romanos", "1 Coríntios", "2 Coríntios", "Gálatas", "Efésios",
		"Filipenses", "Colossenses", "1 Tessalonicenses", "2 Tessalonicenses", "1 Timóteo",
		"2 Timóteo", "Tito", "Filemom", "Hebreus", "Tiago",
		"1 Pedro", "2 Pedro", "1 João", "2 João", "3 João",
		"Judas", "Apocalipse",
	}

	// Create a basic reading plan
	// This is a simplified version - adjust based on your actual Bible 365 plan
	for day := 1; day <= 365; day++ {
		// Distribute Old Testament readings (approximately 231 days)
		var oldTestamentRef string
		if day <= 231 {
			bookIndex := (day - 1) % len(oldTestamentBooks)
			chapter := ((day - 1) / len(oldTestamentBooks)) + 1
			oldTestamentRef = fmt.Sprintf("%s %d", oldTestamentBooks[bookIndex], chapter)
		}

		// Distribute New Testament readings (approximately 89 days)
		var newTestamentRef string
		if day <= 89 {
			bookIndex := (day - 1) % len(newTestamentBooks)
			chapter := ((day - 1) / len(newTestamentBooks)) + 1
			newTestamentRef = fmt.Sprintf("%s %d", newTestamentBooks[bookIndex], chapter)
		}

		// Psalms - read through once (150 chapters, ~150 days)
		var psalmsRef string
		if day <= 150 {
			psalmsRef = fmt.Sprintf("Salmos %d", day)
		}

		// Proverbs - read through multiple times (31 chapters, cycle through)
		proverbsRef := fmt.Sprintf("Provérbios %d", ((day-1)%31)+1)

		plan := &models.ReadingPlan{
			DayOfYear:       day,
			OldTestamentRef: oldTestamentRef,
			NewTestamentRef: newTestamentRef,
			PsalmsRef:       psalmsRef,
			ProverbsRef:     proverbsRef,
		}

		if err := repo.Create(plan); err != nil {
			log.Printf("Failed to create plan for day %d: %v", day, err)
			// Continue with next day
		}

		if day%50 == 0 {
			log.Printf("Processed %d days...", day)
		}
	}

	log.Println("Reading plans populated successfully!")
}


package main

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"flag"
	"fmt"
	"log"
	"strings"
)

// consumeChapters avança sequencialmente pelos capítulos e devolve uma string com as referências do dia.
// Garante que, se precisarmos ler mais de um capítulo no mesmo dia, todos sejam listados.
func consumeChapters(books []struct {
	name     string
	chapters int
}, bookIndex *int, chapter *int, chaptersToRead int) string {
	if chaptersToRead <= 0 || *bookIndex >= len(books) {
		return ""
	}

	var refs []string

	for i := 0; i < chaptersToRead && *bookIndex < len(books); i++ {
		book := books[*bookIndex]
		refs = append(refs, fmt.Sprintf("%s %d", book.name, *chapter))

		*chapter++
		if *chapter > book.chapters {
			*chapter = 1
			*bookIndex++
		}
	}

	return strings.Join(refs, "; ")
}

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

	log.Println("Populating reading plans for 365 days following Bíblia 365 pattern...")

	// Bíblia 365:
	// MANHÃ = Antigo Testamento + Salmos (ambos sequenciais)
	// NOITE = Novo Testamento + Provérbios (ambos sequenciais)
	// Quando Salmos ou Provérbios terminarem, reiniciam do capítulo 1

	// Antigo Testamento (sem Salmos e Provérbios que são tratados separadamente)
	oldTestamentBooks := []struct {
		name     string
		chapters int
	}{
		{"Gênesis", 50}, {"Êxodo", 40}, {"Levítico", 27}, {"Números", 36}, {"Deuteronômio", 34},
		{"Josué", 24}, {"Juízes", 21}, {"Rute", 4}, {"1 Samuel", 31}, {"2 Samuel", 24},
		{"1 Reis", 22}, {"2 Reis", 25}, {"1 Crônicas", 29}, {"2 Crônicas", 36}, {"Esdras", 10},
		{"Neemias", 13}, {"Ester", 10}, {"Jó", 42},
		{"Eclesiastes", 12}, {"Cantares", 8}, {"Isaías", 66}, {"Jeremias", 52}, {"Lamentações", 5},
		{"Ezequiel", 48}, {"Daniel", 12}, {"Oséias", 14}, {"Joel", 3}, {"Amós", 9},
		{"Obadias", 1}, {"Jonas", 4}, {"Miquéias", 7}, {"Naum", 3}, {"Habacuque", 3},
		{"Sofonias", 3}, {"Ageu", 2}, {"Zacarias", 14}, {"Malaquias", 4},
	}

	// Novo Testamento
	newTestamentBooks := []struct {
		name     string
		chapters int
	}{
		{"Mateus", 28}, {"Marcos", 16}, {"Lucas", 24}, {"João", 21}, {"Atos", 28},
		{"Romanos", 16}, {"1 Coríntios", 16}, {"2 Coríntios", 13}, {"Gálatas", 6}, {"Efésios", 6},
		{"Filipenses", 4}, {"Colossenses", 4}, {"1 Tessalonicenses", 5}, {"2 Tessalonicenses", 3}, {"1 Timóteo", 6},
		{"2 Timóteo", 4}, {"Tito", 3}, {"Filemom", 1}, {"Hebreus", 13}, {"Tiago", 5},
		{"1 Pedro", 5}, {"2 Pedro", 3}, {"1 João", 5}, {"2 João", 1}, {"3 João", 1},
		{"Judas", 1}, {"Apocalipse", 22},
	}

	// Calcular total de capítulos
	totalOTChapters := 0
	for _, book := range oldTestamentBooks {
		totalOTChapters += book.chapters
	}

	totalNTChapters := 0
	for _, book := range newTestamentBooks {
		totalNTChapters += book.chapters
	}

	// Distribuir capítulos ao longo de 365 dias
	// Manhã: AT sequencial + Salmos sequencial
	// Noite: NT sequencial + Provérbios sequencial
	otChaptersPerDay := float64(totalOTChapters) / 365.0
	ntChaptersPerDay := float64(totalNTChapters) / 365.0

	// Contadores para distribuição sequencial
	otBookIndex := 0
	otChapter := 1
	otAccumulator := 0.0

	ntBookIndex := 0
	ntChapter := 1
	ntAccumulator := 0.0

	for day := 1; day <= 365; day++ {
		var oldTestamentRef string
		var newTestamentRef string
		var psalmsRef string
		var proverbsRef string

		// MANHÃ: Antigo Testamento (sequencial)
		otAccumulator += otChaptersPerDay
		chaptersToRead := int(otAccumulator)
		otAccumulator -= float64(chaptersToRead)

		oldTestamentRef = consumeChapters(oldTestamentBooks, &otBookIndex, &otChapter, chaptersToRead)

		// MANHÃ: Salmos (sequencial, reinicia quando terminar)
		// Calcular baseado no dia do ano: (day - 1) % 150 + 1
		psalmDay := ((day - 1) % 150) + 1
		psalmsRef = fmt.Sprintf("Salmos %d", psalmDay)

		// NOITE: Novo Testamento (sequencial)
		ntAccumulator += ntChaptersPerDay
		ntChaptersToRead := int(ntAccumulator)
		ntAccumulator -= float64(ntChaptersToRead)

		newTestamentRef = consumeChapters(newTestamentBooks, &ntBookIndex, &ntChapter, ntChaptersToRead)

		// NOITE: Provérbios (sequencial, reinicia quando terminar)
		// Calcular baseado no dia do ano: (day - 1) % 31 + 1
		proverbDay := ((day - 1) % 31) + 1
		proverbsRef = fmt.Sprintf("Provérbios %d", proverbDay)

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

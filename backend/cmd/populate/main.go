package main

import (
	"biblia-am-pm/internal/database"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"flag"
	"log"
	"strings"
)

// isOldTestament verifica se uma leitura é do Antigo Testamento
// Nota: Salmos (Sl) e Provérbios (Pv) são tratados separadamente
func isOldTestament(ref string) bool {
	otBooks := []string{
		"Gn", "Êx", "Lv", "Nm", "Dt", "Js", "Jz", "Rt", "1 Sm", "2 Sm",
		"1 Rs", "2 Rs", "1 Cr", "2 Cr", "Ed", "Ne", "Et", "Jó", "Ec", "Ct",
		"Is", "Jr", "Lm", "Ez", "Dn", "Os", "Jl", "Am", "Ob", "Jn",
		"Mq", "Na", "Hc", "Sf", "Ag", "Zc", "Ml",
	}
	refUpper := strings.ToUpper(ref)
	for _, book := range otBooks {
		if strings.HasPrefix(refUpper, strings.ToUpper(book)) {
			return true
		}
	}
	return false
}

// isNewTestament verifica se uma leitura é do Novo Testamento
func isNewTestament(ref string) bool {
	ntBooks := []string{
		"Mt", "Mc", "Lc", "Jo", "At", "Rm", "1 Co", "2 Co", "Gl", "Ef",
		"Fp", "Cl", "1 Ts", "2 Ts", "1 Tm", "2 Tm", "Tt", "Fl", "Hb", "Tg",
		"1 Pe", "2 Pe", "1 Jo", "2 Jo", "3 Jo", "Jd", "Ap",
	}
	refUpper := strings.ToUpper(ref)
	for _, book := range ntBooks {
		if strings.HasPrefix(refUpper, strings.ToUpper(book)) {
			return true
		}
	}
	return false
}

// isPsalms verifica se uma leitura é de Salmos
func isPsalms(ref string) bool {
	refUpper := strings.ToUpper(ref)
	return strings.HasPrefix(refUpper, "SL") || strings.HasPrefix(refUpper, "SALMOS")
}

// isProverbs verifica se uma leitura é de Provérbios
func isProverbs(ref string) bool {
	refUpper := strings.ToUpper(ref)
	return strings.HasPrefix(refUpper, "PV") || strings.HasPrefix(refUpper, "PROVÉRBIOS")
}

// mapRMMToHybrid mapeia as 4 leituras do RMM para a estrutura híbrida (manhã/noite)
func mapRMMToHybrid(rmmDay RMMDay) (oldTestamentRef, psalmsRef, newTestamentRef, proverbsRef string) {
	var otReadings []string
	var ntReadings []string
	var psalmsReading string
	var proverbsReading string

	// Analisar cada leitura
	readings := []string{rmmDay.Reading1, rmmDay.Reading2, rmmDay.Reading3, rmmDay.Reading4}
	
	for _, reading := range readings {
		if reading == "" {
			continue
		}
		
		if isPsalms(reading) {
			psalmsReading = reading
		} else if isProverbs(reading) {
			proverbsReading = reading
		} else if isOldTestament(reading) {
			otReadings = append(otReadings, reading)
		} else if isNewTestament(reading) {
			ntReadings = append(ntReadings, reading)
		}
	}

	// Agrupar leituras
	oldTestamentRef = strings.Join(otReadings, "; ")
	if oldTestamentRef == "" {
		oldTestamentRef = ""
	}

	psalmsRef = psalmsReading

	// Para NT, usar a primeira leitura encontrada (geralmente Reading2)
	if len(ntReadings) > 0 {
		newTestamentRef = ntReadings[0]
	} else {
		newTestamentRef = ""
	}

	proverbsRef = proverbsReading

	return oldTestamentRef, psalmsRef, newTestamentRef, proverbsRef
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

	log.Println("Populating reading plans for 365 days following Robert Murray M'Cheyne plan...")

	// Carregar plano RMM
	rmmPlan := getRMMPlan()

	// Processar cada dia
	for day := 1; day <= 365; day++ {
		rmmDay, exists := rmmPlan[day]
		if !exists {
			log.Printf("Warning: No RMM plan found for day %d", day)
			continue
		}

		// Mapear para estrutura híbrida
		oldTestamentRef, psalmsRef, newTestamentRef, proverbsRef := mapRMMToHybrid(rmmDay)

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

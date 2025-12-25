package handlers

import (
	"biblia-am-pm/internal/middleware"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type CatechismHandler struct {
	catechismRepo        *repository.CatechismRepository
	catechismProgressRepo *repository.CatechismProgressRepository
}

func NewCatechismHandler() *CatechismHandler {
	return &CatechismHandler{
		catechismRepo:         repository.NewCatechismRepository(),
		catechismProgressRepo: repository.NewCatechismProgressRepository(),
	}
}

// getLocalTime returns the current time in the configured timezone
func getLocalTimeForCatechism() time.Time {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "America/Sao_Paulo" // Default timezone
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		loc = time.UTC
	}
	now := time.Now().In(loc)
	return now
}

// getWeekStart returns the Sunday of the current week
func getWeekStart(date time.Time) time.Time {
	weekday := int(date.Weekday())
	// Sunday is 0, so we need to handle it
	if weekday == 0 {
		// It's already Sunday
		return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	}
	// Calculate days to subtract to get to Sunday
	daysToSubtract := weekday
	return date.AddDate(0, 0, -daysToSubtract)
}

// getCurrentQuestionNumber calculates which question should be active this week
// Based on the number of weeks since a reference date (first Sunday of 2024)
func getCurrentQuestionNumber(now time.Time, totalQuestions int) int {
	// Reference date: First Sunday of 2024 (January 7, 2024)
	referenceDate := time.Date(2024, 1, 7, 0, 0, 0, 0, now.Location())
	
	// Get the Sunday of the current week
	currentWeekStart := getWeekStart(now)
	
	// Calculate weeks since reference
	weeksSinceReference := int(currentWeekStart.Sub(referenceDate).Hours() / 24 / 7)
	
	// Calculate question number (1-totalQuestions, cycling)
	questionNumber := (weeksSinceReference % totalQuestions) + 1
	
	// Ensure it's between 1 and totalQuestions
	if questionNumber < 1 {
		questionNumber = totalQuestions
	}
	if questionNumber > totalQuestions {
		questionNumber = 1
	}
	
	return questionNumber
}

type CurrentQuestionResponse struct {
	Question        *models.CatechismQuestion `json:"question"`
	WeekProgress    []*models.CatechismProgress `json:"week_progress"`
	WeekStart       string                     `json:"week_start"`
	WeekEnd         string                     `json:"week_end"`
	NextQuestionDate string                    `json:"next_question_date"`
	QuestionNumber  int                        `json:"question_number"`
	TotalQuestions  int                        `json:"total_questions"`
}

func (h *CatechismHandler) GetCurrentQuestion(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Get total number of questions from database
	totalQuestions, err := h.catechismRepo.GetMaxQuestionNumber()
	if err != nil || totalQuestions == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total questions. Please populate the catechism first."})
		return
	}

	now := getLocalTimeForCatechism()
	questionNumber := getCurrentQuestionNumber(now, totalQuestions)
	
	// Get the question
	question, err := h.catechismRepo.GetByQuestionNumber(questionNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get question"})
		return
	}
	
	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found. Please populate the catechism first."})
		return
	}
	
	// Get week start (Sunday)
	weekStart := getWeekStart(now)
	weekEnd := weekStart.AddDate(0, 0, 6)
	nextQuestionDate := weekStart.AddDate(0, 0, 7)
	
	// Get progress for this week
	weekProgress, err := h.catechismProgressRepo.GetByUserAndQuestionForWeek(userID, question.ID, weekStart)
	if err != nil {
		log.Printf("Error getting week progress: %v", err)
		weekProgress = []*models.CatechismProgress{}
	}
	
	response := CurrentQuestionResponse{
		Question:         question,
		WeekProgress:     weekProgress,
		WeekStart:        weekStart.Format("2006-01-02"),
		WeekEnd:          weekEnd.Format("2006-01-02"),
		NextQuestionDate: nextQuestionDate.Format("2006-01-02"),
		QuestionNumber:   questionNumber,
		TotalQuestions:   totalQuestions,
	}
	
	c.JSON(http.StatusOK, response)
}

type MarkCatechismCompletedRequest struct {
	Date string `json:"date"` // Optional, defaults to today
}

func (h *CatechismHandler) MarkAsCompleted(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req MarkCatechismCompletedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get total number of questions from database
	totalQuestions, err := h.catechismRepo.GetMaxQuestionNumber()
	if err != nil || totalQuestions == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get total questions. Please populate the catechism first."})
		return
	}

	now := getLocalTimeForCatechism()
	var targetDate time.Time
	
	if req.Date != "" {
		parsedDate, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
			return
		}
		targetDate = parsedDate
	} else {
		targetDate = now
	}
	
	questionNumber := getCurrentQuestionNumber(targetDate, totalQuestions)
	
	// Get the question
	question, err := h.catechismRepo.GetByQuestionNumber(questionNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get question"})
		return
	}
	
	if question == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Question not found"})
		return
	}
	
	// Get or create progress
	progress, err := h.catechismProgressRepo.GetByUserAndDate(userID, question.ID, targetDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}
	
	if progress == nil {
		progress = &models.CatechismProgress{
			UserID:    userID,
			QuestionID: question.ID,
			Date:      targetDate,
			Completed: false,
		}
	}
	
	// Mark as completed
	progress.Completed = true
	
	// Save progress
	err = h.catechismProgressRepo.CreateOrUpdate(progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}
	
	c.JSON(http.StatusOK, progress)
}

func (h *CatechismHandler) GetProgress(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	progresses, err := h.catechismProgressRepo.GetUserProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}

	c.JSON(http.StatusOK, progresses)
}

// Structure for parsing online catechism data
type OnlineCatechismItem struct {
	Number int    `json:"number"`
	Q      string `json:"q"` // Question
	A      string `json:"a"` // Answer
}

func (h *CatechismHandler) PopulateCatechism(c *gin.Context) {
	// Try to fetch from a public API or source
	// For now, we'll use a known source or allow manual population
	// You can use https://www.westminsterconfession.org/resources/westminster-shorter-catechism.php
	// or similar sources
	
	url := "https://raw.githubusercontent.com/ReformedWiki/westminster-shorter-catechism/master/data/catechism.json"
	
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch catechism: %v", err)})
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch catechism: %s", string(body))})
		return
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response: %v", err)})
		return
	}
	
	var items []OnlineCatechismItem
	if err := json.Unmarshal(body, &items); err != nil {
		// Try alternative format or manual parsing
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse JSON: %v. You may need to populate manually.", err)})
		return
	}
	
	// Convert to our model
	questions := make([]*models.CatechismQuestion, 0, len(items))
	for _, item := range items {
		if item.Number >= 1 && item.Number <= 107 {
			questions = append(questions, &models.CatechismQuestion{
				QuestionNumber: item.Number,
				QuestionText:   strings.TrimSpace(item.Q),
				AnswerText:     strings.TrimSpace(item.A),
			})
		}
	}
	
	if len(questions) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "No valid questions found in the response"})
		return
	}
	
	// Save to database
	if err := h.catechismRepo.CreateBatch(questions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to save questions: %v", err)})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message":   fmt.Sprintf("Successfully populated %d questions", len(questions)),
		"count":     len(questions),
	})
}


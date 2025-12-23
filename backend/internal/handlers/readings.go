package handlers

import (
	"biblia-am-pm/internal/middleware"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

type ReadingsHandler struct {
	readingPlanRepo  *repository.ReadingPlanRepository
	userProgressRepo *repository.UserProgressRepository
}

func NewReadingsHandler() *ReadingsHandler {
	return &ReadingsHandler{
		readingPlanRepo:  repository.NewReadingPlanRepository(),
		userProgressRepo: repository.NewUserProgressRepository(),
	}
}

type TodayReadingsResponse struct {
	Period    string               `json:"period"`
	Readings  *models.ReadingPlan  `json:"readings"`
	Progress  *models.UserProgress `json:"progress"`
	DayOfYear int                  `json:"day_of_year"`
	PlanName  string               `json:"plan_name"`
}

type MarkCompletedRequest struct {
	Period string `json:"period"` // "morning" or "evening"
}

// getLocalTime returns the current time in the configured timezone
func getLocalTime() time.Time {
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "America/Sao_Paulo" // Default timezone
	}
	loc, err := time.LoadLocation(tz)
	if err != nil {
		// Fallback to UTC if timezone is invalid
		loc = time.UTC
	}
	now := time.Now().In(loc)
	return now
}

func (h *ReadingsHandler) GetTodayReadings(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	now := getLocalTime()
	dayOfYear := now.YearDay()
	hour := now.Hour()

	// Determine period based on time
	var period string
	if hour >= 6 && hour < 12 {
		period = "morning"
	} else if hour >= 18 && hour < 23 {
		period = "evening"
	} else {
		period = "all"
	}

	// Debug: log current time and period
	tz := os.Getenv("TZ")
	if tz == "" {
		tz = "America/Sao_Paulo"
	}
	log.Printf("[DEBUG] Timezone: %s, Current time: %s, Hour: %d, Period: %s", tz, now.Format("2006-01-02 15:04:05 MST"), hour, period)

	// Get reading plan for today
	plan, err := h.readingPlanRepo.GetByDayOfYear(dayOfYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reading plan"})
		return
	}

	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reading plan not found for today"})
		return
	}

	// Get user progress for today
	progress, err := h.userProgressRepo.GetByUserAndDate(userID, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}

	// If no progress exists, create a new one
	if progress == nil {
		progress = &models.UserProgress{
			UserID:           userID,
			ReadingPlanID:    plan.ID,
			Date:             now,
			MorningCompleted: false,
			EveningCompleted: false,
		}
	}

	response := TodayReadingsResponse{
		Period:    period,
		Readings:  plan,
		Progress:  progress,
		DayOfYear: dayOfYear,
		PlanName:  "Robert Murray M'Cheyne",
	}

	c.JSON(http.StatusOK, response)
}

func (h *ReadingsHandler) MarkCompleted(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req MarkCompletedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Period != "morning" && req.Period != "evening" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Period must be 'morning' or 'evening'"})
		return
	}

	now := getLocalTime()
	dayOfYear := now.YearDay()

	// Get reading plan for today
	plan, err := h.readingPlanRepo.GetByDayOfYear(dayOfYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reading plan"})
		return
	}

	if plan == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reading plan not found for today"})
		return
	}

	// Get or create progress
	progress, err := h.userProgressRepo.GetByUserAndDate(userID, now)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}

	if progress == nil {
		progress = &models.UserProgress{
			UserID:           userID,
			ReadingPlanID:    plan.ID,
			Date:             now,
			MorningCompleted: false,
			EveningCompleted: false,
		}
	}

	// Update progress based on period
	if req.Period == "morning" {
		progress.MorningCompleted = true
	} else {
		progress.EveningCompleted = true
	}

	// Save progress
	err = h.userProgressRepo.CreateOrUpdate(progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save progress"})
		return
	}

	c.JSON(http.StatusOK, progress)
}

func (h *ReadingsHandler) GetProgress(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	progresses, err := h.userProgressRepo.GetUserProgress(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get progress"})
		return
	}

	c.JSON(http.StatusOK, progresses)
}

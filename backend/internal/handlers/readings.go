package handlers

import (
	"biblia-am-pm/internal/middleware"
	"biblia-am-pm/internal/models"
	"biblia-am-pm/internal/repository"
	"encoding/json"
	"net/http"
	"time"
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
	Period      string                  `json:"period"`
	Readings    *models.ReadingPlan     `json:"readings"`
	Progress    *models.UserProgress    `json:"progress"`
	DayOfYear   int                     `json:"day_of_year"`
}

type MarkCompletedRequest struct {
	Period string `json:"period"` // "morning" or "evening"
}

func (h *ReadingsHandler) GetTodayReadings(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		middleware.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	now := time.Now()
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

	// Get reading plan for today
	plan, err := h.readingPlanRepo.GetByDayOfYear(dayOfYear)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get reading plan")
		return
	}

	if plan == nil {
		middleware.RespondError(w, http.StatusNotFound, "Reading plan not found for today")
		return
	}

	// Get user progress for today
	progress, err := h.userProgressRepo.GetByUserAndDate(userID, now)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get progress")
		return
	}

	// If no progress exists, create a new one
	if progress == nil {
		progress = &models.UserProgress{
			UserID:            userID,
			ReadingPlanID:     plan.ID,
			Date:              now,
			MorningCompleted:  false,
			EveningCompleted:  false,
		}
	}

	response := TodayReadingsResponse{
		Period:    period,
		Readings:  plan,
		Progress:  progress,
		DayOfYear: dayOfYear,
	}

	middleware.RespondJSON(w, http.StatusOK, response)
}

func (h *ReadingsHandler) MarkCompleted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		middleware.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		middleware.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req MarkCompletedRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Period != "morning" && req.Period != "evening" {
		middleware.RespondError(w, http.StatusBadRequest, "Period must be 'morning' or 'evening'")
		return
	}

	now := time.Now()
	dayOfYear := now.YearDay()

	// Get reading plan for today
	plan, err := h.readingPlanRepo.GetByDayOfYear(dayOfYear)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get reading plan")
		return
	}

	if plan == nil {
		middleware.RespondError(w, http.StatusNotFound, "Reading plan not found for today")
		return
	}

	// Get or create progress
	progress, err := h.userProgressRepo.GetByUserAndDate(userID, now)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get progress")
		return
	}

	if progress == nil {
		progress = &models.UserProgress{
			UserID:            userID,
			ReadingPlanID:     plan.ID,
			Date:              now,
			MorningCompleted:  false,
			EveningCompleted:  false,
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
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to save progress")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, progress)
}

func (h *ReadingsHandler) GetProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		middleware.RespondError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	userID, err := middleware.GetUserIDFromContext(r)
	if err != nil {
		middleware.RespondError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	progresses, err := h.userProgressRepo.GetUserProgress(userID)
	if err != nil {
		middleware.RespondError(w, http.StatusInternalServerError, "Failed to get progress")
		return
	}

	middleware.RespondJSON(w, http.StatusOK, progresses)
}


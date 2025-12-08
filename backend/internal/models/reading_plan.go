package models

type ReadingPlan struct {
	ID                int    `json:"id"`
	DayOfYear         int    `json:"day_of_year"`
	OldTestamentRef   string `json:"old_testament_ref"`
	NewTestamentRef   string `json:"new_testament_ref"`
	PsalmsRef         string `json:"psalms_ref"`
	ProverbsRef       string `json:"proverbs_ref"`
}


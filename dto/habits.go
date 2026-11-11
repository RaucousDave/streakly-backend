package dto

type HabitType string

const (
	HabitTypeBoolean      HabitType = "BOOLEAN"
	HabitTypeQuantitative HabitType = "QUANTITATIVE"
)

type HabitDto struct {
	Title     string    `json:"title"`
	Type      HabitType `json:"type"`
	TargetQty *int      `json:"target_qty"`
	MVTQty    *int      `json:"mvt_qty"`
	Unit      *string   `json:"unit"`
}

type HabitCompletionDto struct {
	HabitId        string `json:"habit_id"`
	CompletedDate  string `json:"completed_date"`
	CompletedAtUTC string `json:"completed_at_utc"`
	QtyAchieved    int    `json:"qty_achieved"`
	IdempotencyKey string `json:"idempotency_key"`
}

type EditHabitDto struct {
	Title     *string    `json:"title"`
	Type      *HabitType `json:"type"`
	TargetQty *int       `json:"target_qty"`
	MVTQty    *int       `json:"mvt_qty"`
	Unit      *string    `json:"unit"`
}
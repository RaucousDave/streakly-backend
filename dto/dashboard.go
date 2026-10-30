package dto

type HabitDashboardItem struct {
	ID               string  `json:"id"`
	Title            string  `json:"title"`
	Type             string  `json:"type"`
	TargetQty        *int    `json:"target_qty"`
	MVTQty           *int    `json:"mvt_qty"`
	Unit             *string `json:"unit"`
	CurrentStreak    int     `json:"current_streak"`
	IsCompletedToday bool    `json:"is_completed_today"`
	QtyAchievedToday int     `json:"qty_achieved_today"`
	IsMVTMetToday    bool    `json:"is_mvt_met_today"`
	IsTargetMetToday bool    `json:"is_target_met_today"`
}

type HeatmapEntry struct {
	Date           string `json:"date"`
	Completed      bool   `json:"completed"`
	Total          int    `json:"total"`
	CompletedCount int    `json:"completed_count"`
}

type DashboardUser struct {
	DisplayName   string `json:"display_name"`
	Timezone      string `json:"timezone"`
	FreezeBalance int    `json:"freeze_balance"`
}

type TodaySection struct {
	Date   string               `json:"date"`
	Habits []HabitDashboardItem `json:"habits"`
}

type DashboardResponse struct {
	User    DashboardUser  `json:"user"`
	Today   TodaySection   `json:"today"`
	Heatmap []HeatmapEntry `json:"heatmap_entry"`
}

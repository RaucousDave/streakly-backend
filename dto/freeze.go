package dto

type ManualFreezeResponse struct {
	Success       bool   `json:"success"`
	FreezeBalance int    `json:"freeze_balance"`
	TargetDate    string `json:"target_date"`
	Reason        string `json:"reason"`
}

type PerformanceFreezeResponse struct {
	Success        bool   `json:"success"`
	Qualified      bool   `json:"qualified"`
	CompletionRate int    `json:"completion_rate"`
	FreezeGranted  bool   `json:"freeze_granted"`
	AlreadyGranted bool   `json:"already_granted"`
	Reason         string `json:"reason"`
	FreezeBalance  *int   `json:"freeze_balance"`
}

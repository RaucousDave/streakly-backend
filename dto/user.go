package dto

type UserResponseDto struct {
	Email       string  `json:"email"`
	GoogleId    *string `json:"google_id"`
	DisplayName string  `json:"display_name"`
	AvatarUrl   *string `json:"avatar_url"`
	Provider    string  `json:"provider"`
}

package dto

type GoogleLoginDto struct {
	IDToken  string `json:"id_token"`
	Timezone string `json:"timezone"`
}

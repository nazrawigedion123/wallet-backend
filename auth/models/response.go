package models

type RegisterResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Tier  string `json:"tier"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  struct {
		ID    uint   `json:"id"`
		Email string `json:"email"`
		Tier  string `json:"tier"`
	} `json:"user"`
}

type LogoutResponse struct {
	Message string `json:"message"`
}

type ProfileResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Tier  string `json:"tier"`
}

type TierUpgradeResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Tier  string `json:"tier"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

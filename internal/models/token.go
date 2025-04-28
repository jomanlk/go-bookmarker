package models

// Token represents an access token for a user
type Token struct {
	ID        int64  `json:"id"`
	UserID    int64  `json:"user_id"`
	Token     string `json:"token"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
	ExpiresAt int64  `json:"expires_at"`
}

package models

// RefreshToken represents a refresh token for a user
type RefreshToken struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"user_id"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt   int64  `json:"expires_at"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
}

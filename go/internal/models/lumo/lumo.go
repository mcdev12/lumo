package lumo

import (
	"time"

	"github.com/google/uuid"
)

// Lumo represents a container for travel planning in the domain
type Lumo struct {
	// Internal database ID (not exposed in API)
	ID int64 `json:"-"`

	// Unique identifier (UUID)
	LumoID string `json:"lumo_id"`

	// User identifier (UUID)
	UserID string `json:"user_id"`

	// Display title
	Title string `json:"title"`

	// System timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewLumo creates a new Lumo with generated UUID
func NewLumo(userID string, title string) *Lumo {
	now := time.Now()
	return &Lumo{
		LumoID:    uuid.New().String(),
		UserID:    userID,
		Title:     title,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// IsValid performs basic validation on the Lumo
func (l *Lumo) IsValid() bool {
	if l.LumoID == "" {
		return false
	}
	if l.UserID == "" {
		return false
	}
	if l.Title == "" {
		return false
	}
	return true
}
package file_entity

import "time"

type ActivityLog struct {
	ID        uint      `json:"id"`
	UserID    string    `json:"user_id"`
	Action    string    `json:"action"` // "upload" or "delete"
	FileID    *string   `json:"file_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type ActivityByDay struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

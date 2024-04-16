package domain

import "time"

type Order struct {
	ID         int64
	Number     int64 `json:"number,omitempty"`
	UserId     int64
	Status     string    `json:"status,omitempty"`
	Accrual    *int64    `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

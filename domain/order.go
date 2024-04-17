package domain

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID         int64
	Number     int64 `json:"number,omitempty"`
	UserId     int64
	Status     string    `json:"status,omitempty"`
	Accrual    *int64    `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		UploadedAt string `json:"uploaded_at"`
		*Order
	}{
		Order:      o,
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}

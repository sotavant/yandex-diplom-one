package domain

import (
	"encoding/json"
	"time"
)

type Order struct {
	ID         int64     `json:"-"`
	Number     string    `json:"number,omitempty"`
	UserID     int64     `json:"-"`
	Status     string    `json:"status,omitempty"`
	Accrual    *float64  `json:"accrual,omitempty"`
	UploadedAt time.Time `json:"uploaded_at"`
}

func (o *Order) MarshalJSON() ([]byte, error) {
	type Alias Order
	return json.Marshal(&struct {
		UploadedAt string `json:"uploaded_at"`
		*Alias
	}{
		Alias:      (*Alias)(o),
		UploadedAt: o.UploadedAt.Format(time.RFC3339),
	})
}

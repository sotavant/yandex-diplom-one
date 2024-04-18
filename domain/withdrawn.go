package domain

import (
	"encoding/json"
	"time"
)

type Withdrawn struct {
	ID          int64     `json:"-"`
	OrderNum    int64     `json:"order,omitempty"`
	UserId      int64     `json:"-"`
	Sum         float64   `json:"sum,omitempty"`
	ProcessedAt time.Time `json:"uploaded_at"`
}

func (o *Withdrawn) MarshalJSON() ([]byte, error) {
	type Alias Withdrawn
	return json.Marshal(&struct {
		ProcessedAt string `json:"processed_at"`
		*Alias
	}{
		Alias:       (*Alias)(o),
		ProcessedAt: o.ProcessedAt.Format(time.RFC3339),
	})
}

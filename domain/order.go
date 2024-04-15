package domain

type Order struct {
	ID         int64
	Number     int64 `json:"number"`
	UserId     int64
	Status     string `json:"status"`
	Accrual    int64  `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

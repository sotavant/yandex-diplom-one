package domain

type User struct {
	ID        int64   `json:"-"`
	Login     string  `json:"login,omitempty"`
	Password  string  `json:"password,omitempty"`
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

package model

type Location struct {
	ID       int64  `json:"id" db:"id"`
	Code     string `json:"code" db:"code"`
	Name     string `json:"name" db:"name"`
	Capacity int64  `json:"capacity" db:"capacity"`
}

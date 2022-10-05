package entities

import "time"

type Fileprocess struct {
	Id        int64
	Periode   string
	File      string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

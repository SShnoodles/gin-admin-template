package models

import "time"

type Org struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"username"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Org) TableName() string {
	return "org"
}

package models

import "time"

type Org struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:50"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Org) TableName() string {
	return "org"
}

package models

import "time"

type Resource struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:50"`
	Code      string    `json:"code" gorm:"not null;size:50"`
	Url       string    `json:"path" gorm:"not null;size:200"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Resource) TableName() string {
	return "resource"
}

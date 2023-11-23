package models

import "time"

type Resource struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	Url       string    `json:"path"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Resource) TableName() string {
	return "resource"
}

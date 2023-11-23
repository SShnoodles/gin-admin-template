package models

import "time"

type Role struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Code      string    `json:"code"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Role) TableName() string {
	return "role"
}

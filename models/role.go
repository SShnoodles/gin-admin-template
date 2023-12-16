package models

import "time"

type Role struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null;size:50"`
	Code      string    `json:"code" gorm:"not null;size:50"`
	OrgId     int64     `json:"orgId" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Role) TableName() string {
	return "role"
}

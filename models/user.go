package models

import (
	"time"
)

type User struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	OrgId     int64     `json:"orgId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}

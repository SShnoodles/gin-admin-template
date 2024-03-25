package domain

import (
	"time"
)

type User struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	Username  string    `json:"username" gorm:"not null;size:50"`
	RealName  string    `json:"realName" gorm:"not null;size:50"`
	WorkNo    string    `json:"workNo" gorm:"not null;size:50"`
	Password  string    `json:"password" gorm:"not null;size:200"`
	OrgId     int64     `json:"orgId,string"`
	Enabled   bool      `json:"enabled" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}

package domain

import "time"

type Role struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	Name      string    `json:"name" gorm:"not null;size:50"`
	Code      string    `json:"code" gorm:"not null;size:50"`
	OrgId     int64     `json:"orgId,string" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Role) TableName() string {
	return "role"
}

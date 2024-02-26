package domain

import "time"

type Resource struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	Name      string    `json:"name" gorm:"not null;size:50"`
	Method    string    `json:"method" gorm:"not null;size:20"`
	Path      string    `json:"path" gorm:"not null;size:200"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Resource) TableName() string {
	return "resource"
}

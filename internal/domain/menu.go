package domain

import "time"

type Menu struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	Pid       *int64    `json:"pid,string"`
	Name      string    `json:"name" gorm:"not null;size:50"`
	Path      string    `json:"path" gorm:"not null;size:200"`
	Icon      string    `json:"icon" gorm:"not null;size:200"`
	Sort      int32     `json:"sort" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Menu) TableName() string {
	return "menu"
}

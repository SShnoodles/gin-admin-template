package models

import "time"

type Menu struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	Pid       int64     `json:"pid"`
	Name      string    `json:"name"`
	Path      string    `json:"path"`
	Icon      string    `json:"icon"`
	Target    string    `json:"target"`
	Sort      int32     `json:"sort"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (Menu) TableName() string {
	return "menu"
}

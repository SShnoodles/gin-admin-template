package models

import "time"

type MenuResourceRelation struct {
	Id         int64     `json:"id" gorm:"primaryKey"`
	MenuId     int64     `json:"menuId"`
	ResourceId int64     `json:"resourceId"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (MenuResourceRelation) TableName() string {
	return "menu_resource_relation"
}

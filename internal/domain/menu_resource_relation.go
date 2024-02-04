package domain

import "time"

type MenuResourceRelation struct {
	Id         int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	MenuId     int64     `json:"menuId,string"`
	ResourceId int64     `json:"resourceId,string"`
	CreatedAt  time.Time `json:"createdAt"`
}

func (MenuResourceRelation) TableName() string {
	return "menu_resource_relation"
}

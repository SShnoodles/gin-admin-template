package domain

import "time"

type RoleMenuRelation struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	RoleId    int64     `json:"roleId"`
	MenuId    int64     `json:"menuId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (RoleMenuRelation) TableName() string {
	return "role_menu_relation"
}

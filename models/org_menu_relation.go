package models

import "time"

type OrgMenuRelation struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	OrgId     int64     `json:"orgId"`
	MenuId    int64     `json:"menuId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (OrgMenuRelation) TableName() string {
	return "org_menu_relation"
}

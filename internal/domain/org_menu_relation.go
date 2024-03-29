package domain

import "time"

type OrgMenuRelation struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	OrgId     int64     `json:"orgId,string"`
	MenuId    int64     `json:"menuId,string"`
	CreatedAt time.Time `json:"createdAt"`
}

func (OrgMenuRelation) TableName() string {
	return "org_menu_relation"
}

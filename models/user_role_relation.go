package models

import "time"

type UserRoleRelation struct {
	Id        int64     `json:"id" gorm:"primaryKey"`
	UserId    int64     `json:"userId"`
	OrgId     int64     `json:"orgId"`
	RoleId    int64     `json:"roleId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserRoleRelation) TableName() string {
	return "user_role_relation"
}

package domain

import "time"

type UserRoleRelation struct {
	Id        int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	UserId    int64     `json:"userId,string"`
	OrgId     int64     `json:"orgId,string"`
	RoleId    int64     `json:"roleId,string"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserRoleRelation) TableName() string {
	return "user_role_relation"
}

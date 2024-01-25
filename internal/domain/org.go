package domain

import "time"

type Org struct {
	Id         int64     `json:"id,string" gorm:"primaryKey;autoIncrement:false"`
	Name       string    `json:"name" gorm:"not null;size:50"`
	CreditCode string    `json:"creditCode" gorm:"not null;size:18"`
	Address    string    `json:"address,omitempty" gorm:"size:300"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (Org) TableName() string {
	return "org"
}

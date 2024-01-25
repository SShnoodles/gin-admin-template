package service

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
)

func FindOrgByName(name string) (domain.Org, error) {
	var org domain.Org
	err := config.DB.First(&org, "name = ?", org).Error
	if err != nil {
		return org, err
	}
	return org, nil
}

package service

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
)

func FindMenuByPath(path string) (domain.Menu, error) {
	var menu domain.Menu
	err := config.DB.First(&menu, "path = ?", path).Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func FindRootMenu() ([]domain.Menu, error) {
	var menu []domain.Menu
	err := config.DB.Find(&menu, "pid is null or pid = ''").Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

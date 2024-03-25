package service

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"strconv"
)

func FindMenuByPath(path string) (domain.Menu, error) {
	var menu domain.Menu
	err := config.DB.First(&menu, "path = ?", path).Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func FindMenusByPid(pid int64) ([]domain.Menu, error) {
	var menu []domain.Menu
	err := config.DB.Find(&menu, "pid = ?", pid).Error
	if err != nil {
		return menu, err
	}
	return menu, nil
}

func DeleteMenusByPid(pid int64) error {
	return config.DB.Delete(&domain.Menu{}, "pid = ?", pid).Error
}

func FindResourceIdsByMenuId(id int64) ([]string, error) {
	var mr []domain.MenuResourceRelation
	err := config.DB.Find(&mr, "menu_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	if len(mr) == 0 {
		return ids, nil
	}
	for _, m := range mr {
		ids = append(ids, strconv.FormatInt(m.ResourceId, 10))
	}
	return ids, nil
}

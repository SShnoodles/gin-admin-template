package service

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"strconv"
)

func FindMenuIdsByRoleId(id int64) ([]string, error) {
	var om []domain.RoleMenuRelation
	err := config.DB.Find(&om, "role_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	if len(om) == 0 {
		return ids, nil
	}
	for _, m := range om {
		ids = append(ids, strconv.FormatInt(m.MenuId, 10))
	}
	return ids, nil
}

func FindRolesByOrgId(orgId int64) ([]domain.Role, error) {
	var roles []domain.Role
	err := config.DB.Find(&roles, "org_id = ?", orgId).Error
	if err != nil {
		return roles, err
	}
	return roles, nil
}

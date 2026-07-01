package service

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrRoleCodeExists = errors.New("role code exists")
	ErrRoleNotFound   = errors.New("role not found")
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

func CreateRole(role domain.Role, menuIds []string) (int64, error) {
	if err := validateRole(role); err != nil {
		return 0, err
	}
	roleId := config.IdGenerate()
	return roleId, config.DB.Transaction(func(tx *gorm.DB) error {
		role.Id = roleId
		if err := tx.Create(&role).Error; err != nil {
			return err
		}
		return replaceRoleMenus(tx, roleId, menuIds)
	})
}

func UpdateRole(roleId int64, input domain.Role, menuIds []string) error {
	if err := validateRole(input); err != nil {
		return err
	}
	var role domain.Role
	err := FindById(&role, roleId)
	if err != nil {
		return ErrRoleNotFound
	}
	if input.Code != role.Code {
		var existing domain.Role
		err = FindByCode(&existing, input.Code)
		if err == nil {
			return ErrRoleCodeExists
		}
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		role.Name = input.Name
		role.Code = input.Code
		role.OrgId = input.OrgId
		if err := tx.Save(&role).Error; err != nil {
			return err
		}
		return replaceRoleMenus(tx, roleId, menuIds)
	})
}

func DeleteRole(id int64) error {
	var role domain.Role
	err := FindById(&role, id)
	if err != nil {
		return ErrRoleNotFound
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&role).Error; err != nil {
			return err
		}
		if err := tx.Where("role_id = ?", id).Delete(&domain.RoleMenuRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func replaceRoleMenus(tx *gorm.DB, roleId int64, menuIds []string) error {
	if err := tx.Where("role_id = ?", roleId).Delete(&domain.RoleMenuRelation{}).Error; err != nil {
		return err
	}
	if len(menuIds) == 0 {
		return nil
	}
	ids, err := ParsePositiveIds(menuIds)
	if err != nil {
		return err
	}
	var rmr []domain.RoleMenuRelation
	for _, menuId := range ids {
		rmr = append(rmr, domain.RoleMenuRelation{
			Id:     config.IdGenerate(),
			RoleId: roleId,
			MenuId: menuId,
		})
	}
	return tx.Create(&rmr).Error
}

func validateRole(role domain.Role) error {
	if strings.TrimSpace(role.Name) == "" || strings.TrimSpace(role.Code) == "" || role.OrgId <= 0 {
		return ErrInvalidParam
	}
	return nil
}

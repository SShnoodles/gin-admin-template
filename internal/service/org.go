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
	ErrOrgNameExists = errors.New("org name exists")
	ErrOrgNotFound   = errors.New("org not found")
)

func FindMenuIdsByOrgId(id int64) ([]string, error) {
	var om []domain.OrgMenuRelation
	err := config.DB.Find(&om, "org_id = ?", id).Error
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

func CreateOrg(org domain.Org, menuIds []string) (int64, error) {
	if err := validateOrg(org); err != nil {
		return 0, err
	}
	var existing domain.Org
	err := FindByName(&existing, org.Name)
	if err == nil {
		return 0, ErrOrgNameExists
	}

	orgId := config.IdGenerate()
	return orgId, config.DB.Transaction(func(tx *gorm.DB) error {
		org.Id = orgId
		if err := tx.Create(&org).Error; err != nil {
			return err
		}
		return replaceOrgMenus(tx, orgId, menuIds)
	})
}

func UpdateOrg(orgId int64, input domain.Org, menuIds []string) error {
	if err := validateOrg(input); err != nil {
		return err
	}
	var org domain.Org
	err := FindById(&org, orgId)
	if err != nil {
		return ErrOrgNotFound
	}
	if input.Name != org.Name {
		var existing domain.Org
		err = FindByName(&existing, input.Name)
		if err == nil {
			return ErrOrgNameExists
		}
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		org.Name = input.Name
		org.CreditCode = input.CreditCode
		org.Address = input.Address
		if err := tx.Save(&org).Error; err != nil {
			return err
		}
		return replaceOrgMenus(tx, orgId, menuIds)
	})
}

func DeleteOrg(id int64) error {
	var org domain.Org
	err := FindById(&org, id)
	if err != nil {
		return ErrOrgNotFound
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&org).Error; err != nil {
			return err
		}
		if err := tx.Where("org_id = ?", id).Delete(&domain.OrgMenuRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func replaceOrgMenus(tx *gorm.DB, orgId int64, menuIds []string) error {
	if err := tx.Where("org_id = ?", orgId).Delete(&domain.OrgMenuRelation{}).Error; err != nil {
		return err
	}
	if len(menuIds) == 0 {
		return nil
	}
	ids, err := ParsePositiveIds(menuIds)
	if err != nil {
		return err
	}
	var omr []domain.OrgMenuRelation
	for _, menuId := range ids {
		omr = append(omr, domain.OrgMenuRelation{
			Id:     config.IdGenerate(),
			OrgId:  orgId,
			MenuId: menuId,
		})
	}
	return tx.Create(&omr).Error
}

func validateOrg(org domain.Org) error {
	if strings.TrimSpace(org.Name) == "" || strings.TrimSpace(org.CreditCode) == "" {
		return ErrInvalidParam
	}
	return nil
}

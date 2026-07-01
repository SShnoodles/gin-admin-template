package service

import (
	"errors"
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/util"
	"strconv"

	"gorm.io/gorm"
)

var (
	ErrUserExists        = errors.New("user exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserPasswordWrong = errors.New("user password wrong")
)

func FindUserByUsername(username string) (domain.User, error) {
	var user domain.User
	err := config.DB.First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func FindRoleIdsByUserId(id int64) ([]string, error) {
	var ur []domain.UserRoleRelation
	err := config.DB.Find(&ur, "user_id = ?", id).Error
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0)
	if len(ur) == 0 {
		return ids, nil
	}
	for _, m := range ur {
		ids = append(ids, strconv.FormatInt(m.RoleId, 10))
	}
	return ids, nil
}

func CreateUser(user domain.User, roleIds []string) (int64, error) {
	existing, _ := FindUserByUsername(user.Username)
	if existing != (domain.User{}) {
		return 0, ErrUserExists
	}

	userId := config.IdGenerate()
	password, err := util.EncryptedPassword(util.DefaultPassword)
	if err != nil {
		return 0, err
	}
	return userId, config.DB.Transaction(func(tx *gorm.DB) error {
		user.Id = userId
		user.Password = password
		user.Enabled = true
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		return replaceUserRoles(tx, userId, user.OrgId, roleIds)
	})
}

func UpdateUserInfo(userId int64, input domain.User, roleIds []string) error {
	var user domain.User
	err := FindById(&user, userId)
	if err != nil {
		return ErrUserNotFound
	}
	if user.Username != input.Username {
		_, err := FindUserByUsername(input.Username)
		if err == nil {
			return ErrUserExists
		}
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		user.Username = input.Username
		user.RealName = input.RealName
		user.WorkNo = input.WorkNo
		user.OrgId = input.OrgId
		if err := tx.Save(&user).Error; err != nil {
			return err
		}
		return replaceUserRoles(tx, userId, input.OrgId, roleIds)
	})
}

func ToggleUserEnabled(userId int64) error {
	var user domain.User
	err := FindById(&user, userId)
	if err != nil {
		return ErrUserNotFound
	}
	return config.DB.Model(&user).UpdateColumn("enabled", !user.Enabled).Error
}

func ChangeUserPassword(userId int64, oldPassword string, newPassword string) error {
	var user domain.User
	err := FindById(&user, userId)
	if err != nil {
		return ErrUserNotFound
	}
	if !util.VerifyPassword(oldPassword, user.Password) {
		return ErrUserPasswordWrong
	}
	password, err := util.EncryptedPassword(newPassword)
	if err != nil {
		return err
	}
	return config.DB.Model(&user).UpdateColumn("password", password).Error
}

func DeleteUser(id int64) error {
	var user domain.User
	err := FindById(&user, id)
	if err != nil {
		return ErrUserNotFound
	}
	return config.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", id).Delete(&domain.UserRoleRelation{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func replaceUserRoles(tx *gorm.DB, userId int64, orgId int64, roleIds []string) error {
	if err := tx.Where("user_id = ?", userId).Delete(&domain.UserRoleRelation{}).Error; err != nil {
		return err
	}
	if len(roleIds) == 0 {
		return nil
	}
	var urr []domain.UserRoleRelation
	for _, id := range roleIds {
		roleId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return err
		}
		urr = append(urr, domain.UserRoleRelation{
			Id:     config.IdGenerate(),
			UserId: userId,
			RoleId: roleId,
			OrgId:  orgId,
		})
	}
	return tx.Create(&urr).Error
}

func InitAdminUser() error {
	var user domain.User
	err := config.DB.First(&user, "username = ?", "superadmin").Error
	if err != nil {
		if err.Error() == "record not found" {
			password, err := util.EncryptedPassword(util.DefaultPassword)
			if err != nil {
				return err
			}
			user = domain.User{
				Id:       config.IdGenerate(),
				Username: "superadmin",
				Password: password,
				RealName: "Administrator",
				Enabled:  true,
			}
			err = config.DB.Create(&user).Error
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	if !user.Enabled {
		return config.DB.Model(&user).UpdateColumn("enabled", true).Error
	}
	return nil
}

package service

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/util"
	"strconv"
)

func FindUserByUsername(username string) (domain.User, error) {
	var user domain.User
	err := config.DB.First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdateUser(user domain.User, updatePassword bool) error {
	var err error
	if updatePassword {
		user.Password, err = util.EncryptedPassword(user.Password)
		if err != nil {
			return err
		}
	}
	err = config.DB.Model(user).Updates(&user).Error
	return err
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

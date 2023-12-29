package service

import (
	"gin-admin-template/internal/config"
	"gin-admin-template/internal/domain"
	"gin-admin-template/internal/util"
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

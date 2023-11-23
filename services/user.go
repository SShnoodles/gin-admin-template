package services

import (
	"gin-admin-template/initializations"
	"gin-admin-template/models"
	"gin-admin-template/utils"
)

func FindUserByUsername(username string) (models.User, error) {
	var user models.User
	err := initializations.DB.First(&user, "username = ?", username).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func UpdateUser(user models.User, updatePassword bool) error {
	var err error
	if updatePassword {
		user.Password, err = utils.EncryptedPassword(user.Password)
		if err != nil {
			return err
		}
	}
	err = initializations.DB.Model(user).Updates(&user).Error
	return err
}

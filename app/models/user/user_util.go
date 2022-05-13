package user

import "GDColumn/pkg/database"

func IsEmailExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("email = ?", email).Count(&count)
	return count > 0
}

func IsPhoneExist(email string) bool {
	var count int64
	database.DB.Model(User{}).Where("phone = ?", email).Count(&count)
	return count > 0
}
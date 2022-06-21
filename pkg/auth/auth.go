package auth

import (
	"GDColumn/app/models/user"
	"GDColumn/pkg/logger"
	"errors"
	"github.com/gin-gonic/gin"
)

// Attempt 尝试登录
func Attempt(email, password string)(user.User,error){
	userModel := user.GetByMulti(email)
	if userModel.ID == 0{
		return user.User{},errors.New("用户不存在")
	}
	if !userModel.ComparePassword(password){
		return user.User{},errors.New("密码错误")
	}
	return userModel,nil
}

func LoginByPhone(phone string) (user.User, error) {
	userModel := user.GetByPhone(phone)
	if userModel.ID == 0 {
		return user.User{}, errors.New("手机号未注册")
	}

	return userModel, nil
}

// CurrentUser 从 gin.context 中获取当前登录用户
func CurrentUser(c *gin.Context) user.User {
	userModel, ok := c.MustGet("current_user").(user.User)
	if !ok {
		logger.LogIf(errors.New("无法获取用户"))
		return user.User{}
	}
	// db is now a *DB value
	return userModel
}

// CurrentUserAvatar 从 gin.context 中获取当前登录用户头像
func CurrentUserAvatar(c *gin.Context) user.Avatar {
	avatarModel, ok := c.MustGet("current_avatar").(user.Avatar)
	if !ok {
		logger.LogIf(errors.New("无法获取用户头像"))
		return user.Avatar{}
	}
	// db is now a *DB value
	return avatarModel
}

// CurrentUID 从 gin.context 中获取当前登录用户 ID
func CurrentUID(c *gin.Context) string {
	return c.GetString("current_user_id")
}
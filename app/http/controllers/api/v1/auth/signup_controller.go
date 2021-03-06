package auth

import (
	v1 "GDColumn/app/http/controllers/api/v1"
	"GDColumn/app/models/user"
	"GDColumn/app/models/column"
	"GDColumn/app/requests"
	"GDColumn/pkg/jwt"
	"GDColumn/pkg/logger"
	"GDColumn/pkg/response"
	"GDColumn/pkg/snowflake"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type SignupController struct {
	v1.BaseAPIController
}

func (sc *SignupController) IsPhoneExist(c *gin.Context) {
	request := requests.SignupPhoneExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupPhoneExist); !ok {
		return
	}

	response.JSON(c, gin.H{
		"exist": user.IsPhoneExist(request.Phone),
	})
}

func (sc *SignupController) IsEmailExist(c *gin.Context) {
	request := requests.SignupEmailExistRequest{}
	if ok := requests.Validate(c, &request, requests.ValidateSignupEmailExist); !ok {
		return
	}
	response.JSON(c, gin.H{
		"exist": user.IsEmailExist(request.Email),
	})
}

func (sc *SignupController) SignupUsingPhone(c *gin.Context) {

	request := requests.SignupUsingPhoneRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingPhone); !ok {
		return
	}

	uid, err := snowflake.GetID()
	userID := cast.ToString(uid)
	if err != nil {
		logger.ErrorString("Auth","GetID",err.Error())
		response.Abort500(c, "创建用户失败，请稍后尝试~")
		return
	}
	userModel := user.User{
		NickName:     request.NickName,
		Phone:    	  request.Phone,
		Password: 	  request.Password,
		ID:		      userID,
	}
	userModel.Create()

	if userModel.ID != "" {
		token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.NickName)
		response.CreatedJSON(c, gin.H{
			"token": token,
			"data":  userModel,
		})
	} else {
		response.Abort500(c, "创建用户失败，请稍后尝试~")
	}
}

// SignupUsingEmail 注册用户
// @Summary 注册用户
// @Description 需要用户名密码
// @Tags user 关于用户的路由，登录，注册，获取当前用户等等
// @Accept application/json
// @Produce application/json
// @Param body body requests.SignupUsingEmailRequest true "用户注册，需要提供用户的邮箱和密码"
// @Success 200 {object} user.User
//@Router /auth/signup/using-email [POST]
func (sc *SignupController) SignupUsingEmail(c *gin.Context) {

	request := requests.SignupUsingEmailRequest{}
	if ok := requests.Validate(c, &request, requests.SignupUsingEmail); !ok {
		return
	}

	uid, err := snowflake.GetID()
	userID := cast.ToString(uid)
	if err != nil {
		logger.ErrorString("Auth","GetID",err.Error())
		response.Abort500(c, "创建用户失败，请稍后尝试~")
		return
	}

	cid, err := snowflake.GetID()
	columnID := cast.ToString(cid)
	if err != nil {
		logger.ErrorString("Auth","GetID",err.Error())
		response.Abort500(c, "创建专栏失败，请稍后尝试~")
		return
	}
	columnModel := column.Column{
		ID:                    columnID,
		Title:                 fmt.Sprintf("这里是%v专栏，有一段非常有意思的简介，可以更新一下欧",request.NickName),
		Description:           fmt.Sprintf("%v的专栏",request.NickName),
		Author:                userID,
	}
	columnModel.Create()

	userModel := user.User{
		NickName:     request.NickName,
		Email:    	  request.Email,
		Password:     request.Password,
		ID:           userID,
		ColumnID:     columnID,
	}
	userModel.Create()

	if userModel.ID != "" {
		//token := jwt.NewJWT().IssueToken(userModel.GetStringID(), userModel.NickName)
		response.Created(c,userModel)
	} else {
		response.Abort500(c, "创建用户失败，请稍后尝试~")
	}
}
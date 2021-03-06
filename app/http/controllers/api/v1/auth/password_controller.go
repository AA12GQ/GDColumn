package auth

import (
	v1 "GDColumn/app/http/controllers/api/v1"
	"GDColumn/app/policies"
	"GDColumn/app/requests"
	"GDColumn/pkg/auth"
	"GDColumn/pkg/response"
	"GDColumn/app/models/user"
	"github.com/gin-gonic/gin"
)

type PasswordController struct {
	v1.BaseAPIController
}

func (pc *PasswordController) ResetByEmail(c *gin.Context) {
	request := requests.ResetByEmailRequest{}
	if ok := requests.Validate(c, &request, requests.ResetByEmail); !ok {
		return
	}

	id := auth.CurrentUID(c)
	if ok := policies.CanModifyPost(c, id); !ok {
		response.Abort403(c)
		return
	}

	// 2. 更新密码
	userModel := user.GetByEmail(request.Email)
	if userModel.ID == "" {
		response.Abort404(c)
	} else {
		userModel.Password = request.Password
		userModel.Save()
		response.Success(c)
	}
}

package v1

import (
    "GDColumn/app/models/image"
    "GDColumn/app/models/user"
    "GDColumn/app/requests"
    "GDColumn/pkg/auth"
    "GDColumn/pkg/response"
    "github.com/gin-gonic/gin"
)


type UsersController struct {
    BaseAPIController
}

func (ctrl *UsersController) CurrentUser(c *gin.Context) {
    userModel := auth.CurrentUser(c)

    currentUser := user.User{
        ID : userModel.ID,
        Email : userModel.Email,
        NickName : userModel.NickName,
        ColumnID : userModel.ColumnID,
    }

    response.Data(c, currentUser)
}

func (ctrl *UsersController) Index(c *gin.Context) {
    request := requests.PaginationRequest{}
    if ok := requests.Validate(c, &request, requests.Pagination); !ok {
        return
    }

    data, pager := user.Paginate(c, 10)
    response.JSON(c, gin.H{
        "data":  data,
        "pager": pager,
    })
}

func (ctrl *UsersController) UpdateProfile(c *gin.Context) {

    request := requests.UserUpdateProfileRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdateProfile); !ok {
        return
    }
    currentUser := auth.CurrentUser(c)
    var avatar *user.Image
    switch  {
    case request.NickName != "":
        currentUser.NickName = request.NickName
        fallthrough
    case request.Description != "":
        currentUser.Description = request.Description
        fallthrough
    case request.AvatarID == "":
        imgModel := image.Get(currentUser.AvatarID)
        avatar = &user.Image{
           ID:  imgModel.ID,
           URL: imgModel.URL,
        }
        currentUser.Avatar = avatar
    case request.AvatarID != "":
        imgModel := image.Get(request.AvatarID)
        avatar = &user.Image{
            ID:  imgModel.ID,
            URL: imgModel.URL,
        }
        currentUser.AvatarID = imgModel.ID
        currentUser.Avatar = avatar
    }

    rowsAffected := currentUser.Updates(currentUser.ID,request.AvatarID,request.NickName,request.Description)
    currentUser = user.Get(currentUser.ID)
    if rowsAffected > 0 {
        response.Data(c, currentUser)
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

func (ctrl *UsersController) UpdateEmail(c *gin.Context) {

    request := requests.UserUpdateEmailRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdateEmail); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    currentUser.Email = request.Email
    rowsAffected := currentUser.Save()

    if rowsAffected > 0 {
        response.Success(c)
    } else {
        // 失败，显示错误提示
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

func (ctrl *UsersController) UpdatePhone(c *gin.Context) {

    request := requests.UserUpdatePhoneRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdatePhone); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    currentUser.Phone = request.Phone
    rowsAffected := currentUser.Save()

    if rowsAffected > 0 {
        response.Success(c)
    } else {
        response.Abort500(c, "更新失败，请稍后尝试~")
    }
}

func (ctrl *UsersController) UpdatePassword(c *gin.Context) {

    request := requests.UserUpdatePasswordRequest{}
    if ok := requests.Validate(c, &request, requests.UserUpdatePassword); !ok {
        return
    }

    currentUser := auth.CurrentUser(c)
    // 验证原始密码是否正确
    _, err := auth.Attempt(currentUser.NickName, request.Password)
    if err != nil {
        response.Unauthorized(c, "原密码不正确")
    } else {
        // 更新密码为新密码
        currentUser.Password = request.NewPassword
        currentUser.Save()

        response.Success(c)
    }
}

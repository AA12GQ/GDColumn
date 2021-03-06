package routes

import (
	controllers "GDColumn/app/http/controllers/api/v1"
	"GDColumn/app/http/controllers/api/v1/auth"
	"GDColumn/app/http/middlewares"
	"GDColumn/pkg/config"
	"github.com/gin-gonic/gin"
	_ "GDColumn/docs"
	gs "github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

func RegisterAPIRoutes(r *gin.Engine) {

	var v1 *gin.RouterGroup
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	v1.Use(middlewares.LimitIP("200-H"))
	{
		authGroup := v1.Group("/auth")
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			suc := new(auth.SignupController)
			authGroup.POST("/signup/phone/exist",
				middlewares.LimitPerRoute("60-H"), middlewares.GuestJWT(), suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist",
				middlewares.LimitPerRoute("60-H"), middlewares.GuestJWT(), suc.IsEmailExist)
			authGroup.POST("/signup/using-phone",
				middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email",
				middlewares.GuestJWT(), suc.SignupUsingEmail)

			lgc := new(auth.LoginController)
			authGroup.POST("/login/using-email", middlewares.GuestJWT(), lgc.LoginByEmail)
			authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), lgc.RefreshToken)

			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-email", middlewares.AuthJWT(), pwc.ResetByEmail)
		}

		uc := new(controllers.UsersController)
		usersGroup := v1.Group("/users")
		v1.GET("/current", middlewares.AuthJWT(),middlewares.Cors, uc.CurrentUser)
		{
			usersGroup.GET("", uc.Index)
			usersGroup.PUT("", middlewares.AuthJWT(), uc.UpdateProfile)
			usersGroup.PUT("/email", middlewares.AuthJWT(), uc.UpdateEmail)
			usersGroup.PUT("/phone", middlewares.AuthJWT(), uc.UpdatePhone)
			usersGroup.PUT("/password", middlewares.AuthJWT(), uc.UpdatePassword)
		}
		imc := new(controllers.ImagesController)
		imcGroup := v1.Group("/upload")
		{
			imcGroup.POST("",middlewares.LimitPerRoute("20-H"),imc.Upload)
		}
		clc := new(controllers.ColumnsController)
		clcGroup := v1.Group("/columns")
		{
			clcGroup.GET("", clc.Index)
			clcGroup.GET("/:id", clc.ShowColumn)
			clcGroup.POST("", middlewares.AuthJWT(), clc.Store)
			clcGroup.PUT("", middlewares.AuthJWT(), clc.Update)
			clcGroup.DELETE("/:id", middlewares.AuthJWT(), clc.Delete)
		}
		poc := new(controllers.PostsController)
		v1.GET("/columns/:id/posts",poc.Index)
		pocGroup := v1.Group("/posts")
		{
			pocGroup.POST("", middlewares.AuthJWT(), poc.Store)
			pocGroup.PUT("/:id",middlewares.AuthJWT(),poc.Update)
			pocGroup.DELETE("/:id", middlewares.AuthJWT(), poc.Delete)
			pocGroup.GET("/:id", poc.Show)
		}
		lsc := new(controllers.LinksController)
		linksGroup := v1.Group("/links")
		{
			linksGroup.GET("", lsc.Index)
		}
	}
}
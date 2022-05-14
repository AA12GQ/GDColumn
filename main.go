package main

import (
	"GDColumn/bootstrap"
	btsConfig "GDColumn/config"
	"GDColumn/pkg/config"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

func init() {
	btsConfig.Initialize()
}

func main() {

	// 配置初始化，依赖命令行 --env 参数
	var env string
	flag.StringVar(&env, "env", "", "加载 .env 文件，如 --env=testing 加载的是 .env.testing 文件")
	flag.Parse()
	config.InitConfig(env)

	bootstrap.SetupLogger()

	gin.SetMode(gin.ReleaseMode)


	router := gin.New()

	bootstrap.SetupDB()

	bootstrap.SetupSnowflake()

	bootstrap.SetupRoute(router)


	err := router.Run(":" + config.Get("app.port"))
	if err != nil {
		fmt.Println(err.Error())
	}
}
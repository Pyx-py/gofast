package initialize

import (
	"github.com/Pyx-py/gofast/global"
	"github.com/Pyx-py/gofast/middleware"
	"github.com/gin-gonic/gin"
	"{{.ModuleName}}/router"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouters(middlewares ...string) {
	Router := gin.Default()
	Router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	group := Router.Group("")
	for _, middlewareName := range middlewares {
		switch middlewareName {
		case "loadtls":
			Router.Use(middleware.LoadTls())
		case "cors":
			Router.Use(middleware.Cors())
		case "error":
			Router.Use(middleware.GinRecovery(true))
		default:

		}
	}
    router.InitHealthCheckRouter(group)
	// router code genarate. **BEGIN !DON'T EDIT IT

	// router code genarate. **END !DON'T EDIT IT

	global.GF_ROUTER = Router
	global.GF_GROUP = group
}

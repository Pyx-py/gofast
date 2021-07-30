package router

import (
    "{{.ModuleName}}/api"
    "github.com/gin-gonic/gin"
)

func InitHealthCheckRouter(Router *gin.RouterGroup) {
    Router.GET("healthCheck", v1.HealthCheck)
}

package global

import (
    "{{.ModuleName}}/config"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GF_DB     *gorm.DB
	GF_LOG    *zap.Logger
    GF_CONF   config.Server
    GF_VP     *viper.Viper
	GF_ROUTER *gin.Engine
	GF_GROUP  *gin.RouterGroup
	GF_REDIS  *redis.Client
)

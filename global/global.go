package global

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	GF_DB     *gorm.DB
	GF_LOG    *zap.Logger
	GF_ROUTER *gin.Engine
	GF_GROUP  *gin.RouterGroup
)

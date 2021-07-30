package initialize

import (
	"fmt"

	"github.com/pyx-py/gofast/initialize/internal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type MysqlConfig struct {
	LogMode      string // gorm日志等级
	LogZap       bool   // 是否将sql日志加载到全局日志文件中
	LogPath      string // 日志路径，为空则不启用默认日志
	Dsn          string
	StringSize   int
	MaxIdleConns int // 空闲连接池中连接的最大数量
	MaxOpenConns int // 打开数据库连接的最大数量
}

func (mc MysqlConfig) initMysqlConfig() (nmc MysqlConfig) {
	if mc.Dsn == "" {
		panic(fmt.Errorf("mysql dsn can't be empty"))
	}
	if mc.LogMode == "" {
		nmc.LogMode = "info"
	} else {
		nmc.LogMode = mc.LogMode
	}
	nmc.LogZap = mc.LogZap
	if mc.StringSize == 0 {
		nmc.StringSize = 255
	} else {
		nmc.StringSize = mc.StringSize
	}
	if mc.MaxIdleConns == 0 {
		nmc.MaxIdleConns = 10
	} else {
		nmc.MaxIdleConns = mc.MaxIdleConns
	}
	if mc.MaxOpenConns == 0 {
		nmc.MaxOpenConns = 100
	} else {
		nmc.MaxOpenConns = mc.MaxOpenConns
	}
	return nmc
}

func (mc MysqlConfig) GormMysql() *gorm.DB {
	nmc := mc.initMysqlConfig()
	mysqlConfig := mysql.Config{
		DSN:                       nmc.Dsn,
		DefaultStringSize:         uint(nmc.StringSize), // string类型字段的默认长度
		DisableDatetimePrecision:  true,                 // 禁用datetime精度， mysql 5.6之前的数据库不支持
		DontSupportRenameIndex:    true,                 // 重命名索引时采用删除并新建的方式，mysql 5.7之前的数据库和mariadb不支持重命名索引
		DontSupportRenameColumn:   true,                 // 用`change`重命名列，mysql 8 之前的数据库和mariadb不支持重命名列
		SkipInitializeWithVersion: false,                // 根据版本自动配置
	}

	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(nmc.LogMode, nmc.LogZap)); err != nil {
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(nmc.MaxIdleConns)
		sqlDB.SetMaxOpenConns(nmc.MaxOpenConns)
		return db
	}

}

func gormConfig(logMode string, logZap bool) *gorm.Config {
	internal.InitDefaultLogger(logZap) // 初始化logger，判断是否将数据库的sql日志写入日志文件
	config := &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}
	switch logMode {
	case "silent", "Silent":
		config.Logger = internal.Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = internal.Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = internal.Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = internal.Default.LogMode(logger.Info)
	default:
		config.Logger = internal.Default.LogMode(logger.Info)
	}
	return config
}

// **INIT_MAIN**  TAG. DO NOT EDIT
package main

import (
	{{- if ne .LogPath ""}}
	"{{.ModuleName}}/core"
	{{- end}}
	"{{.ModuleName}}/global"
	"{{.ModuleName}}/initialize"
	"github.com/fvbock/endless"
	"time"
	"fmt"
)


func main() {
	s := initServer()
	{{if eq .LogPath ""}}
	fmt.Printf("server run on %s", global.GF_CONF.System.Addr)
	{{else}}
	global.GF_LOG.Info(fmt.Sprintf("server run on %s", global.GF_CONF.System.Addr))
	{{end}}
	if err := s.ListenAndServe();err != nil {
		panic(err)
	}

	//程序结束前关闭mysql链接
	db, _ := global.GF_DB.DB()
	defer db.Close()
}

type server interface {
	ListenAndServe() error
}

func init() {
	global.GF_VP = core.Viper()
	{{- if ne .LogPath "" }}
	logConfig := core.LogConfig{
		LogPath:      global.GF_CONF.Zap.LogPath,
		LogLevel:     global.GF_CONF.Zap.Level,
		LogName:      global.GF_CONF.Zap.LogName,
		EncoderLevel: global.GF_CONF.Zap.EncoderLevel,
		TextType:     global.GF_CONF.Zap.TextType,
		Day:          global.GF_CONF.Zap.Day,
		LogInConsole: global.GF_CONF.Zap.LogInConsole,
	}
	global.GF_LOG = logConfig.Zap() // 初始化zap日志库
	{{- else}}
	/* logConfig := core.LogConfig{
		LogPath:      global.GF_CONF.Zap.LogPath,
		LogLevel:     global.GF_CONF.Zap.Level,
		LogName:      global.GF_CONF.Zap.LogName,
		EncoderLevel: global.GF_CONF.Zap.EncoderLevel,
		TextType:     global.GF_CONF.Zap.TextType,
		Day:          global.GF_CONF.Zap.Day
		LogInConsole: global.GF_CONF.Zap.LogInConsole,
	}
	global.GF_LOG = logConfig.Zap() // 初始化zap日志库 */
	{{- end}}
	

	mysqlConfig := initialize.MysqlConfig{
		LogMode: global.GF_CONF.Mysql.LogMode,    // gorm的日志等级
		LogZap:  global.GF_CONF.Mysql.LogZap,    // 是否将sql的详情日志加载到全局日志文件中，默认只会在控制台打印
		LogPath:  global.GF_CONF.Zap.LogPath,    // 日志文件路径，此参数为空sql日志也不会加载到全局日志文件
		Dsn:          global.GF_CONF.Mysql.Dsn(),     // 数据库的连接地址（必填，如: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local）
		StringSize:   global.GF_CONF.Mysql.StringSize,     // 字符类型字段的默认长度
		MaxIdleConns: global.GF_CONF.Mysql.MaxIdleConns,      // 数据库最大连接数
		MaxOpenConns: global.GF_CONF.Mysql.MaxOpenConns,     // 打开数据库连接的最大数量
	}
	global.GF_DB = mysqlConfig.GormMysql() // 初始化数据库连接

	// 初始化数据库表
	models := make([]interface{}, 0)
	/*  models = append(models, &model.System{}) //示例代码,此表无需生成，有需要生成的表在此处添加对应的model，append即可 */
	initialize.MysqlTables(global.GF_DB, models) // 初始化表

	// 初始化redis连接
	/* initialize.Redis() */

	// 初始化router
	initialize.InitRouters()
}

func initServer() server {
	s := endless.NewServer(global.GF_CONF.System.Addr, global.GF_ROUTER)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	return s
}

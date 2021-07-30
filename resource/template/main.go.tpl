// **INIT_MAIN**  TAG. DO NOT EDIT
package main

import (
	{{- if ne .LogPath ""}}
	"github.com/Pyx-py/gofast/core"
	{{- end}}
	"github.com/Pyx-py/gofast/global"
	gf_init "github.com/Pyx-py/gofast/initialize"
	"{{.ModuleName}}/initialize"
	"github.com/fvbock/endless"
	"time"
	"fmt"
)

const (
	address = "127.0.0.1:8888"
)

func main() {
	s := initServer()
	{{if eq .LogPath ""}}
	fmt.Printf("server run on %s", address)
	{{else}}
	global.GF_LOG.Info(fmt.Sprintf("server run on %s", address))
	{{end}}
	if err := s.ListenAndServe();err != nil {
		panic(err)
	}
}

type server interface {
	ListenAndServe() error
}

func init() {
	{{- if ne .LogPath "" }}
	logConfig := core.LogConfig{
		LogPath: "{{.LogPath}}",
		LogLevel:     "",
		LogName:      "",
		EncoderLevel: "",
		TextType:     "",
		Day:          0,
		LogInConsole: true,
	}
	global.GF_LOG = logConfig.Zap() // 初始化zap日志库
	{{- else}}
	/* logConfig := core.LogConfig{
		LogPath:      "",
		LogLevel:     "",
		LogName:      "",
		EncoderLevel: "",
		TextType:     "",
		Day:          0,
		LogInConsole: true,
	}
	global.GF_LOG = logConfig.Zap() // 初始化zap日志库 */
	{{- end}}
	

	mysqlConfig := gf_init.MysqlConfig{
		LogMode: "info",    // gorm的日志等级
		LogZap:  false,    // 是否将sql的详情日志加载到全局日志文件中，默认只会在控制台打印
		LogPath:  {{- if eq .LogPath ""}}""{{- else}}"{{.LogPath}}"{{end}},    // 日志文件路径，此参数为空sql日志也不会加载到全局日志文件
		Dsn:          "",     // 数据库的连接地址（必填，如: user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local）
		StringSize:   255,     // 字符类型字段的默认长度
		MaxIdleConns: 10,      // 数据库最大连接数 （默认10）
		MaxOpenConns: 100,     // 打开数据库连接的最大数量（默认100）
	}
	global.GF_DB = mysqlConfig.GormMysql() // 初始化数据库连接
	// 初始化router
	initialize.InitRouters()
}

func initServer() server {
	s := endless.NewServer(address, global.GF_ROUTER)
	s.ReadHeaderTimeout = 10 * time.Millisecond
	s.WriteTimeout = 10 * time.Second
	return s
}

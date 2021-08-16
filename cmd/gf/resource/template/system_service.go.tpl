package service

import (
	"{{.ModuleName}}/config"
	"{{.ModuleName}}/global"
	"{{.ModuleName}}/model"
	"{{.ModuleName}}/utils"
	"go.uber.org/zap"
)

//@author: [piexlmax](https://github.com/piexlmax)
//@function: GetSystemConfig
//@description: 读取配置文件
//@return: err error, conf config.Server

func GetSystemConfig() (err error, conf config.Server) {
	return nil, global.GF_CONF
}

// @description   set system config,
//@author: [piexlmax](https://github.com/piexlmax)
//@function: SetSystemConfig
//@description: 设置配置文件
//@param: system model.System
//@return: err error

func SetSystemConfig(system model.System) (err error) {
	cs := utils.StructToMap(system.Config)
	for k, v := range cs {
		global.GF_VP.Set(k, v)
	}
	err = global.GF_VP.WriteConfig()
	return err
}

//@author: [SliverHorn](https://github.com/SliverHorn)
//@function: GetServerInfo
//@description: 获取服务器信息
//@return: server *utils.Server, err error

func GetServerInfo() (server *utils.Server, err error) {
	var s utils.Server
	s.Os = utils.InitOS()
	if s.Cpu, err = utils.InitCPU(); err != nil{
        {{- if ne .LogPath ""}}
		global.GF_LOG.Error("func utils.InitCPU() Failed!", zap.String("err", err.Error()))
        {{- end}}
		return &s, err
	}
	if s.Rrm, err = utils.InitRAM(); err != nil{
        {{- if ne .LogPath ""}}
		global.GF_LOG.Error("func utils.InitRAM() Failed!", zap.String("err", err.Error()))
        {{- end}}
		return &s, err
	}
	if s.Disk, err = utils.InitDisk(); err != nil{
        {{- if ne .LogPath ""}}
		global.GF_LOG.Error("func utils.InitDisk() Failed!", zap.String("err", err.Error()))
        {{- end}}
		return &s, err
	}

	return &s, nil
}

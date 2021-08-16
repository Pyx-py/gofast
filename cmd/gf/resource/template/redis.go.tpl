package initialize

import (
	"{{.ModuleName}}/global"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

func Redis() {
	redisCfg := global.GF_CONF.Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisCfg.Addr,
		Password: redisCfg.Password, // no password set
		DB:       redisCfg.DB,       // use default DB
	})
	pong, err := client.Ping().Result()
	if err != nil {
    {{- if ne .LogPath "" }}
		global.GF_LOG.Error("redis connect ping failed, err:", zap.Any("err", err))
    {{- else}}
        panic(fmt.Errorf("redis connect ping failed, err:%s", err.Error()))
    {{- end}}
	} else {
        {{- if ne .LogPath "" }}
		global.GF_LOG.Info("redis connect ping response:", zap.String("pong", pong))
        {{- end}}
		global.GF_REDIS = client
	}
}

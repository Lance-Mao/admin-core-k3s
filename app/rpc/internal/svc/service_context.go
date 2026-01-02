package svc

import (
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/ent"
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/internal/config"
	"github.com/redis/go-redis/v9"

	"github.com/zeromicro/go-zero/core/logx"

	_ "github.com/Lance-Mao/admin-core-k3s/app/rpc/ent/runtime"
)

type ServiceContext struct {
	Config config.Config
	DB     *ent.Client
	Redis  redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	entOpts := []ent.Option{
		ent.Log(logx.Info),
		ent.Driver(c.DatabaseConf.NewNoCacheDriver()),
	}

	if c.DatabaseConf.Debug {
		entOpts = append(entOpts, ent.Debug())
	}

	db := ent.NewClient(entOpts...)

	return &ServiceContext{
		Config: c,
		DB:     db,
		Redis:  c.RedisConf.MustNewUniversalRedis(),
	}
}

package configs

import (
	"github.com/bilibili/kratos/pkg/conf/paladin"
)

// ApplicationConfig ...
type ApplicationConfig struct {
	Peers []string
}

// Application ...
func Application() ApplicationConfig {

	app := ApplicationConfig{}

	if err := paladin.Get("application.toml").UnmarshalTOML(&app); err != nil {
		// 不存在时，将会为nil使用默认配置
		if err != paladin.ErrNotExist {
			panic(err)
		}
	}

	return app
}

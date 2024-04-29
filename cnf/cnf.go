package cnf

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Http struct {
			Port int
		}
	}
	Redis struct {
		Addr     string
		Username string
		Password string
		DB       int
	}
}

func init() {
	viper.SetDefault("Server.Http.Port", 8080)
	viper.SetDefault("Redis.Addr", "localhost:6379")
	viper.SetDefault("Redis.Username", "default")
	viper.SetDefault("Redis.Password", "")
	viper.SetDefault("Redis.DB", 0)
}

var cnf = new(Config)

func Get() *Config {
	return cnf
}

func InitConfig(config interface{}) {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("unmarshal config failed err: %s", err.Error()))
	}
}

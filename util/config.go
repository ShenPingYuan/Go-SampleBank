package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER" default:"mysql"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS" default:":8080"`
	SymmetricKey        string        `mapstructure:"SYMMETRIC_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION" default:"15m"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	// 设置配置文件名
	viper.SetConfigName("appSetting")
	// 设置配置文件类型
	viper.SetConfigType("env") //or "json" or "yaml" or "toml"
	// 读取环境变量
	viper.AutomaticEnv()
	// 读取配置文件
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	// 将配置文件中的配置信息保存到config变量中
	err = viper.Unmarshal(&config)
	return
}

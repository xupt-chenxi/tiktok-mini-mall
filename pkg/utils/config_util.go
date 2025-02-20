package utils

import (
	"github.com/spf13/viper"
	"log"
)

type DatabaseConfig struct {
	DSN string `mapstructure:"dsn"`
}

type ServiceConfig struct {
	IP       string         `mapstructure:"ip"`
	Port     string         `mapstructure:"port"`
	Database DatabaseConfig `mapstructure:"database"`
}

type RedisConfig struct {
	IP       string `mapstructure:"ip"`
	Port     string `mapstructure:"port"`
	DB       string `mapstructure:"db"`
	Password string `mapstructure:"password"`
}

type RocketConfig struct {
	Endpoint  string `mapstructure:"endpoint"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	TopicProd string `mapstructure:"topic_prod"`
}

type AllConfig struct {
	App struct {
		Name string `mapstructure:"name"`
	} `mapstructure:"app"`
	User     ServiceConfig `mapstructure:"user"`
	Product  ServiceConfig `mapstructure:"product"`
	Cart     ServiceConfig `mapstructure:"cart"`
	Shop     ServiceConfig `mapstructure:"shop"`
	Redis    RedisConfig   `mapstructure:"redis"`
	RocketMQ RocketConfig  `mapstructure:"rocketmq"`
}

var Config AllConfig

func init() {
	InitConfig("configs/config.yaml")
}

func InitConfig(configFile string) {
	viper.SetConfigFile(configFile)
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	// 将配置文件的内容映射到结构体
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("配置文件转换失败: %v", err)
	}
}

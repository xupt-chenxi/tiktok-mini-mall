package utils

import (
	"github.com/spf13/viper"
	"log"
)

func InitViper(configFile string) {
	viper.SetConfigFile(configFile)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	log.Println("读取配置文件成功: ", viper.ConfigFileUsed())
}

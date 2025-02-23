package utils

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

func NewNamingClient() (naming_client.INamingClient, error) {
	clientConfig := constant.ClientConfig{
		NamespaceId:         Config.Nacos.NamespaceID,
		TimeoutMs:           uint64(Config.Nacos.TimeoutMs),
		NotLoadCacheAtStart: true,
		LogDir:              Config.Nacos.LogDir,
		CacheDir:            Config.Nacos.CacheDir,
	}

	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: Config.Nacos.IP,
			Port:   uint64(Config.Nacos.Port),
		},
	}

	namingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		return nil, err
	}
	return namingClient, nil
}

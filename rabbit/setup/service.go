package setup

import (
	"fmt"
	"strconv"

	"alsritter.icu/rabbit"
	"alsritter.icu/rabbit/internal/config"
	"alsritter.icu/rabbit/internal/config/etcdconfig"
	"alsritter.icu/rabbit/internal/util"
)

const (
	defaultServicePort = "9000"
	defaultHttpPort    = "10000"
)

type ServiceInfo struct {
	ServicePort string
	HttpPort    string
	SSL         bool
}

// NewServiceInfo init the service info.
func NewServiceInfo(app *rabbit.Application) (*ServiceInfo, error) {
	currentServicePort := strconv.Itoa(util.RandInt(50000, 60000))
	currentHttpPort := defaultHttpPort
	currentIsSSL := false

	etcdServerUrls := config.GetEtcdServerURL()
	if etcdServerUrls == "" {
		return nil, fmt.Errorf("Can't not found env '%s'", config.ENV_ETCD_SERVER_URL)
	}

	// get etcd config server
	serviceConfig := etcdconfig.NewServiceConfig(etcdServerUrls, app.Name)
	serviceConfig.InitCache()

	serviceConfigs, err := serviceConfig.GetConfigs()
	if err != nil {
		return nil, fmt.Errorf("serviceConfig.GetConfigs err: %v", err)
	}

	currentKey := serviceConfig.GetKeyName(app.Name)

	// get server port.
	for key, value := range serviceConfigs {
		// If this service already exists, so the direct reuse of the original configuration.
		if currentKey == key {
			currentServicePort = value.ServicePort
			break
		}

		if value.ServicePort == currentServicePort {
			return nil, fmt.Errorf("The service port is duplicated, please try again")
		}
	}

	currentHttpPort = currentServicePort
	// Maybe the rabbit version has been updated, so here need to update the ETCD
	err = serviceConfig.WriteConfig(etcdconfig.Config{
		ServiceVersion: rabbit.RABBIT_VERSION,
		ServicePort:    currentServicePort,
		HttpPort:       currentHttpPort,
		IsSsl:          currentIsSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("serviceConfig.WriteConfig err: %v", err)
	}

	return &ServiceInfo{
		ServicePort: currentServicePort,
		HttpPort:    currentHttpPort,
		SSL:         currentIsSSL,
	}, nil
}

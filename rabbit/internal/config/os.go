package config

import "os"

const (
	ENV_ETCD_SERVER_URL = "ETCD_SERVER_URL"
	ENV_ETCD_USERNAME   = "ETCD_USERNAME"
	ENV_ETCD_PASSWORD   = "ETCD_PASSWORD"
)

// GetEtcdServerURL get the etcd server URL from env
func GetEtcdServerURL() string {
	return os.Getenv(ENV_ETCD_SERVER_URL)
}

// GetEtcdUsername gets etcd v3 auth username config from env.
func GetEtcdUsername() string {
	return os.Getenv(ENV_ETCD_USERNAME)
}

// GetEtcdPassword gets etcd v3 auth password config from env.
func GetEtcdPassword() string {
	return os.Getenv(ENV_ETCD_PASSWORD)
}

package etcdconfig

type Config struct {
	ServiceVersion string `json:"service_version"`
	ServicePort    string `json:"service_port"`
	HttpPort       string `json:"http_port"`
	IsSsl          bool   `json:"is_ssl"`
}

func NewServiceConfig() {

}

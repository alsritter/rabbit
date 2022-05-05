package rabbit

const (
	APP_OPTIONS_CONFIGER AppOptionType = iota + 1
	APP_OPTIONS_REGISTRAR
	APP_OPTIONS_DISCOVERY
)

const (
	APP_TYPE_GRPC AppType = iota + 1
	APP_TYPE_QUEUE
	APP_TYPE_HTTP
	APP_TYPE_GIN
)

const (
	// After executing "NewRunner"
	POS_NEW_RUNNER CallbackPos = iota + 1
)

type Application struct {
	ApplicationInfo

	RegisterCallback map[CallbackPos]func() error
	AppOptions       map[AppOptionType]OptionIface
}

type ApplicationInfo struct {
	Name string
	Port int64
	Type AppType
}

type ServiceInstance struct {
	ServiceName    string
	ServiceIP      string
	ServicePort    int
	HttpPort       int
	ServiceVersion string
	SSL            bool
}

// OptionDisable used to close app option
type OptionDisable struct {
	OptionIface
}

// ConfigerIface is Service Config Management Interface.
type ConfigerIface interface {
	OptionIface
	Init() error
	IsExistAppConfig(*Application) (bool, error)
	LoadAppConfig(*Application) error
	WatchUpdateConfig()
}

// RegistrarIface is Service Register Standard Interface
type RegistrarIface interface {
	OptionIface
	Init() error
	Register(service *ApplicationInfo) (*ServiceInstance, error)
	Deregister(service *ApplicationInfo) error
}

// DiscoveryIface is Service Discovery Standard Interface
type DiscoveryIface interface {
	OptionIface
	Init() error
	GetService(serviceName string) ([]*ServiceInstance, error)
	LoadBalance(instances []*ServiceInstance) (*ServiceInstance, error)
}

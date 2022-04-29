package rabbit

const (
	APP_OPTIONS_CONFIGER AppOptionType = iota + 1
	APP_OPTIONS_REGISTRAR
	APP_OPTIONS_DISCOVERY
)

const (
	// After executing "NewRunner"
	POS_NEW_RUNNER CallbackPos = iota + 1
)

type Application struct {
	Name             string
	Port             int64
	RegisterCallback map[CallbackPos]func() error
	AppOptions       map[AppOptionType]OptionIface
}

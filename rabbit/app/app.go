package app

import (
	"errors"

	"alsritter.icu/rabbit"
)

type Runner struct {
	App *rabbit.Application
}

func NewRunner(app *rabbit.Application) (*Runner, error) {
	if app.Name == "" {
		return nil, errors.New("Application name can't not be empty")
	}

	if app.Type <= 0 {
		return nil, errors.New("Application type can't not be empty")
	}

	if err := validAppOption(app.AppOptions); err != nil {
		return nil, err
	}

}

// ServiceRegister server register to etcd.
func (r *Runner) ServiceRegister(application *rabbit.Application) {

}

func validAppOption(appOptions map[rabbit.AppOptionType]rabbit.OptionIface) error {
	for optionType, option := range appOptions {
		if option == nil {
			continue
		}

		if optionType == rabbit.APP_OPTIONS_CONFIGER {
			switch option.(type) {
			case rabbit.OptionDisable:
			case rabbit.ConfigerIface:
			default:
				return errors.New("the type of Application Configer must be rabbit.ConfigerIface")
			}
		}

		if optionType == rabbit.APP_OPTIONS_REGISTRAR {
			switch option.(type) {
			case rabbit.RegistrarIface:
			case rabbit.OptionDisable:
			default:
				return errors.New("the type of Application Registrar must be rabbit.RegistrarIface")
			}
		}

		if optionType == rabbit.APP_OPTIONS_DISCOVERY {
			switch option.(type) {
			case rabbit.DiscoveryIface:
			case rabbit.OptionDisable:
			default:
				return errors.New("the type of Application Discovery Type must be rabbit.DiscoveryIface")
			}
		}
	}

	return nil
}

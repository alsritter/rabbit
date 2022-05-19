package app

import (
	"errors"
	"fmt"
	"strconv"

	"alsritter.icu/rabbit"
	"alsritter.icu/rabbit/internal/util"
	"alsritter.icu/rabbit/setup"
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

	return &Runner{App: app}, nil
}

// ServiceRegister server register to etcd.
// To avoid collisions, the registered service ports are randomly generated.
func (r *Runner) ServiceRegister() (*setup.ServiceInfo, error) {
	registrar := r.App.AppOptions[rabbit.APP_OPTIONS_REGISTRAR]

	// if register is empty, default write to etcd register center.
	if registrar == nil {
		serviceInfo, err := setup.NewServiceInfo(r.App)
		if err != nil {
			return nil, err
		}
		return serviceInfo, nil
	}

	switch registrar.(type) {
	// The service registry empty interface.
	case rabbit.OptionDisable:
		port := r.App.Port
		// Register
		if port <= 0 {
			port = int64(util.RandInt(50000, 60000))
		}
		return &setup.ServiceInfo{
			ServicePort: strconv.Itoa(int(port)),
			HttpPort:    strconv.Itoa(int(port)),
		}, nil

	// Service registry interface with callback function.
	case rabbit.RegistrarIface:
		registrar := registrar.(rabbit.RegistrarIface)
		err := registrar.Init()
		if err != nil {
			return nil, err
		}

		// Register
		if r.App.Port <= 0 {
			r.App.Port = int64(util.RandInt(50000, 60000))
		}

		serviceInstance, err := registrar.Register(&rabbit.ApplicationInfo{
			Name: r.App.Name,
			Port: r.App.Port,
		})

		if err != nil {
			return nil, errors.New(fmt.Sprintf("register service failed, reason: %v", err))
		}

		// Get Instance Info
		if serviceInstance == nil {
			return nil, errors.New("register service failed, please retry !")
		}

		return &setup.ServiceInfo{
			ServicePort: strconv.Itoa(serviceInstance.ServicePort),
			HttpPort:    strconv.Itoa(serviceInstance.HttpPort),
			SSL:         serviceInstance.SSL,
		}, nil

	default:
		return nil, errors.New("the type of Application Registrar must be rabbit.RegistrarIface or rabbit.OptionDisable")
	}
}

// check the callback implementation a rabbit.OptionIface.
func validAppOption(appOptions map[rabbit.AppOptionType]rabbit.OptionIface) error {
	for optionType, option := range appOptions {
		if option == nil {
			continue
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

package service

import (
	"github.com/MichalMichalak/arx2/log"
	"github.com/pkg/errors"
)

// Set logger before anything else
type ServiceBuilder struct {
	name        string
	logger      log.Logger
	providers   map[string]Provider
	error       error
	configPaths []string
}

func (sb *ServiceBuilder) validateLogger() {
	if sb.logger != nil {
		return
	}
	sb.logger = log.NewServiceLogger(log.SeverityInfo)
	sb.logger.Warn("logger not defined, setting default one")
}

func (sb *ServiceBuilder) Name(name string) *ServiceBuilder {
	// Pass through if error's already there
	if sb.error != nil {
		return sb
	}

	// Validate logger
	sb.validateLogger()

	// Can be set only once
	if sb.name != "" {
		sb.logger.Warnf("trying to set name `%s` but `%s` it's already set", name, sb.name)
		return sb
	}

	// Can't have empty name
	if name == "" {
		sb.error = errors.New("missing service name")
		return sb
	}

	// Set name
	sb.name = name
	sb.logger.Debugf("setting name `%s`", name)
	return sb
}

func (sb *ServiceBuilder) Logger(logger log.Logger) *ServiceBuilder {
	sb.logger = logger
	logger.Debug("setting logger")
	return sb
}

func (sb *ServiceBuilder) Provider(provider Provider) *ServiceBuilder {
	// Pass through if error's already there
	if sb.error != nil {
		return sb
	}

	// Validate logger
	sb.validateLogger()

	// Provider can't be empty
	if provider == nil {
		sb.error = errors.New("trying to register empty provider")
		return sb
	}

	// Provider's name can't be empty
	providerName := provider.Name()
	if providerName == "" {
		sb.error = errors.New("provider's name can't be empty")
		return sb
	}

	// Can't register providers with the same name
	if _, exist := sb.providers[providerName]; exist {
		sb.error = errors.Errorf("provider with name `%s` already registered", providerName)
		return sb
	}

	// Register provider
	sb.logger.Infof("registering provider `%s` of `%T` type", provider.Name(), provider)
	if sb.providers == nil {
		sb.providers = map[string]Provider{providerName: provider}
	} else {
		sb.providers[providerName] = provider
	}
	return sb
}

func (sb *ServiceBuilder) ConfigPaths(paths []string) *ServiceBuilder {
	// Pass through if error's already there
	if sb.error != nil {
		return sb
	}

	// Validate logger
	sb.validateLogger()

	sb.configPaths = paths
	return sb
}

func (sb *ServiceBuilder) Build() (*Service, error) {
	// Validate logger
	sb.validateLogger()

	// Cover error
	if sb.error != nil {
		return nil, errors.Wrap(sb.error, "initialization error")
	}

	// Name
	if sb.name == "" {
		return nil, errors.Wrap(sb.error, "empty name")
	}

	// Providers
	if sb.providers == nil {
		sb.providers = map[string]Provider{}
		sb.logger.Warn("no providers registered")
	}

	// Create and return the service
	s, err := newService(sb.name, sb.logger, sb.configPaths, sb.providers)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create service `%s`", sb.name)
	}
	return &s, nil
}

package service

import (
	"github.com/MichalMichalak/arx2/log"
)

type Context interface {
	ServiceName() string
	Logger() log.Logger
	Providers() map[string]Provider
}

type ServiceContext struct {
	serviceName string
	logger      log.Logger
	providers   map[string]Provider
}

func NewServiceContext(serviceName string, logger log.Logger, providers map[string]Provider) ServiceContext {
	return ServiceContext{serviceName: serviceName, logger: logger, providers: providers}
}

func (sc ServiceContext) Logger() log.Logger {
	return sc.logger
}

func (sc ServiceContext) ServiceName() string {
	return sc.serviceName
}

func (sc ServiceContext) Providers() map[string]Provider {
	return sc.providers
}

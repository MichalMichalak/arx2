package svc

import (
	"sync"

	"github.com/MichalMichalak/arx2/cnf"
	"github.com/MichalMichalak/arx2/log"
	"github.com/pkg/errors"
)

type Provider interface {
	Name() string
	Run() error
	ShutdownRequest()
	Configure(ctx Context, resolver cnf.Resolver) error
}

type Service struct {
	status  Status
	context Context
}

func newService(name string, logger log.Logger, configPaths []string, providers map[string]Provider) (Service, error) {
	context := NewServiceContext(name, logger, providers)
	confResolver := cnf.NewResolver(logger, configPaths)
	err := configureProviders(context, confResolver, providers)
	if err != nil {
		return Service{}, errors.Wrapf(err, "failed to initialize service `%s`", name)
	}
	return Service{context: context, status: Created}, nil
}

func configureProviders(context Context, resolver cnf.Resolver, providers map[string]Provider) error {
	for name, provider := range providers {
		err := provider.Configure(context, resolver)
		if err != nil {
			return errors.Wrapf(err, "failed to configure provider `%s`", name)
		}
	}
	return nil
}

func (s Service) Run() error {
	defer s.shutdown()
	if !s.status.IsCreated() {
		return errors.Errorf("invalid state, expected %s, got %s", Created, s.status)
	}
	s.status = Running
	wg := sync.WaitGroup{}
	providers := s.context.Providers()
	wg.Add(len(providers))
	for _, p := range providers {
		go runProvider(p, &wg)
	}
	wg.Wait()
	return nil
}

func (s Service) shutdown() {
	s.status = Closing
	providers := s.context.Providers()
	wg := sync.WaitGroup{}
	wg.Add(len(providers))
	for _, p := range providers {
		go shutdownProvider(p, &wg)
	}
	wg.Wait()
	s.status = Closed
}

func shutdownProvider(p Provider, wg *sync.WaitGroup) {
	p.ShutdownRequest()
	wg.Done()
}

func runProvider(p Provider, wg *sync.WaitGroup) {
	err := p.Run()
	if err != nil {
		err = errors.Wrapf(err, "[%s] execution interrupted", p.Name())
		panic(err)
		// TODO find better way to handle error
	}
	wg.Done()
}

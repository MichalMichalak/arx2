package arx2_test

import (
	"github.com/MichalMichalak/arx2/log"
	"github.com/MichalMichalak/arx2/provider"
	"github.com/MichalMichalak/arx2/service"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func TestArx2(t *testing.T) {
	numChan := make(chan int)

	p1 := provider.NewMyDoer()
	p2 := provider.NewMyConsumer(numChan)
	l := log.NewServiceLogger(log.SeverityDebug)

	sb := &service.ServiceBuilder{}
	s, err := sb.Logger(l).Name("s1").Provider(p1).Provider(p2).ConfigPaths([]string{"../_temp/config.yaml"}).Build()
	require.NoError(t, err)
	require.NotNil(t, s)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go f1(s, t, wg)
	go f2(numChan, wg)
	wg.Wait()
}

func f2(numChan chan int, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	for i := 1; i < 10; i++ {
		numChan <- i
	}
	numChan <- 0
	wg.Done()
}

func f1(s *service.Service, t *testing.T, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	err := s.Run()
	require.NoError(t, err)
	wg.Done()
}

func TestConsumer(t *testing.T) {
	numChan := make(chan int)
	p := provider.NewMyConsumer(numChan)
	wg := sync.WaitGroup{}
	wg.Add(2)

	ctx := service.NewServiceContext("s1", log.NewServiceLogger(log.SeverityDebug), map[string]service.Provider{})

	go func() {
		err := p.Run(ctx)
		require.NoError(t, err)
		wg.Done()
	}()

	go func() {
		numChan <- 1
		numChan <- 2
		numChan <- 3
		numChan <- 0
		wg.Done()
	}()

	wg.Wait()
}

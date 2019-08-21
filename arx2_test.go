package arx2_test

import (
	"sync"
	"testing"
	"time"

	"github.com/MichalMichalak/arx2/cnf"
	"github.com/MichalMichalak/arx2/log"
	"github.com/MichalMichalak/arx2/provider"
	"github.com/MichalMichalak/arx2/svc"
	"github.com/stretchr/testify/require"
)

func TestArx2(t *testing.T) {
	numChan := make(chan int)

	p1 := provider.NewMyDoer()
	p2 := provider.NewMyConsumer(numChan)
	l := log.NewServiceLogger(log.SeverityDebug)

	sb := &svc.Builder{}
	s, err := sb.Logger(l).Name("s1").Provider(p1).Provider(p2).ConfigPaths([]string{"../_temp/config.yaml"}).Build()
	require.NoError(t, err)
	require.NotNil(t, s)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go f1(s, t, wg)
	go f2(numChan, wg)
	wg.Wait()
}

func f1(s *svc.Service, t *testing.T, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	err := s.Run()
	require.NoError(t, err)
	wg.Done()
}

func f2(numChan chan int, wg *sync.WaitGroup) {
	time.Sleep(1 * time.Second)
	for i := 1; i < 10; i++ {
		numChan <- i
	}
	numChan <- 0
	wg.Done()
}

func TestConsumer(t *testing.T) {
	numChan := make(chan int)
	p := provider.NewMyConsumer(numChan)
	wg := sync.WaitGroup{}
	wg.Add(2)

	logger := log.NewServiceLogger(log.SeverityDebug)
	ctx := svc.NewServiceContext("s1", logger, map[string]svc.Provider{})
	err := p.Configure(ctx, cnf.NewResolver(logger, nil))
	require.NoError(t, err)

	go func() {
		err := p.Run()
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

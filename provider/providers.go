package provider

import (
	"time"

	"github.com/MichalMichalak/arx2/cnf"
	"github.com/MichalMichalak/arx2/log"
	"github.com/MichalMichalak/arx2/svc"
)

// --------------------
// Provider for testing
// --------------------
type MyDoer struct {
	name   string
	logger log.Logger
}

func NewMyDoer() *MyDoer {
	return &MyDoer{name: "doer"}
}

func (mp *MyDoer) Name() string { return mp.name }

func (mp *MyDoer) Run() error {
	mp.logger.Log(log.SeverityDebug, mp.name+" - Run")
	defer mp.logger.Log(log.SeverityDebug, mp.name+" - Run complete")
	return nil
}

func (mp *MyDoer) DoSomething(severity log.Severity, msg string) {
	mp.logger.Error("error, just for fun")
}

func (mp *MyDoer) ShutdownRequest() {
	mp.logger.Log(log.SeverityInfo, mp.name+" - ShutdownRequest")
	defer mp.logger.Log(log.SeverityInfo, mp.name+" - ShutdownRequest complete")
}

func (mp *MyDoer) Configure(ctx svc.Context, resolver cnf.Resolver) error {
	mp.logger = ctx.Logger()
	return nil
}

// --------------------
// Provider for testing
// --------------------
type myConsConf struct {
	Topic string `conf:"some.fake.topic"`
}

type MyConsumer struct {
	name          string
	numberChannel chan int
	logger        log.Logger
	conf          myConsConf
}

func NewMyConsumer(numberChannel chan int) *MyConsumer {
	return &MyConsumer{numberChannel: numberChannel, name: "consumer"}
}

func (mc *MyConsumer) Name() string { return mc.name }

func (mc *MyConsumer) Run() error {
	mc.logger.Infof("Config %+v", mc.conf)
	mc.logger.Debug("run")
	for {
		mc.logger.Debug("run waiting")
		i := <-mc.numberChannel
		if i == 0 {
			close(mc.numberChannel)
			mc.logger.Warn("run got 0 stop")
			break
		}
		mc.logger.Info(i)
	}
	mc.logger.Info("run done")
	return nil
}

func (mc *MyConsumer) ShutdownRequest() {
	mc.logger.Warnf("%s shutting down...", mc.name)
	time.Sleep(5 * time.Second)
	mc.logger.Warnf("%s shutdown", mc.name)
}

func (mc *MyConsumer) Configure(ctx svc.Context, resolver cnf.Resolver) error {
	mc.logger = ctx.Logger()
	return cnf.Configure(&mc.conf, resolver)
}

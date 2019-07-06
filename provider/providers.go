package provider

import (
	"github.com/unravelin/core/arx2/conf"
	"github.com/unravelin/core/arx2/log"
	"github.com/unravelin/core/arx2/service"
	"time"
)

type MyDoer struct {
	name   string
	logger log.Logger
}

func NewMyDoer() *MyDoer {
	return &MyDoer{name: "doer"}
}

func (mp *MyDoer) Name() string { return mp.name }

func (mp *MyDoer) Run(ctx service.Context) error {
	mp.logger = ctx.Logger()
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

func (mp *MyDoer) Configure(resolver conf.Resolver) error {
	return nil
}

// #############

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

func (mc *MyConsumer) Run(ctx service.Context) error {
	mc.logger = ctx.Logger()
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

func (mc *MyConsumer) Configure(resolver conf.Resolver) error {
	err := conf.Configure(&mc.conf, resolver)
	if err != nil {
		return err
	}
	return nil
}

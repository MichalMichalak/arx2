package main

import (
	"github.com/MichalMichalak/arx2/idgen"
	"github.com/MichalMichalak/arx2/log"
	"github.com/MichalMichalak/arx2/provider/nsq"
	"github.com/MichalMichalak/arx2/svc"
)

func main() {
	logger := log.NewServiceLogger(log.SeverityDebug)
	cons := nsq.NewConsumer(idgen.UUIDGenerator{}, handler(logger))
	service, err := svc.NewBuilder().
		Logger(logger).
		Name("svc-1").
		ConfigPath("config/nsq-test.yaml").
		Provider(cons).
		Build()
	if err != nil {
		panic(err)
	}
	err = service.Run()
	if err != nil {
		panic(err)
	}
}

func handler(logger log.Logger) nsq.Handler {
	return func(payload []byte) error {
		logger.Infof("MESSAGE RECEIVED: `%s`\n", string(payload))
		return nil
	}
}

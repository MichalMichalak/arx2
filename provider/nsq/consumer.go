package nsq

import (
	"fmt"
	"sync"

	"github.com/MichalMichalak/arx2/cnf"
	"github.com/MichalMichalak/arx2/idgen"
	"github.com/MichalMichalak/arx2/log"
	"github.com/MichalMichalak/arx2/svc"
	"github.com/google/uuid"
	"github.com/nsqio/go-nsq"
)

type ConsConfig struct {
	Topic   string `conf:"nsq.consumer.topic"`
	NsqHost string `conf:"nsq.consumer.host"`
	id      string
	temp    string
}

type Handler func(payload []byte) error

type Consumer struct {
	idGen   idgen.IDGenerator
	handler Handler
	config  ConsConfig
	logger  log.Logger
	blocker sync.WaitGroup
}

//type Handler func()

func NewConsumer(idGen idgen.IDGenerator, handler Handler) *Consumer {
	consumer := Consumer{idGen: idGen, handler: handler}
	consumer.config.id = "nsq-consumer-" + idGen.NewId()
	return &consumer
}

func (c *Consumer) Name() string {
	return c.config.id
}

func (c *Consumer) Run() error {
	// Start listening and return
	nsqConf := nsq.NewConfig()
	nsqConsumer, err := nsq.NewConsumer(c.config.Topic, c.Name(), nsqConf)
	if err != nil {
		return err
	}
	var h nsq.HandlerFunc = func(message *nsq.Message) error {
		return c.handler(message.Body)
	}
	nsqConsumer.AddHandler(h)
	go c.start(nsqConsumer)
	c.blocker = sync.WaitGroup{}
	c.blocker.Add(1)
	c.blocker.Wait()
	return nil
}

func (c *Consumer) start(consumer *nsq.Consumer) {
	err := consumer.ConnectToNSQD(c.config.NsqHost)
	if err != nil {
		c.logger.Errorf("Provider `%s` failed to connect to NSQ server `%s`: %v", c.Name(), c.config.NsqHost, err)
		return
	}
}

func (c *Consumer) ShutdownRequest() {
	c.logger.Info("closing nsq consumer")
	c.blocker.Done()
}

func (c *Consumer) Configure(ctx svc.Context, resolver cnf.Resolver) error {
	// Setup anything you need from context
	c.logger = ctx.Logger()

	// Directly
	c.config.temp = uuid.New().String()

	// From resolver, but you have to cast
	t := resolver.Conf()["consumer.temp"]
	if tt, ok := t.(string); ok {
		c.config.temp = tt
	}

	// Through injection
	var err error
	err = cnf.Configure(&c.config, resolver)
	fmt.Printf("CONS CONF=%+v\n", c.config)
	return err
}

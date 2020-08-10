package collect

import (
	"errors"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"time"
)

type Consumer interface {
	Consume(e Event) error
	Shutdown() error
	//UUID() string
}

type Event interface {
}

type ConsumerFactory interface {
	CreateConsumer() (Consumer, error)
}

type ConsumerManager struct {
	factory       ConsumerFactory
	freeConsumers chan Consumer
	consumers     chan Consumer
	producer      chan struct{}

	Max, Min          int64
	Timeout, idleTime time.Duration // 等待consumer时间

	isShutdown bool
}

func (c *ConsumerManager) Shutdown() error {
	c.isShutdown = true
	for freeConsumer := range c.freeConsumers {
		if err := freeConsumer.Shutdown(); err != nil {
			inline.WithFields("consumer", freeConsumer).Errorln("shutDown fail")
		}
	}

	return nil
}

func (c *ConsumerManager) outOfRange() bool {
	return int64(len(c.consumers)) >= c.Max
}

func (c *ConsumerManager) shouldFastAdd() bool {
	return int64(len(c.consumers)) < c.Min
}

func (c *ConsumerManager) shouldClose() bool {
	return int64(len(c.freeConsumers)) > c.Min
}

// 这里策略可以有很多，比如尽可能的往下访问，又或者是直接在当前报错
func (c *ConsumerManager) Consume(e Event) error {
	select {
	case c.producer <- struct{}{}:
	case <-time.NewTimer(c.Timeout).C:
		return errors.New("wait consumer timeout")
	}

	defer func() {
		<-c.producer
	}()

	if c.isShutdown {
		return errors.New("manager has shutdown")
	}
	if c.shouldFastAdd() {
		return c.createAndConsumer(e)
	}

	timer := time.NewTimer(c.Timeout)
	select {
	case consumer := <-c.freeConsumers:
		inline.WithFields("e", e).Infoln("borrow consumer")
		return c.execute(consumer, e)
	case <-timer.C:
		if c.outOfRange() {
			return errors.New("out of range")
		} else {
			return c.createAndConsumer(e)
		}
	}
}

func (c *ConsumerManager) closeIdle() {
	ticker := time.NewTicker(c.idleTime)
	for range ticker.C {
		if c.shouldClose() {
			consumer := <-c.freeConsumers
			<-c.consumers
			if err := consumer.Shutdown(); err != nil {
				inline.WithFields("consumer", consumer).Infoln("shutdown idle err")
			} else {
				inline.WithFields("consumer", consumer).Infoln("destroy success")
			}
		}
	}
}

func (c *ConsumerManager) createAndConsumer(e Event) (err error) {
	defer func() {
		r, ok := recover().(error)
		if ok {
			err = r
		}
	}()
	consumer, err := c.factory.CreateConsumer()
	if err != nil {
		return inline.PrependErrorFmt(err, "create fail")
	}

	c.consumers <- consumer
	fmt.Println("create consumer")
	if err := c.execute(consumer, e); err != nil {
		return inline.PrependErrorFmt(err, "execute fail")
	}
	return nil
}

func (c *ConsumerManager) execute(consumer Consumer, e Event) (err error) {
	defer func() {
		r, ok := recover().(error)
		if ok {
			err = r
		}
	}()
	if err := consumer.Consume(e); err != nil {
		<-c.consumers
		if err := consumer.Shutdown(); err != nil {
			return inline.PrependErrorFmt(err, "shut down %+v", consumer)
		}
		return inline.PrependErrorFmt(err, "consume %+v", consumer)
	}
	fmt.Println("add to free")

	c.freeConsumers <- consumer
	return nil
}

type managerBuilder struct {
	consumer *ConsumerManager
}

func NewManagerBuilder() *managerBuilder {
	return &managerBuilder{consumer: &ConsumerManager{
		Timeout:  time.Second,
		idleTime: time.Second,
		Max:      20, Min: 10,
		freeConsumers: make(chan Consumer, 20),
		producer:      make(chan struct{}, 20),
		consumers:     make(chan Consumer, 20),
	}}
}

func (b *managerBuilder) Max(max int64) *managerBuilder {
	b.consumer.Max = max
	b.consumer.freeConsumers = make(chan Consumer, max)
	b.consumer.consumers = make(chan Consumer, max)
	return b
}

func (b *managerBuilder) Backend(cnt int64) *managerBuilder {
	b.consumer.producer = make(chan struct{}, cnt)
	return b
}

func (b *managerBuilder) Min(min int64) *managerBuilder {
	b.consumer.Min = min
	return b
}

func (b *managerBuilder) Timeout(timeout time.Duration) *managerBuilder {
	b.consumer.Timeout = timeout
	return b
}

func (b *managerBuilder) IdleTimeout(timeout time.Duration) *managerBuilder {
	b.consumer.idleTime = timeout
	return b
}

func (b *managerBuilder) Factory(factory ConsumerFactory) *managerBuilder {
	b.consumer.factory = factory
	return b
}

func (c *ConsumerManager) valid() error {
	if c.Min > c.Max {
		return fmt.Errorf("min %d > max %d", c.Min, c.Max)
	}
	if c.factory == nil {
		return errors.New("factory cannot be nil")
	}
	return nil
}

func (b *managerBuilder) Build() (*ConsumerManager, error) {
	if err := b.consumer.valid(); err != nil {
		return nil, err
	}
	go b.consumer.closeIdle()
	return b.consumer, nil
}

package collect

import (
	"errors"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"sync/atomic"
	"time"
)

type Consumer interface {
	Consume(e Event) error
	Shutdown() error
}

type Event interface {
}

type ConsumerFactory interface {
	CreateConsumer() (Consumer, error)
}

type ConsumerManager struct {
	factory       ConsumerFactory
	freeConsumers chan Consumer

	Max, Min, Count int64
	Timeout         time.Duration

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
	return atomic.LoadInt64(&c.Count) >= c.Max
}

func (c *ConsumerManager) shouldFastAdd() bool {
	return atomic.LoadInt64(&c.Count) < c.Min
}

func (c *ConsumerManager) incr() error {
	return c.add(func() error {
		return nil
	})
}

func (c *ConsumerManager) decr(f func() error) error {
	for old := c.Count; !atomic.CompareAndSwapInt64(&c.Count, old, old-1); old = c.Count {
		if old <= 0 {
			return errors.New("less zero")
		}
	}
	return f()
}

func (c *ConsumerManager) add(f func() error) error {
	for old := c.Count; !atomic.CompareAndSwapInt64(&c.Count, old, old+1); old = c.Count {
		if c.outOfRange() {
			return errors.New("out of range")
		}
	}

	return f()
}

func (c *ConsumerManager) returnConsumer(consumer Consumer) {
	if !c.isShutdown {
		c.freeConsumers <- consumer
	}
}

// 这里策略可以有很多，比如尽可能的往下访问，又或者是直接在当前报错
func (c *ConsumerManager) Consume(e Event) error {
	if c.shouldFastAdd() {
		return c.createAndConsumer(e)
	}

	timer := time.NewTimer(c.Timeout)
	select {
	case consumer := <-c.freeConsumers:
		return c.execute(consumer, e)
	case <-timer.C:
		if c.outOfRange() {
			return errors.New("out of range")
		} else {
			return c.createAndConsumer(e)
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
	if err := consumer.Consume(e); err != nil {
		if err := consumer.Shutdown(); err != nil {
			return inline.PrependErrorFmt(err, "shutdown err")
		}
		return inline.PrependErrorFmt(err, "consumer err event %+v", e)
	}
	return c.add(func() error {
		c.returnConsumer(consumer)
		return nil
	})

}

func (c *ConsumerManager) execute(consumer Consumer, e Event) (err error) {
	defer func() {
		r, ok := recover().(error)
		if ok {
			err = r
		}
	}()
	if err := consumer.Consume(e); err != nil {
		if err := c.decr(func() error {
			return consumer.Shutdown()
		}); err != nil {
			return inline.PrependErrorFmt(err, "decr")
		}

		return inline.PrependErrorFmt(err, "consume %+v", e)
	}

	c.returnConsumer(consumer)
	return nil
}

type managerBuilder struct {
	consumer *ConsumerManager
}

func ManagerBuilder() *managerBuilder {
	return &managerBuilder{consumer: &ConsumerManager{Timeout: time.Second, Max: 20, Min: 10, freeConsumers: make(chan Consumer, 20)}}
}

func (b *managerBuilder) Max(max int64) *managerBuilder {
	b.consumer.Max = max
	b.consumer.freeConsumers = make(chan Consumer, max)
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
	return b.consumer, nil
}

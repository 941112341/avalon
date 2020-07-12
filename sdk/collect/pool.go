package collect

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"sync"
	"time"
)

type IDGetter interface {
	ID() string
}

type ConsumerFactory func() (Consumer, error)

type Pool interface {
	IDGetter
	Producer(event interface{}) error
	IsShutdown() bool
	Shutdown() error
}

// key=serviceName value=chan
type pool struct {
	lock sync.Mutex
	id   string

	idleTime       time.Duration
	minConsumerNum int
	maxConsumerNum int // for each key
	consumers      []Consumer
	Timeout        time.Duration

	ch chan interface{}

	factory ConsumerFactory

	isShutdown bool
}

func (p *pool) ID() string {
	return p.id
}

func (p *pool) IsShutdown() bool {
	return p.isShutdown
}

func (p *pool) addConsumer() error {
	if p.maxConsumerNum > len(p.consumers) {
		consumer, err := p.factory()
		if err != nil {
			return errors.Wrap(err, "create consumer fail")
		}
		p.lock.Lock()
		defer p.lock.Unlock()
		consumer.Consumer(p.ch)
		p.consumers = append(p.consumers, consumer)
	}
	return nil
}

func (p *pool) Producer(event interface{}) error {
	if err := p.addConsumer(); err != nil {
		return errors.Wrap(err, "addConsumer")
	}
	select {
	case p.ch <- event:
		return nil
	case <-time.NewTimer(p.Timeout).C:
		return errors.New("consumer timeout")
	}
}

func (p *pool) Shutdown() (err error) {
	p.isShutdown = true
	close(p.ch)
	for _, consumer := range p.consumers {
		err = consumer.Close()
		if err != nil {
			return errors.Wrap(err, "consumer ")
		}
	}
	p.consumers = []Consumer{}
	return nil
}

func (p *pool) check() {
	ticker := time.NewTicker(p.idleTime)
	for !p.IsShutdown() {
		select {
		case <-ticker.C:
			p.removeIdleConsumer()
		}
	}
}

func (p *pool) removeIdleConsumer() {
	p.lock.Lock()
	defer p.lock.Unlock()
	consumers := make([]Consumer, 0)
	for _, c := range p.consumers {
		if time.Now().Sub(c.LastUseTime()) < p.idleTime {
			consumers = append(consumers, c)
		}
	}
	p.consumers = consumers
}

func NewPool(idleTime, timeout time.Duration, minNum, maxNum, backUp int, factory ConsumerFactory) Pool {
	return &pool{
		lock:           sync.Mutex{},
		id:             inline.RandString(32),
		idleTime:       idleTime,
		minConsumerNum: minNum,
		maxConsumerNum: maxNum,
		consumers:      []Consumer{},
		Timeout:        timeout,
		ch:             make(chan interface{}, backUp),
		factory:        factory,
		isShutdown:     false,
	}
}

type Consumer interface {
	Closable
	IDGetter
	Start()
	Consumer(ch chan interface{})
	LastUseTime() time.Time
}

type Closable interface {
	Close() error
}

type element struct {
	id          string
	lastUseTime time.Time
	Client      Closable
	ch          chan interface{}
	IsClosed    bool

	consumer func(event interface{}) error // err should be close this closable
}

func (e *element) Start() {
	for !e.IsClosed {
		select {
		case event := <-e.ch:
			err := e.consumer(event)
			if err != nil {
				e.Close()
			}
		}
	}
}

func (e *element) ID() string {
	return e.id
}

func (e *element) LastUseTime() time.Time {
	return e.lastUseTime
}

func (e *element) Consumer(ch chan interface{}) {
	e.ch = ch
}

func (e *element) Close() error {
	e.IsClosed = true
	return e.Client.Close()
}

func NewConsumer(closable Closable, consumer func(event interface{}) error) Consumer {
	elem := &element{
		id:          inline.RandString(32),
		lastUseTime: time.Now(),
		Client:      closable,
		consumer:    consumer,
	}
	elem.Start()
	return elem
}

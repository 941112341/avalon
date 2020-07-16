package collect

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"sync/atomic"
	"time"
)

var (
	ErrShutdown       = errors.New("pool: has shutdown")
	ErrSizeOutOfLimit = errors.New("pool: out of size")
)

const (
	NORMAL   int32 = 0
	SHUTDOWN int32 = 1
)

type Introspect interface {
	String() string
}

type IDGetter interface {
	ID() string
}

type idGetter struct {
	id string
}

func (i *idGetter) ID() string {
	return i.id
}

func NewIDGetter() IDGetter {
	return &idGetter{id: inline.RandString(32)}
}

type ConsumerFactory interface {
	Create() (Consumer, error)
}

type Pool interface {
	IDGetter
	Introspect
	IsShutdown() bool
	Shutdown() error
	GetConsumer() (Consumer, error)
	GetConsumerBlock(timeout time.Duration) (Consumer, error)
	PutConsumer(consumer Consumer) error
	DelConsumer(consumer Consumer) error
	Size() int
}

// key=serviceName value=chan
type pool struct {
	IDGetter

	idleTime       time.Duration
	minConsumerNum int
	maxConsumerNum int // for each key

	size       int32
	consumers  chan *element
	factory    ConsumerFactory
	isShutdown int32 // 0 . 1
}

func (p *pool) Size() int {
	return int(atomic.LoadInt32(&p.size))
}

func (p *pool) String() string {

	return inline.ToJsonString(map[string]interface{}{
		"min":  p.minConsumerNum,
		"max":  p.maxConsumerNum,
		"idle": p.idleTime.String(),
		"size": p.Size(),
	})
}

func (p *pool) incr(delta int32) {
	atomic.AddInt32(&p.size, delta)
}

func (p *pool) GetConsumer() (Consumer, error) {

	return p.GetConsumerBlock(10 * time.Millisecond)
}

func (p *pool) MoreThanLimit() bool {
	return p.Size() > p.maxConsumerNum
}

func (p *pool) LessThanLimit() bool {
	return p.Size() < p.minConsumerNum
}

func (p *pool) GetConsumerBlock(timeout time.Duration) (Consumer, error) {
	if p.IsShutdown() {
		return nil, ErrShutdown
	}
	if len(p.consumers) > 0 {
		select {
		case consumer := <-p.consumers:
			return consumer, nil
		case <-time.NewTimer(10 * time.Millisecond).C:
		}
	}
	if p.LessThanLimit() {
		return p.createNewElement()
	}
	select {
	case element := <-p.consumers:
		inline.Debugln("borrow element", inline.NewPair("id", element.ID()))
		return element, nil
	case <-time.NewTimer(timeout).C:
		return p.createNewElement()
	}
}

func (p *pool) PutConsumer(consumer Consumer) error {
	if p.IsShutdown() {
		return ErrShutdown
	}
	element, ok := consumer.(*element)
	if !ok {
		element = p.newElement(consumer)
	}
	select {
	case <-time.NewTimer(100 * time.Millisecond).C:
		p.incr(-1)
		inline.Warningln("put consumer fail", inline.NewPair("id", element.ID()), inline.NewPair("len", len(p.consumers)))
		return nil
	case p.consumers <- element:
		if !ok {
			p.incr(1)
		}
		inline.Debugln("return consumer success", inline.NewPair("id", element.ID()))
	}
	return nil
}

func (p *pool) DelConsumer(consumer Consumer) error {
	if p.IsShutdown() {
		return nil
	}
	p.incr(-1)
	return nil
}

// incr and lock
func (p *pool) withLock(f func() error) error {

	for size := p.Size(); size < p.maxConsumerNum; size = p.Size() {
		if atomic.CompareAndSwapInt32(&p.size, int32(size), int32(size+1)) {
			return f()
		}
	}
	return ErrSizeOutOfLimit
}

func (p *pool) createNewElement() (e *element, err error) {
	err = p.withLock(func() error {

		var consumer Consumer
		consumer, err = p.factory.Create()
		if err != nil {
			p.incr(-1)
			return errors.Wrap(err, "consumer factory create err")
		}
		e = p.newElement(consumer)
		inline.Debugln("pool create new element", inline.NewPair("id", e.ID()))
		return nil
	})

	return
}

func (p *pool) newElement(consumer Consumer) *element {
	return &element{
		IDGetter:    NewIDGetter(),
		lastUseTime: time.Now(),
		p:           p,
		c:           consumer,
	}
}

func (p *pool) IsShutdown() bool {
	return atomic.LoadInt32(&p.isShutdown) == SHUTDOWN
}

func (p *pool) Shutdown() (err error) {
	atomic.StoreInt32(&p.isShutdown, SHUTDOWN)

	for ok := true; ok && err == nil; {
		var elem *element
		elem, ok = <-p.consumers
		if ok {
			err = elem.ShutDown()
		}
	}
	if err != nil {
		close(p.consumers)
	}
	return
}

func (p *pool) checkLoop() {
	ticker := time.NewTicker(p.idleTime)

loop:
	for {
		select {
		case <-ticker.C:
			l := len(p.consumers)
			for i := 0; i < l; i++ {
				element, ok := <-p.consumers
				if !ok {
					break loop
				}
				if element.isIdle(p.idleTime) && !p.LessThanLimit() {
					err := element.ShutDown()
					if err != nil {
						inline.Errorln("shutdown err", inline.NewPair("err", err))
					}
				} else {
					p.consumers <- element
				}
			}
		}
	}
}

func NewPool(idleTime time.Duration, minNum, maxNum int, factory ConsumerFactory) Pool {
	pool := &pool{
		IDGetter:       NewIDGetter(),
		idleTime:       idleTime,
		minConsumerNum: minNum,
		maxConsumerNum: maxNum,
		consumers:      make(chan *element, maxNum),
		factory:        factory,
	}
	go pool.checkLoop()
	return pool
}

type Consumer interface {
	Closable
	Do(ctx context.Context, args ...interface{}) error
}

type Closable interface {
	Close() error
}

type element struct {
	IDGetter
	lastUseTime time.Time

	p Pool

	c        Consumer
	errTimes int
}

func (e *element) isIdle(idleTimeout time.Duration) bool {
	return time.Now().Sub(e.lastUseTime) > idleTimeout
}

func (e *element) Do(ctx context.Context, args ...interface{}) error {
	e.lastUseTime = time.Now()
	err := e.c.Do(ctx, args...)
	if err != nil {
		e.errTimes++
		if e.errTimes > 3 {
			e.ShutDown()
		}
		return errors.Wrap(err, "do err")
	}
	return nil
}

// fact return to p
func (e *element) Close() error {
	if e.errTimes <= 3 {
		return e.p.PutConsumer(e)
	}
	return nil
}

func (e *element) ShutDown() error {
	inline.Infoln("element shutdown", inline.NewPair("id", e.ID()))

	err := e.p.DelConsumer(e)
	if err != nil {
		return err
	}
	return e.c.Close()
}

package collect

import (
	"errors"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type CF struct {
}

func (C *CF) CreateConsumer() (Consumer, error) {
	return &TestConsumer{}, nil
}

type TestConsumer struct {
}

func (t *TestConsumer) Consume(e Event) error {
	r := rand.Intn(100)
	time.Sleep(time.Millisecond * 100 * time.Duration(r))
	if r > 50 {
		return errors.New("err")
	}
	return nil
}

func (t *TestConsumer) Shutdown() error {
	return nil
}

func TestManager(t *testing.T) {
	manager, err := NewManagerBuilder().Factory(&CF{}).Build()
	if err != nil {
		panic(err)
	}

	for i := 0; i < 30; i++ {
		go func() {
			err := manager.Consume(nil)
			fmt.Println(err)
		}()
	}

	time.Sleep(5 * time.Second)

	for i := 0; i < 30; i++ {
		go func() {
			err := manager.Consume(nil)
			fmt.Println(err)
		}()
	}

	time.Sleep(5 * time.Second)
}

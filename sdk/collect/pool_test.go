package collect

import (
	"context"
	"fmt"
	"github.com/941112341/avalon/sdk/inline"
	"math/rand"
	"testing"
	"time"
)

type TestConsumer struct {
}

func (t *TestConsumer) Close() error {
	fmt.Println("close")
	return nil
}

func (t *TestConsumer) Do(ctx context.Context, args ...interface{}) error {
	time.Sleep(time.Duration(rand.Int63n(3)) * time.Second)
	return nil
}

func TestPool(t *testing.T) {
	pool := NewPool(5*time.Second, 10, 20, func() (Consumer, error) {
		time.Sleep(1 * time.Second)
		return &TestConsumer{}, nil
	})
	ctx := context.Background()
	for i := 0; i < 30; i++ {
		go func() {
			consumer, err := pool.GetConsumerBlock(3 * time.Second)
			if err != nil {
				inline.Errorln("get consumer err", inline.NewPair("err", err))
				return
			}
			consumer.Do(ctx)
			consumer.Close()
		}()
	}
	time.Sleep(20 * time.Second)
	fmt.Println(pool.String())

}

func TestChan(t *testing.T) {
	ch := make(chan struct{}, 10)
	s, ok := <-ch
	fmt.Println(s, ok)
}

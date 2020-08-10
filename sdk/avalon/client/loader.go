package client

import (
	"errors"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/inline"
	"time"
)

type LoadBalancer interface {
	Balancer
	collect.Consumer
}

type RetryLoadBalancer struct {
	Balancer     Balancer
	hostConsumer map[string]*collect.ConsumerManager

	statistic both.RecordStatistic
	threshold both.Threshold
	timeout   time.Duration // 客户端超时时间
}

func (r *RetryLoadBalancer) Choices() []string {
	return r.Balancer.Choices()
}

// 策略todo抽象
func (r *RetryLoadBalancer) allFailIn(hostPort string) bool {
	now := time.Now()
	duration := r.threshold.Duration()
	before := now.Add(-duration)
	statistic := r.statistic.StatisticResultBetween(hostPort, before, now)

	return r.threshold.Threshold(statistic)
}

func (r *RetryLoadBalancer) Choice() (string, error) {
	hostList := r.Balancer.Choices()

	normalHosts := make([]string, 0)
	for _, host := range hostList {
		if r.allFailIn(host) {
			continue
		}
		normalHosts = append(normalHosts, host)
	}
	if len(normalHosts) == 0 {
		return "", errors.New("no normal hostports")
	}
	return inline.RandomList(normalHosts).(string), nil
}

func (r *RetryLoadBalancer) Consume(e collect.Event) error {
	return inline.Retry(func() error {
		hostport, err := r.Choice()
		if err != nil {
			return inline.PrependErrorFmt(err, "no hostport %s", inline.ToJsonString(e))
		}

		manager, ok := r.hostConsumer[hostport]
		if !ok {
			factory := NewDefaultFactory(hostport, r.timeout)
			manager, err = NewPool(100*time.Millisecond, 10, 20, factory)
			if err != nil {
				return inline.PrependErrorFmt(err, "new pool fail")
			}
			r.hostConsumer[hostport] = manager
		}

		if err := manager.Consume(e); err != nil {
			r.statistic.AddErrLog(hostport, e, err)
			return inline.PrependErrorFmt(err, "consume err")
		}
		r.statistic.AddSuc(hostport)
		return nil
	}, 3, time.Millisecond*10)
}

func (r *RetryLoadBalancer) Shutdown() error {
	for _, c := range r.hostConsumer {
		if err := c.Shutdown(); err != nil {
			return inline.PrependErrorFmt(err, "shutdown")
		}
	}
	return nil
}

func NewLoadBalancer(balancer Balancer, timeout time.Duration) LoadBalancer {
	return &RetryLoadBalancer{
		Balancer:     balancer,
		hostConsumer: map[string]*collect.ConsumerManager{},
		statistic:    both.NewMemoryStatistic(),
		threshold: &both.BaseThreshold{
			LastDuration: time.Second * 3,
			Rate:         0.8,
			Total:        10,
		},
		timeout: timeout,
	}
}

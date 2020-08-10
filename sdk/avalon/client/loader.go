package client

import (
	"errors"
	"github.com/941112341/avalon/sdk/avalon/both"
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/inline"
	"strings"
	"time"
)

type LoadBalancer interface {
	Balancer
	collect.Consumer
}

type RetryLoadBalancer struct {
	Balancer     Balancer
	hostConsumer map[string]collect.Consumer

	statistic both.RecordStatistic
	threshold both.Threshold
	timeout   time.Duration // 客户端超时时间

	min, max, backend int
	rate              float64
	total             int
	duration          time.Duration
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

	// 这里可以使用代理模式，偷个懒省略
	hostPort := inline.RandomList(normalHosts).(string)
	ip := inline.GetIP()
	if strings.Contains(hostPort, ip) {
		hostPort = strings.ReplaceAll(hostPort, ip, "localhost")
	}
	return hostPort, nil
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
			manager, err = NewPool(100*time.Millisecond, r.min, r.max, factory)
			if err != nil {
				return inline.PrependErrorFmt(err, "new pool fail")
			}
			r.hostConsumer[hostport] = manager

		}

		if err := manager.Consume(e); err != nil {
			r.statistic.AddErrLog(hostport, e, err)
			return inline.PrependErrorFmt(err, "consume err %s", hostport)
		}
		r.statistic.AddSuc(hostport)
		return nil
	}, 3, 0)
}

func (r *RetryLoadBalancer) Shutdown() error {
	for _, c := range r.hostConsumer {
		if err := c.Shutdown(); err != nil {
			return inline.PrependErrorFmt(err, "shutdown")
		}
	}
	return nil
}

type loadBalancerBuilder struct {
	b *RetryLoadBalancer
}

func NewLoadBalancerBuilder() *loadBalancerBuilder {
	return &loadBalancerBuilder{b: &RetryLoadBalancer{
		Balancer:     nil,
		hostConsumer: map[string]collect.Consumer{},
		statistic:    both.NewMemoryStatistic(),
		threshold:    &both.BaseThreshold{},
		timeout:      time.Second * 2,
		min:          10,
		max:          20,
		backend:      20,
		rate:         0.8,
		total:        20,
		duration:     time.Second * 10,
	}}
}

func (b *loadBalancerBuilder) Min(min int) *loadBalancerBuilder {
	b.b.min = min
	return b
}

func (b *loadBalancerBuilder) Max(max int) *loadBalancerBuilder {
	b.b.max = max
	return b
}

func (b *loadBalancerBuilder) Rate(rate float64) *loadBalancerBuilder {
	b.b.rate = rate
	return b
}

func (b *loadBalancerBuilder) Backend(backend int) *loadBalancerBuilder {
	b.b.backend = backend
	return b
}

func (b *loadBalancerBuilder) Total(total int) *loadBalancerBuilder {
	b.b.total = total
	return b
}

func (b *loadBalancerBuilder) Duration(duration time.Duration) *loadBalancerBuilder {
	b.b.duration = duration
	return b
}

func (b *loadBalancerBuilder) Statistic(statistic both.RecordStatistic) *loadBalancerBuilder {
	b.b.statistic = statistic
	return b
}

func (b *loadBalancerBuilder) Threshold(threshold both.Threshold) *loadBalancerBuilder {
	b.b.threshold = threshold
	return b
}

func (b *loadBalancerBuilder) Balancer(balancer Balancer) *loadBalancerBuilder {
	b.b.Balancer = balancer
	return b
}

func (b *loadBalancerBuilder) Build() LoadBalancer {
	return b.b
}

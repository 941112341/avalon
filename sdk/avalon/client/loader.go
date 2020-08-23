package client

import (
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
	"math/rand"
	"strings"
)

type LoadBalancer interface {
	avalon.Bean
	GetIP(ips []string) string
}

type RandomBalancer struct {
	avalon.TodoBean
}

func (r *RandomBalancer) GetIP(ips []string) string {
	if len(ips) == 0 {
		return ""
	}

	idx := rand.Intn(len(ips))
	str := ips[idx]
	ip := inline.GetIP()
	if strings.HasPrefix(str, ip) {
		return "localhost" + strings.TrimPrefix(str, ip)
	}
	return str
}

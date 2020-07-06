package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/pkg/errors"
	"math/rand"
)

type Discover interface {
	Get(serviceName string) (string, error)
	Put(serviceName, hostPort string)
}

type RandomDiscover struct {
	hostMap inline.SaveMapList
}

func (r *RandomDiscover) Get(serviceName string) (string, error) {
	v, ok := r.hostMap.Get(serviceName)
	if !ok {
		return "", errors.Errorf("serviceName %s not found", serviceName)
	}

	host, ok := v.([]interface{})
	if !ok {
		return "", errors.Errorf("serviceName %s, hostMap %s type not match", serviceName, inline.JsonString(v))
	}
	idx := rand.Intn(len(host))
	hostPortString, ok := host[idx].(string)
	if !ok {
		return "", errors.Errorf("host=%s, idx=%d is not string", inline.JsonString(host), idx)
	}
	return hostPortString, nil
}

func (r *RandomDiscover) Put(serviceName, hostPort string) {
	r.hostMap.Append(serviceName, hostPort)
}

var DiscoverInstance = RandomDiscover{hostMap: inline.SaveMapList{}}

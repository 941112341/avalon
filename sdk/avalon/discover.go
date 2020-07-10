package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/collect"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/zookeeper"
	"github.com/pkg/errors"
	"math/rand"
	"sync"
)

// key=serviceName; value=[]hostPort
var syncMap collect.SyncMap
var lock sync.Mutex

func init() {
	syncMap = *collect.NewSyncMap()
}

func DiscoverMiddleware(cfg *ClientConfig, call Call) Call {
	return func(ctx context.Context, method string, args, result interface{}) error {
		if cfg.ServiceName == "" {
			return errors.New("no service name")
		}

		if !syncMap.Contains(cfg.ServiceName) {
			lock.Lock()
			defer lock.Unlock()
			if !syncMap.Contains(cfg.ServiceName) {
				zkCli, err := zookeeper.GetZkClientInstance(&cfg.ZkConfig)
				if err != nil {
					return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
				}
				node := zookeeper.NewZkNodeBuilder(inline.JoinPath(cfg.Path, cfg.ServiceName)).Build()
				err = node.ListWL(zkCli, true)
				if err != nil {
					return err
				}
				syncMap.Put(cfg.ServiceName, node)
			}
		}

		i, _ := syncMap.Get(cfg.ServiceName)
		node := i.(*zookeeper.ZkNode)
		hostPorts := node.GetChildrenKey()
		if len(hostPorts) == 0 {
			return errors.New(cfg.ServiceName + " has none server")
		}
		idx := rand.Intn(len(hostPorts))

		cfg.HostPort = hostPorts[idx]
		return call(ctx, method, args, result)
	}
}

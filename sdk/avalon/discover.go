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

func DiscoverMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		if cfg.Client.HostPort != "" { // pass if set hostPort
			return call(ctx, method, args, result)
		}
		if cfg.Psm == "" {
			return errors.New("no service name")
		}

		if !syncMap.Contains(cfg.Psm) {
			lock.Lock()
			defer lock.Unlock()
			if !syncMap.Contains(cfg.Psm) {
				zkCli, err := zookeeper.GetZkClientInstance(cfg.ZkConfig)
				if err != nil {
					return errors.WithMessage(err, inline.ToJsonString(cfg.ZkConfig))
				}
				node := zookeeper.NewZkNodeBuilder(inline.JoinPath(cfg.ZkConfig.Path, cfg.Psm)).Build()
				err = node.ListWL(zkCli, true)
				if err != nil {
					return err
				}
				syncMap.Put(cfg.Psm, node)
			}
		}

		i, _ := syncMap.Get(cfg.Psm)
		node := i.(*zookeeper.ZkNode)
		hostPorts := node.GetChildrenKey()
		if len(hostPorts) == 0 {
			return errors.New(cfg.Psm + " has none server")
		}
		idx := rand.Intn(len(hostPorts))
		inline.Infoln("hostPort", inline.NewPair("hostPort", hostPorts[idx]))

		// set session
		SetHostPort(ctx, hostPorts[idx])
		return call(ctx, method, args, result)
	}
}

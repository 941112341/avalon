package avalon

import (
	"github.com/941112341/avalon/sdk/inline"
	"github.com/941112341/avalon/sdk/log"
	"github.com/pkg/errors"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
)

type server struct {
	start func() error
	stop  func() error
}

func (s *server) Start() error {
	return s.start()
}

func (s *server) Stop() error {
	return s.stop()
}

func ServiceRegisterWrapper(cfg *ServerConfig, coreServer Server) Server {
	return &server{
		start: func() error {

			servicePath := cfg.ZkConfig.Path + "/" + cfg.ServiceName
			exist, _, err := ZkClient.Conn.Exists(servicePath)
			if err != nil {
				return errors.Cause(err)
			}
			if !exist {
				str, err := ZkClient.Conn.Create(servicePath, []byte(""), 0, zk.WorldACL(zk.PermAll))
				log.New().WithField("path", str).WithField("err", err).Debugf("create service node\n")
			}

			hostPort := cfg.HostPort
			idx := strings.LastIndex(hostPort, ":")
			port := hostPort[idx:]
			ip, err := inline.GetIp()
			if err != nil {
				return errors.Cause(err)
			}

			completeHostPort := ip + port
			path := servicePath + "/" + completeHostPort

			exist, _, ch, err := ZkClient.Conn.ExistsW(path)
			if err != nil {
				return errors.WithMessage(err, "existsW "+path)
			}
			if !exist {
				log.New().Infof("create service node %s by path %s", completeHostPort, path)
				_, err = ZkClient.Conn.Create(path, []byte(""), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
				if err != nil {
					return errors.WithMessage(err, inline.JsonString(cfg))
				}
			} else {
				go func() {
					select {
					case event := <-ch:
						if event.Type == zk.EventNodeDeleted {
							_, err = ZkClient.Conn.Create(path, []byte(""), zk.FlagEphemeral, zk.WorldACL(zk.PermAll))
							if err != nil {
								log.New().WithField("err", err).WithField("path", path).Fatalln("reListening fail")
							}
						} else {
							log.New().WithField("event", event).Infoln("receive event")
						}
					}
				}()
			}
			return coreServer.Start()
		},
		stop: coreServer.Stop,
	}
}

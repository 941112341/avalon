package server

func DefaultServer() Bootstrap {

	return (&MyServer{}).AddBootstrap(&ThriftServer{}).
		AddBootstrapHook(&Zookeeper{}).
		AddWrapper(&LogWrapper{}).AddWrapper(&ErrorWrapper{})
}

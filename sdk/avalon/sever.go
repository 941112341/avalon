package avalon

type Server interface {
	Start(cfg *ServerConfig) error
	Stop(cfg *ServerConfig) error
}

type ServerWrapper func(server Server) Server

type Bootstrap struct {
	server   Server
	wrappers []ServerWrapper
}

func (b *Bootstrap) Start(cfg *ServerConfig) error {

	s := b.server
	for _, wrapper := range b.wrappers {
		s = wrapper(s)
	}
	return s.Start(cfg)
}

func (b *Bootstrap) Stop(cfg *ServerConfig) error {
	s := b.server
	for _, wrapper := range b.wrappers {
		s = wrapper(s)
	}
	return s.Stop(cfg)
}

func NewBootstrap(core Server) *Bootstrap {
	return &Bootstrap{
		server: core,
		wrappers: []ServerWrapper{
			ServiceRegisterWrapper,
		},
	}
}

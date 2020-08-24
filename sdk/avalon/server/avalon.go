package server

import (
	"github.com/941112341/avalon/sdk/avalon"
	"github.com/941112341/avalon/sdk/inline"
)

type ReactBootstrap interface {
	Bootstrap
	Omit(o interface{})
}

type Bootstrap interface {
	avalon.Bean
	Run(handler interface{}) error

	wrap(handler interface{}, wrapper avalon.Wrapper) interface{}
}

type BootstrapHook interface {
	avalon.Bean

	BeforeRun() error
}

// server implements
type MyServer struct {
	bootstrapHooks []BootstrapHook
	bootstraps     []Bootstrap
	handlers       avalon.WrapperComposite

	signal chan interface{}
}

func (MyServer MyServer) wrap(handler interface{}, wrapper avalon.Wrapper) interface{} {
	return handler
}

func (MyServer *MyServer) Omit(o interface{}) {
	MyServer.signal <- o
}

func (MyServer *MyServer) Initial() error {
	if err := avalon.InitialSlice(MyServer.bootstrapHooks); err != nil {
		return err
	}
	if err := avalon.InitialSlice(MyServer.bootstraps); err != nil {
		return err
	}
	if err := MyServer.handlers.Initial(); err != nil {
		return err
	}
	// 默认thriftAdapter todo

	MyServer.signal = make(chan interface{})
	return nil
}

func (MyServer *MyServer) Destroy() error {
	if err := avalon.DestroySlice(MyServer.bootstrapHooks); err != nil {
		return err
	}
	if err := avalon.DestroySlice(MyServer.bootstraps); err != nil {
		return err
	}
	if err := MyServer.handlers.Destroy(); err != nil {
		return err
	}
	return nil
}

func (MyServer *MyServer) Run(handler interface{}) error {
	if err := MyServer.Initial(); err != nil {
		return err
	}

	for _, hook := range MyServer.bootstrapHooks {
		if err := hook.BeforeRun(); err != nil {
			return err
		}
	}

	for _, bootstrap := range MyServer.bootstraps {

		go func(bootstrap Bootstrap) {

			if err := bootstrap.Run(bootstrap.wrap(handler, MyServer.handlers)); err != nil {
				MyServer.Omit(err)
			}
		}(bootstrap)
	}

signal:
	for o := range MyServer.signal {
		switch i := o.(type) {
		case error:
			inline.WithFields("err", i).Errorln("server run error")
			break signal
		case int:
			switch i {
			case 9:
				panic("force kill, panic!!")
			case 15:
				break
			}
		}
	}
	if err := MyServer.Destroy(); err != nil {
		inline.WithFields("err", err).Errorln("destroy fail")
	}
	return nil
}

func (MyServer *MyServer) AddBootstrapHook(hook BootstrapHook) *MyServer {
	MyServer.bootstrapHooks = append(MyServer.bootstrapHooks, hook)
	return MyServer
}

func (MyServer *MyServer) AddBootstrap(bootstrap Bootstrap) *MyServer {
	MyServer.bootstraps = append(MyServer.bootstraps, bootstrap)
	return MyServer
}

func (MyServer *MyServer) AddWrapper(handler avalon.Wrapper) *MyServer {
	MyServer.handlers = append(MyServer.handlers, handler)
	return MyServer
}

func DefaultServer() Bootstrap {

	return (&MyServer{}).AddBootstrap(&ThriftServer{}).
		AddBootstrapHook(&Zookeeper{}).
		AddWrapper(&ErrorWrapper{}).AddWrapper(&LogWrapper{})
}

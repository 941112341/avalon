package avalon

import "context"

type Call func(ctx context.Context, invoke *Invoke) error

type Invoke struct {
	MethodName string
	Request    interface{}
	Response   interface{}
}

type Wrapper interface {
	Bean
	Middleware(call Call) Call
}

// handler implements
type WrapperComposite []Wrapper

func (h WrapperComposite) Initial() error {
	return InitialSlice(h)
}

func (h WrapperComposite) Destroy() error {
	return DestroySlice(h)
}

func (h WrapperComposite) Middleware(call Call) Call {

	for _, wrapper := range h {
		call = wrapper.Middleware(call)
	}
	return call
}

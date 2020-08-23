package avalon

import "github.com/941112341/avalon/sdk/inline"

type Bean interface {
	Initial() error
	Destroy() error
}

type CompositeBeans []Bean

func (c CompositeBeans) Initial() error {
	for _, bean := range c {
		if err := bean.Initial(); err != nil {
			return err
		}
	}
	return nil
}

func (c CompositeBeans) Destroy() error {
	for _, bean := range c {
		if err := bean.Destroy(); err != nil {
			return err
		}
	}
	return nil
}

func InitialSlice(obj interface{}) error {
	return inline.RangeSlice(obj, func(i interface{}) error {
		if bean, ok := i.(Bean); ok {
			if err := NewBean(bean).Initial(); err != nil {
				return err
			}
		}
		return nil
	})

}

func DestroySlice(obj interface{}) error {
	return inline.RangeSlice(obj, func(i interface{}) error {
		if bean, ok := i.(Bean); ok {
			if err := NewBean(bean).Destroy(); err != nil {
				return err
			}
		}
		return nil
	})
}

type TodoBean struct {
}

func (t *TodoBean) Initial() error {
	return nil
}

func (t *TodoBean) Destroy() error {
	return nil
}

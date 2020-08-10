package both

type Hostport interface {
	Port() string
}

type CurrentHostPort struct {
	HostPort string
}

func (c CurrentHostPort) Hostport() string {
	return c.HostPort
}

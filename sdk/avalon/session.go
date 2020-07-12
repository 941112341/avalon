package avalon

import "context"

type Session struct {
	HostPort string
}

func GetSession(ctx context.Context) *Session {
	r, _ := ctx.Value("session").(*Session)
	return r
}

func WithSession(ctx context.Context, session *Session) context.Context {
	return context.WithValue(ctx, "session", session)
}

func GetHostPort(ctx context.Context) string {
	return GetSession(ctx).HostPort
}

func SetHostPort(ctx context.Context, hostPort string) {
	GetSession(ctx).HostPort = hostPort
}

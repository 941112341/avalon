package avalon

import "context"

type Session struct {
	HostPort string

	Attachments map[string]interface{}
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

func SetAttachment(ctx context.Context, key string, value interface{}) {
	GetSession(ctx).Attachments[key] = value
}

func GetAttachment(ctx context.Context, key string) interface{} {
	value, _ := GetSession(ctx).Attachments[key]
	return value
}

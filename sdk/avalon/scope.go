package avalon

import (
	"context"
	"github.com/941112341/avalon/sdk/inline"
	"github.com/json-iterator/go"
	"sync"
)

var application = &Scope{
	cache:     sync.Map{},
	parent:    nil,
	scopeType: FromApplication,
}

type Scope struct {
	cache     sync.Map
	parent    *Scope
	scopeType ScopeType
}

const (
	ScopeKey    = "__SCOPE_KEY__"
	HostPortKey = "HostPort"
)

type ScopeType int

const (
	_ ScopeType = iota
	FromSession
	FromCrossRPC
	FromApplication
)

func (s *Scope) find(scopeType ScopeType) *Scope {
	if s == nil {
		return nil
	}
	if s.scopeType == scopeType {
		return s
	}
	return s.parent.find(scopeType)
}

func (s *Scope) get(key string, scopeType ScopeType, extend bool) (string, bool) {
	scope := s.find(scopeType)
	if scope == nil {
		return "", false
	}
	ret, ok := scope.cache.Load(key)
	if ok {
		return ret.(string), ok
	}
	if !extend {
		return "", false
	}
	return scope.parent.get(key, scopeType, extend)
}

func (s *Scope) Get(key string, extend bool) (string, bool) {
	return s.get(key, FromSession, extend)
}

func (s *Scope) Set(key, value string, scopeType ScopeType) {
	if value == "" {
		return
	}
	scope := s.find(scopeType)
	if scope == nil {
		return
	}
	scope.cache.Store(key, value)
}

func GetScope(ctx context.Context) *Scope {
	scope, _ := ctx.Value(ScopeKey).(*Scope)
	return scope
}

func Get(ctx context.Context, key string) (string, bool) {
	scope := GetScope(ctx)
	return scope.Get(key, true)
}

func ConsistentSet(ctx context.Context, key, value string) {
	GetScope(ctx).Set(key, value, FromCrossRPC)
}

func scopeMiddleware(cfg Config, call Endpoint) Endpoint {
	return func(ctx context.Context, method string, args, result interface{}) error {
		scope := &Scope{
			cache: sync.Map{},
			parent: &Scope{
				cache:     sync.Map{},
				parent:    application,
				scopeType: FromCrossRPC,
			},
			scopeType: FromSession,
		}

		cfgStr, _ := jsoniter.Marshal(cfg)
		cfgMap := make(map[string]interface{})
		err := jsoniter.Unmarshal(cfgStr, &cfgMap)
		if err != nil {
			inline.Errorln("unmarshal err", inline.NewPairs("err", err.Error())...)
		}
		scope.Set(HostPortKey, cfg.Client.HostPort, FromCrossRPC)

		ctx = context.WithValue(ctx, ScopeKey, scope)
		return call(ctx, method, args, result)
	}
}

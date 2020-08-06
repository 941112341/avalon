package impl

import (
	"errors"
	"github.com/941112341/avalon/gateway/registry"
	"net/http"
)

const Token = "bbv14s5d4zs__xpf3bwpmszm0000gn"

func init() {
	_ = registry.Registry("Auth", &AuthController{
		RuleNeedMap: map[string]Authenticate{
			"Upload":   &HeaderChecker{Token: Token},
			"Registry": &HeaderChecker{Token: Token},
		},
	})
}

type AuthController struct {
	RuleNeedMap map[string]Authenticate
}

func (a *AuthController) CanAccess(method string, request *http.Request) error {
	authenticate, ok := a.RuleNeedMap[method]
	if !ok {
		return nil
	}

	return authenticate.Authenticate(request)
}

type Authenticate interface {
	Authenticate(request *http.Request) error
}

type HeaderChecker struct {
	Token string
}

func (h *HeaderChecker) Authenticate(request *http.Request) error {
	if authenticate := request.Header.Get("authenticate"); authenticate == "" {
		return errors.New("need authenticate")
	} else if authenticate != h.Token {
		return errors.New("token err")
	}
	return nil
}

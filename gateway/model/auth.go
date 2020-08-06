package model

import "net/http"

type AuthManager interface {
	CanAccess(method string, request *http.Request) error
}

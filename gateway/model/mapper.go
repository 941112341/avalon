package model

type MapperRules interface {
	GetApplication(request *HttpRequest) (*Application, error)
	AddRule(rule MapperRule) error
	Len() int
	Swap(i, j int)
	Less(i, j int) bool
}

type MapperRule interface {
	GetType() MapperRuleType
	Order() int
	Match(request *HttpRequest) (ApplicationKey, bool)
}

type MapperRuleType int32

const (
	Absolute MapperRuleType = iota
)

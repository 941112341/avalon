package impl

import (
	"fmt"
	"github.com/941112341/avalon/gateway/model"
	"github.com/941112341/avalon/gateway/repository"
	"github.com/941112341/avalon/sdk/inline"
	"math"
	"net/http"
	"regexp"
	"sort"
)

type Ruler struct {
	Rules []model.MapperRule
}

func (r *Ruler) GetApplication(request *http.Request) (model.Application, error) {
	sort.Sort(r)

	for _, rule := range r.Rules {
		key, ok := rule.Match(request)
		if !ok {
			continue
		}
		return key.GetApplication(), nil
	}

	return nil, fmt.Errorf("request not match %+v", request)
}

func (r *Ruler) Len() int {
	return len(r.Rules)
}

func (r *Ruler) Swap(i, j int) {
	r.Rules[i], r.Rules[j] = r.Rules[j], r.Rules[i]
}

func (r *Ruler) Less(i, j int) bool {
	return r.Rules[i].Order() < r.Rules[j].Order()
}

const (
	AbsoluteMatch model.MapperRuleType = iota
	PartialMatch
	RegexpMatch
)

type AbsRules struct {
	repository.MapperVo
}

func (r *AbsRules) GetType() model.MapperRuleType {
	return AbsoluteMatch
}

func (r *AbsRules) Order() int {
	return int(r.GetType())
}

func (r *AbsRules) Match(request *http.Request) (model.ApplicationKey, bool) {
	match := request.URL.Path == r.URL
	if !match {
		return nil, false
	}

	return &ExecutorFactory{
		ExecutorData: ExecutorData{
			psm:     r.PSM,
			version: r.Version,
			method:  r.Method,
			base:    r.Base,
		},
		MapperArgs: map[string]interface{}{},
	}, true
}

type PartialRule struct {
	repository.MapperVo
}

func (p *PartialRule) GetType() model.MapperRuleType {
	return PartialMatch
}

func (p *PartialRule) Order() int {
	return int(p.GetType())
}

func (p *PartialRule) Match(request *http.Request) (model.ApplicationKey, bool) {
	args, err := inline.TemplateExtract(p.URL, request.URL.Path)
	if err != nil {
		inline.WithFields("request", request, "rule", p).Errorln("extra fail")
		return nil, false
	}

	return &ExecutorFactory{
		ExecutorData: ExecutorData{
			psm:     p.PSM,
			version: p.Version,
			method:  p.Method,
			base:    p.Base,
		},
		MapperArgs: inline.MustStringIfaceMap(args),
	}, true
}

type RegexpRule struct {
	repository.MapperVo
}

func (r *RegexpRule) GetType() model.MapperRuleType {
	return RegexpMatch
}

func (r *RegexpRule) Order() int {
	return int(r.GetType())
}

func (r *RegexpRule) Match(request *http.Request) (model.ApplicationKey, bool) {
	pattern, err := regexp.Compile(r.URL)
	if err != nil {
		inline.WithFields("url", r.URL, "err", err).Errorln("compile err")
		return nil, false
	}

	args := inline.SubNameMatchMap(pattern, request.URL.Path)
	return &ExecutorFactory{
		ExecutorData: ExecutorData{
			psm:     r.PSM,
			version: r.Version,
			method:  r.Method,
			base:    r.Base,
		},
		MapperArgs: inline.MustStringIfaceMap(args),
	}, true
}

type NilMatchRule struct {
}

func (n *NilMatchRule) GetType() model.MapperRuleType {
	return math.MaxInt32
}

func (n *NilMatchRule) Order() int {
	return math.MaxInt32
}

func (n *NilMatchRule) Match(request *http.Request) (model.ApplicationKey, bool) {
	return nil, false
}

func NewRule(mapper repository.MapperVo) model.MapperRule {
	switch mapper.Type {
	case int16(AbsoluteMatch):
		return &AbsRules{mapper}
	case int16(PartialMatch):
		return &PartialRule{mapper}
	case int16(RegexpMatch):
		return &RegexpRule{mapper}
	default:
		return &NilMatchRule{}
	}
}

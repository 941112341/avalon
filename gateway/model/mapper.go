package model

import (
	"github.com/941112341/avalon/gateway/util"
	"github.com/941112341/avalon/sdk/inline"
	"regexp"
	"strings"
	"time"
)

type Request struct {
	Headers    map[string]string
	Body       string
	URL        string
	HTTPMethod string
}

type Response struct {
	HTTPCode int
	Headers  map[string]string
	Body     string
}

type MapperMatch struct {
	Mapper *Mapper
	URL    string
	Args   map[string]string
}

func (m *MapperMatch) valid() error {
	if m.Mapper == nil {
		return inline.NewError(inline.ErrArgs, "mapper nil")
	}
	return nil
}

func (m *MapperMatch) IDLFile(repo UploadRepository) (*IDLFile, error) {
	if err := m.valid(); err != nil {
		return nil, err
	}
	return repo.FindByKey(m.Mapper.IDLFileID)
}

func (m *MapperMatch) Transfer(repo UploadRepository, request Request) (*Response, error) {
	_, err := m.IDLFile(repo)
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "idlFile match %s", inline.ToJsonString(m))
	}
	return nil, nil
}

type MapperList struct {
	Mappers []*Mapper

	Absolute map[string]*Mapper
	Partial  []*Mapper
	Regexp   []*Mapper
}

func (m *MapperList) IsEmpty() bool {
	return len(m.Mappers) == 0
}

func (m *MapperList) All(repo MapperRepository) (*MapperList, error) {
	list, err := repo.AllMapper()
	if err != nil {
		return nil, inline.PrependErrorFmt(err, "all err")
	}
	for _, mapper := range list.Mappers {
		switch mapper.Type {
		case Absolute:
			list.Absolute[mapper.URL] = mapper
		case Partial:
			list.Partial = append(list.Partial, mapper)
		case Regexp:
			list.Regexp = append(list.Regexp, mapper)
		}
	}
	return list, nil
}

func (m *MapperList) Save(repo MapperRepository) error {
	return repo.AddMapper(m)
}

func (m *MapperList) Delete(repo MapperRepository) error {
	return repo.DelMapper(m)
}

func (m *MapperList) FindMapper(url string) (*MapperMatch, error) {
	url = util.StandardURL(url)
	mapper, ok := m.Absolute[url]
	if ok {
		return &MapperMatch{
			Mapper: mapper,
			URL:    url,
			Args:   map[string]string{},
		}, nil
	}

	var err error
	for _, pmap := range m.Partial {
		facts := strings.Split(pmap.URL, "/")
		expects := strings.Split(url, "/")
		if len(facts) != len(expects) {
			continue
		}

		length := len(expects)

		result := make(map[string]string)
		match := true // 全部循环完 为true
		for i := 0; i < length; i++ {
			fact, expect := facts[i], expects[i]
			if fact == expect {
				continue
			}

			result, err = inline.ParseFromTemplate(expect, fact, result)
			if err != nil {
				inline.WithFields("err", err).Debugln("not match fact %s", fact)
				match = false
				break
			}

		}

		if match {
			return &MapperMatch{
				Mapper: pmap,
				URL:    url,
				Args:   result,
			}, nil
		}
	}

	for _, r := range m.Regexp {
		pattern, err := regexp.Compile(r.URL)
		if err != nil {
			inline.WithFields("err", err).Errorln("compile err %s", r.URL)
		}
		matches := inline.SubNameMatchMap(pattern, url)
		if matches == nil {
			continue
		}
		return &MapperMatch{
			Mapper: mapper,
			URL:    url,
			Args:   matches,
		}, nil
	}

	return nil, inline.NewError(ErrNoMapperFound, "no match mapper found %s", url)
}

type MapperRule int16

const (
	Absolute MapperRule = iota
	Partial
	Regexp
)

type Mapper struct {
	ID      int64
	Deleted *bool
	URL     string
	Type    MapperRule
	IDLFileID
	Created time.Duration
	Updated time.Duration
	Method  string
}

func (Mapper) TableName() string {
	return "mapper"
}

func (m *Mapper) IDLFile(repo UploadRepository) (*IDLFile, error) {
	return repo.FindByKey(m.IDLFileID)
}

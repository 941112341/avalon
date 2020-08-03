package inline

import (
	"fmt"
	"github.com/json-iterator/go"
	"regexp"
	"strings"
	"unicode"
)

func ToJsonString(o interface{}) string {
	body, _ := jsoniter.MarshalToString(o)
	return body
}

func VString(o interface{}) string {
	return fmt.Sprintf("%+v", o)
}

func JoinPath(paths ...string) string {
	return strings.Join(paths, "/")
}

func String(o interface{}) string {
	return fmt.Sprintf("%s", o)
}

func Ucfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func Lcfirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}
	return ""
}

func IsEmpty(str string) bool {
	s := strings.Trim(str, " ")
	s = strings.Trim(str, "\t")
	return len(s) == 0
}

func ParseFromTemplate(template, fact string, result map[string]string) (map[string]string, error) {
	pattern := regexp.MustCompile(`{\w+}`)
	regexpTemplate := pattern.ReplaceAllStringFunc(template, func(s string) string {
		s = s[1 : len(s)-1]
		return fmt.Sprintf(`(?P<%s>\w+)`, s)
	})
	var err error
	pattern, err = regexp.Compile(regexpTemplate)
	if err != nil {
		return nil, PrependErrorFmt(err, "pattern = %s", ToJsonString(regexpTemplate))
	}

	m := SubNameMatchMap(pattern, fact)
	for k, v := range m {
		result[k] = result[v]
	}
	return result, nil
}

// return nil if un match
func SubNameMatchMap(r *regexp.Regexp, s string) map[string]string {
	matches := r.FindStringSubmatch(s)
	if len(matches) == 0 {
		return nil
	}
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i < len(matches) && i != 0 {
			result[name] = matches[i]
		}
	}
	return result
}

func SubNameMatchStruct(r *regexp.Regexp, s string, any interface{}) error {
	matches := r.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i < len(matches) && i != 0 {
			result[name] = matches[i]
		}
	}
	return Copy(result, any)
}

func JSONAny(any string) jsoniter.Any {
	return jsoniter.ParseString(jsoniter.ConfigDefault, any).ReadAny()
}

func Unwrap(r, content string) (s string) {
	pattern := regexp.MustCompile(r)
	ss := pattern.FindStringSubmatch(content)
	if len(ss) == 1 {
		s = ss[0]
	} else if len(ss) == 2 {
		s = ss[1]
	} else {
		return content
	}

	return strings.Trim(s, " ")
}

func Unwraps(r, content string) (ss []string) {
	pattern := regexp.MustCompile(r)
	ss = pattern.FindStringSubmatch(content)
	if len(ss) <= 0 {
		return
	}
	return ss[1:]
}

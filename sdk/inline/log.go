package inline

import (
	"errors"
	"fmt"
	"github.com/941112341/avalon/sdk/log"
	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

type Pair struct {
	Right string
	Left  string
}

func NewPair(a string, b interface{}) Pair {
	return Pair{Left: a, Right: VString(b)}
}

type Pairs []Pair

func NewPairs(args ...interface{}) Pairs {
	pairs := make(Pairs, 0)
	for i := 0; i+1 < len(args); i = i + 2 {
		key, ok := args[i].(string)
		if !ok {
			continue
		}
		pairs = append(pairs, NewPair(key, args[i+1]))
	}
	return pairs
}

func (p Pairs) Fields() logrus.Fields {
	fields := logrus.Fields{}
	for _, pair := range p {
		fields[pair.Left] = pair.Right
	}
	return fields
}

func (p Pairs) Infoln(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	Infoln(s, p...)
}

func (p Pairs) Debugln(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	Debugln(s, p...)
}

func (p Pairs) Warnln(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	Warningln(s, p...)
}

func (p Pairs) Errorln(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	p = append(p, stackPair()...)
	Errorln(s, p...)
}

func (p Pairs) String(msg string) string {
	m := map[string]interface{}{
		"__msg__": msg,
	}
	for _, pair := range p {
		m[pair.Left] = pair.Right
	}
	return ToJsonString(m)
}

func stackPair() []Pair {
	return []Pair{
		{
			Left:  "__local__",
			Right: RecordStack(2),
		},
		{
			Left:  "__parent__",
			Right: RecordStack(3),
		},
		{
			Left:  "grand",
			Right: RecordStack(4),
		},
	}
}

func (p Pairs) Fatalln(f string, args ...interface{}) {
	s := fmt.Sprintf(f, args...)
	p = append(p, stackPair()...)
	Fatalln(s, p...)
}

func Errorln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Errorln(msg)
	sentry.CaptureException(errors.New(Pairs(pairs).String(msg)))
}

func Infoln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Infoln(msg)
	sentry.CaptureMessage(Pairs(pairs).String(msg))
}

func Debugln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Debugln(msg)
}

func Warningln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Warningln(msg)
	sentry.CaptureException(errors.New(Pairs(pairs).String(msg)))
}

func Fatalln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Fatalln(msg)
	sentry.CaptureException(errors.New(Pairs(pairs).String(msg)))
}

func WithFields(args ...interface{}) Pairs {
	return NewPairs(args...)
}

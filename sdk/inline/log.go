package inline

import (
	"github.com/941112341/avalon/sdk/log"
	"github.com/sirupsen/logrus"
)

type Pair struct {
	Right string
	Left  interface{}
}

func NewPair(a string, b interface{}) Pair {
	return Pair{a, b}
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
		fields[pair.Right] = pair.Left
	}
	return fields
}

func Errorln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Errorln(msg)
}

func Infoln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Infoln(msg)
}

func Debugln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Debugln(msg)
}

func Warningln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Warningln(msg)
}

func Fatalln(msg string, pairs ...Pair) {
	log.New().WithFields(Pairs(pairs).Fields()).Fatalln(msg)
}

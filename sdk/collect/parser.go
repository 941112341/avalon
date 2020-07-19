package collect

import "regexp"

type Parser interface {
	ReadUntil(str string)
	ReadUntilFunc(func(str string) bool)
	HasNext() bool
	ReadLine() string
	NextMatch(regexp regexp.Regexp) Parser

	Reset()
}

type StringParser struct {
	Content []rune
	Index   int
}

func (s *StringParser) ReadUntil(str string) {
	panic("implement me")
}

func (s *StringParser) ReadUntilFunc(f func(str string) bool) {
	for i := s.Index; i < len(s.Content); i++ {
		if f(s.Content[i]) {

		}
	}
}

func (s StringParser) HasNext() bool {
	panic("implement me")
}

func (s StringParser) ReadLine() string {
	panic("implement me")
}

func (s StringParser) NextMatch(regexp regexp.Regexp) Parser {
	panic("implement me")
}

func (s StringParser) Reset() {
	panic("implement me")
}

package kana

import (
	"strings"
	"unicode/utf8"
)

type stream struct {
	buf  []rune
	end  bool
	next func(buf *[]rune)
}

func (s *stream) fill(demand int) {
	if s.end {
		return
	}
	for len(s.buf) < demand {
		oldSize := len(s.buf)
		s.next(&s.buf)
		if len(s.buf) == oldSize {
			s.end = true
			break
		}
	}
}

func (s *stream) consume(num int) {
	newSize := len(s.buf) - num
	for i := 0; i < newSize; i++ {
		s.buf[i] = s.buf[i+num]
	}
	s.buf = s.buf[:newSize]
}

func (s *stream) readOne() (rune, bool) {
	s.fill(1)
	if len(s.buf) == 0 {
		return 0, false
	}
	ch := s.buf[0]
	s.consume(1)
	return ch, true
}

func (s *stream) peekOne() (rune, bool) {
	s.fill(1)
	if len(s.buf) == 0 {
		return 0, false
	}
	return s.buf[0], true
}

func (s *stream) readAll() string {
	builder := strings.Builder{}
	if !s.end {
		for {
			s.readCurrentTo(&builder)

			oldSize := len(s.buf)
			s.next(&s.buf)
			if len(s.buf) == oldSize {
				s.end = true
				break
			}
		}
	}
	s.readCurrentTo(&builder)
	return builder.String()
}
func (s *stream) readCurrentTo(builder *strings.Builder) {
	for _, ch := range s.buf {
		builder.WriteRune(ch)
	}
	s.buf = s.buf[:0]
}

func newStream(next func(buf *[]rune)) *stream {
	return &stream{
		buf:  nil,
		end:  false,
		next: next,
	}
}

func stringStream(s string) *stream {
	pos := 0
	return newStream(func(buf *[]rune) {
		if pos < len(s) {
			ch, size := utf8.DecodeRuneInString(s[pos:])
			*buf = append(*buf, ch)
			pos += size
		}
	})
}

func mapStream(s *stream, f func(rune) rune) *stream {
	return newStream(func(buf *[]rune) {
		s.fill(1)
		for _, ch := range s.buf {
			*buf = append(*buf, f(ch))
		}
		s.consume(len(s.buf))
	})
}

package util

import "bytes"

type runeBuffer struct {
	buf *bytes.Buffer
}

func NewRuneBuffer() *runeBuffer {
	return &runeBuffer{
		buf: &bytes.Buffer{},
	}
}

func (rb *runeBuffer) WriteRune(r rune) *runeBuffer {
	rb.buf.WriteRune(r)
	return rb
}

func (rb *runeBuffer) WriteString(s string) *runeBuffer {
	rb.buf.WriteString(s)
	return rb
}

func (rb *runeBuffer) RemoveLastRune() *runeBuffer {
	s := rb.buf.String()
	rs := []rune(s)
	now := rs[:len(rs)-1]
	rb.buf.Truncate(len([]byte(string(now))))
	return rb
}

func (rb *runeBuffer) String() string {
	return rb.buf.String()
}

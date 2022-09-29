package util

import "testing"

func TestBuf(t *testing.T) {
	rb := NewRuneBuffer()
	rb.WriteString("你好世界").RemoveLastRune()
	if rb.String() != "你好世" {
		t.Errorf("expected: 你好世, actual: %s", rb.String())
	}
}

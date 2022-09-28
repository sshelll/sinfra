package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	f, err := Parse("./mock_test.go")
	assert.Nil(t, err)
	f.Print()
}

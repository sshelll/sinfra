package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile(t *testing.T) {
	assert.True(t, IsFileExist("file.go"))
	assert.False(t, IsFileExist("file.go1"))
	assert.True(t, IsDirExist("../util"))
	assert.False(t, IsDirExist("./util"))
	assert.NotNil(t, LastModTime("file.go"))
}

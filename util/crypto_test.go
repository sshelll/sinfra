package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAesEncrypt(t *testing.T) {
	origData := []byte("hello, world!")
	key := []byte("0123456789abcdef")

	crypted, err := AesEncrypt(origData, key)
	assert.Nil(t, err)

	decrypted, err := AesDecrypt(crypted, key)
	assert.Nil(t, err)
	assert.Equal(t, origData, decrypted)
}

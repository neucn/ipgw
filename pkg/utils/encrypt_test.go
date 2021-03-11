package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDESEncryptAndDecrypt(t *testing.T) {
	a := assert.New(t)
	type testCase struct {
		src    []byte
		secret []byte
	}
	testCases := []testCase{
		{src: []byte("hello_world"), secret: []byte{}},
		{src: []byte("hello_world"), secret: []byte("secret")},
	}
	for _, c := range testCases {
		dist, err := DESEncrypt(c.src, c.secret)
		a.Nil(err)
		decrypted, err := DESDecrypt(dist, c.secret)
		a.Nil(err)
		a.Equal(decrypted, c.src)

		decrypted, err = DESDecrypt(dist, []byte("x"))
		a.NotNil(err)
	}
}

func TestBase64EncodeAndDecode(t *testing.T) {
	a := assert.New(t)
	testCases := []string{
		"hello", "world",
	}
	for _, c := range testCases {
		dist := Base64Encode([]byte(c))
		a.NotZero(len(dist))
		decoded, err := Base64Decode(dist)
		a.Nil(err)
		a.Equal(string(decoded), c)
	}
}

func TestEncryptAndDecrypt(t *testing.T) {
	a := assert.New(t)
	type testCase struct {
		src    []byte
		secret []byte
	}
	testCases := []testCase{
		{src: []byte("hello_world"), secret: []byte{}},
		{src: []byte("hello_world"), secret: []byte("secret")},
	}
	for _, c := range testCases {
		dist, err := Encrypt(c.src, c.secret)
		a.Nil(err)
		decrypted, err := Decrypt(dist, c.secret)
		a.Nil(err)
		a.Equal(decrypted, string(c.src))
	}
}

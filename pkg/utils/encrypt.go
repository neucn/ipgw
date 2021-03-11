package utils

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/forgoer/openssl"
)

func sha(src []byte) []byte {
	h := sha1.New()
	h.Write(src)
	result := h.Sum(nil)
	return []byte(fmt.Sprintf("%x", result))
}

func DESEncrypt(source []byte, secret []byte) ([]byte, error) {
	key := sha(secret)
	result, err := openssl.Des3CBCEncrypt(source, key[:24], key[32:], openssl.PKCS7_PADDING)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func DESDecrypt(encrypted []byte, secret []byte) (result []byte, err error) {
	defer func() {
		if e := recover(); e != nil {
			err = errors.New("wrong secret")
		}
	}()
	key := sha(secret)
	result, err = openssl.Des3CBCDecrypt(encrypted, key[:24], key[32:], openssl.PKCS7_PADDING)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func Base64Encode(source []byte) string {
	return base64.StdEncoding.EncodeToString(source)
}

func Base64Decode(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

func Encrypt(source []byte, secret []byte) (string, error) {
	encrypt, err := DESEncrypt(source, secret)
	if err != nil {
		return "", err
	}
	return Base64Encode(encrypt), nil
}

func Decrypt(encoded string, secret []byte) (string, error) {
	encrypted, err := Base64Decode(encoded)
	if err != nil {
		return "", err
	}
	decrypted, err := DESDecrypt(encrypted, secret)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

package model

import (
	"errors"

	"github.com/neucn/ipgw/pkg/utils"
)

type Account struct {
	Username          string `json:"username"`
	Password          string `json:"-"`
	Cookie            string `json:"-"`
	Secret            string `json:"-"`
	EncryptedPassword string `json:"encrypted_password"`
}

func (a *Account) GetPassword() (string, error) {
	if a.Password != "" {
		return a.Password, nil
	}
	if a.EncryptedPassword == "" {
		return "", errors.New("no password stored")
	}
	result, err := utils.Decrypt(a.EncryptedPassword, []byte(a.Secret))
	if err != nil {
		return "", err
	}
	a.Password = result
	return result, nil
}

func (a *Account) SetPassword(password string, secret []byte) error {
	result, err := utils.Encrypt([]byte(password), secret)
	if err != nil {
		return err
	}
	a.EncryptedPassword = result
	return nil
}

func (a *Account) String() string {
	return a.Username
}

package model

import (
	"fmt"
	"github.com/neucn/ipgw/pkg/utils"
)

type Config struct {
	DefaultAccount string     `json:"default_account"`
	Accounts       []*Account `json:"accounts"`
}

func (c *Config) AddAccount(username, password, secret string) error {
	for _, account := range c.Accounts {
		if account.Username == username {
			return fmt.Errorf("account %s already exists \n", username)
		}
	}
	encryptedPassword, err := utils.Encrypt([]byte(password), []byte(secret))
	if err != nil {
		return err
	}
	c.Accounts = append(c.Accounts, &Account{
		Username:          username,
		EncryptedPassword: encryptedPassword,
	})
	return nil
}

func (c *Config) GetAccount(username string) *Account {
	for _, account := range c.Accounts {
		if account.Username == username {
			return account
		}
	}
	return nil
}

func (c *Config) DelAccount(username string) bool {
	if c.DefaultAccount == username {
		c.DefaultAccount = ""
	}
	for i, account := range c.Accounts {
		if account.Username == username {
			c.Accounts = append(c.Accounts[:i], c.Accounts[i+1:]...)
			return true
		}
	}
	return false
}

func (c *Config) SetDefaultAccount(username string) bool {
	for _, account := range c.Accounts {
		if account.Username == username {
			c.DefaultAccount = username
			return true
		}
	}
	return false
}

func (c *Config) GetDefaultAccount() *Account {
	if c.DefaultAccount != "" {
		return c.GetAccount(c.DefaultAccount)
	}
	if len(c.Accounts) > 0 {
		return c.Accounts[0]
	}
	return nil
}

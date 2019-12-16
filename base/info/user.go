package info

import (
	"net/http"
)

type UserInfo struct {
	Username string
	Password string
	Cookie   *http.Cookie
}

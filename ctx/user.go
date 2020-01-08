package ctx

import (
	"net/http"
)

type User struct {
	Username string
	Password string
	// 一网通的Cookie
	Cookie *http.Cookie
}

func (u *User) SetCookie(c string) {
	u.Cookie = &http.Cookie{
		Name:   "CASTGC",
		Value:  c,
		Domain: "pass.neu.edu.cn",
		Path:   "/tpass/",
	}
}

package info

import (
	"net/http"
)

type UserInfo struct {
	Username string
	Password string
	Cookie   *http.Cookie
	CAS      *http.Cookie
}

func (u *UserInfo) SetCookie(c string) {
	u.Cookie = &http.Cookie{
		Name:   "session_for%3Asrun_cas_php",
		Value:  c,
		Domain: "ipgw.neu.edu.cn",
	}
}

func (u *UserInfo) SetCAS(c string) {
	u.CAS = &http.Cookie{
		Name:   "CASTGC",
		Value:  c,
		Domain: "pass.neu.edu.cn",
		Path:   "/tpass/",
	}
}

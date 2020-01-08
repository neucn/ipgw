package ctx

import (
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

var (
	FullView bool

	t *Ctx
	e sync.Once
)

// 单例获取Ctx对象
func DefaultCtx() *Ctx {
	e.Do(func() {
		t = NewCtx()
	})
	return t
}

// 无参获取新client对象
func NewClient() *http.Client {
	n := &http.Client{Timeout: 3 * time.Second}
	jar, _ := cookiejar.New(nil)
	// 绑定session
	n.Jar = jar
	return n
}

// 无参获取新ctx对象
func NewCtx() *Ctx {
	n := &Ctx{
		User:   &User{},
		Net:    &Net{},
		Client: NewClient(),
		Option: &Option{},
	}
	// 初始化避免空指针
	n.User.SetCookie("")
	n.Net.SetCookie("")
	return n
}

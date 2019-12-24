package ctx

import (
	"ipgw/base/info"
	"net/http"
	"net/http/cookiejar"
	"sync"
	"time"
)

var (
	c *http.Client
	o sync.Once

	t *Ctx
	e sync.Once
)

// 单例获取client对象
func GetClient() *http.Client {
	o.Do(func() {
		// 3秒超时
		c = &http.Client{Timeout: 3 * time.Second}
		jar, _ := cookiejar.New(nil)
		// 绑定session
		c.Jar = jar
	})
	return c
}

// 单例获取Ctx对象
func GetCtx() *Ctx {
	e.Do(func() {
		t = &Ctx{
			User: &info.UserInfo{},
			Net:  &info.NetInfo{},
		}
	})
	return t
}

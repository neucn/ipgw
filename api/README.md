# 目录
- [简介](#简介)

- [API设计](#API设计)

- [API说明](#API说明)

    - [Login](#Login)
    
    - [Proxy](#Proxy)
    
- [返回码](#返回码)

- [SDK](#SDK)

# 简介

当前API版本: `v1`

目前提供两个API:

- `login` 获取登陆后的`CASTGC`值
- `proxy` 代为请求指定页面并返回页面内容

# API设计

通过`ipgw api <version> <command> [arguments]`调用相应的api

`ipgw`会判断api的版本是否兼容，若不兼容则返回错误

# API说明

## Login

调用命令:
```shell script
ipgw api v1 login
```
可用参数
```
-u username -p password -c cookie -v webvpn
```


> 当`u`, `p`, `c`同时给出时，优先使用`u`, `p`登陆
>
> 当`u`, `p`, `c`未给出或给出空值时，使用本地保存的账号登陆
>
> cookie 为一网通平台cookie中`CASTGC`键对应的**值**
>
> *该api提供cookie参数仅为验证cookie是否失效*

## Proxy

调用命令
```shell script
ipgw api v1 proxy
```
可用参数
```
-u username -p password -c cookie

-s service-url -m method -h headers -b body
```

> `u`, `p`, `c`的规则同`login`命令
>
> 若未指定`m`，则默认为`GET`
>
> 无需指定是否使用webvpn，`ipgw`可根据`s`自动判断
>
> `h`的值应为`http.Header`对象使用`json.Marshal()`序列化后的值，即`map[string][]string`类型的json形式
>
> `b`为请求携带的请求体


# 返回码

```
	0 成功
	1 无对应命令
	2 调用时没有给出参数
	3 API版本不兼容
	4 读入本地配置失败
	5 网络错误

	11 使用账号登陆时未指定密码
	12 本地无已保存账号
	13 账号或密码错误
	14 Cookie已失效
	15 账号被Ban 或 服务未授权
	16 登陆失败

	21 未指定service url
```

# SDK

> ipgw.go
```go
package ipgw

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 登陆
func Login(config *LoginConfig) (cookie string, code int) {
	if config.Webvpn {
		return execCommand("ipgw", "api", "v1", "login", "-u", config.User.Username, "-p", config.User.Password, "-c", config.User.Cookie, "-v")
	}
	return execCommand("ipgw", "api", "v1", "login", "-u", config.User.Username, "-p", config.User.Password, "-c", config.User.Cookie)
}

// 代理请求
func Proxy(config *ProxyConfig) (body string, code int) {
	headers, _ := json.Marshal(config.Headers)
	return execCommand("ipgw", "api", "v1", "proxy",
		"-u", config.User.Username, "-p", config.User.Password, "-c", config.User.Cookie,
		"-s", config.ServiceUrl,
		"-m", config.Method,
		"-h", string(headers),
		"-b", config.Body)
}

// 创建客户端
func NewCasClient(cookie string, webvpn bool) (client *http.Client) {
	n := &http.Client{Timeout: 6 * time.Second}
	jar, _ := cookiejar.New(nil)
	// 绑定session
	n.Jar = jar
	if webvpn {
		jar.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "pass-443.webvpn.neu.edu.cn",
			Path:   "/tpass/",
		}, []*http.Cookie{
			{
				Name:   "CASTGC",
				Value:  cookie,
				Domain: "pass-443.webvpn.neu.edu.cn",
				Path:   "/tpass/",
			},
		})
	} else {
		jar.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "pass.neu.edu.cn",
			Path:   "/tpass/",
		}, []*http.Cookie{
			{
				Name:   "CASTGC",
				Value:  cookie,
				Domain: "pass.neu.edu.cn",
				Path:   "/tpass/",
			},
		})
	}
	return n
}

// 当调用命令失败时返回-1。
// 当命令返回错误码非0时返回空字符串与错误码。
// 否则返回调用结果与0.
func execCommand(name string, params ...string) (result string, code int) {
	var outbuf, errbuf bytes.Buffer

	cmd := exec.Command(name, params...)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	if err != nil && !strings.Contains(err.Error(), "exit status") {
		return "", -1
	}
	result = string(outbuf.Bytes())
	c := errbuf.Bytes()
	c = c[:len(c)-1]
	code, _ = strconv.Atoi(string(c))
	return
}

```

> config.go
```go
package ipgw

import "net/http"

type User struct {
	Username string
	Password string
	Cookie   string
}

type LoginConfig struct {
	User   *User
	Webvpn bool
}

type ProxyConfig struct {
	User       *User
	ServiceUrl string
	Method     string
	Headers    *http.Header
	Body       string
}

func NewLoginConfig() *LoginConfig {
	return &LoginConfig{User: &User{}}
}

func NewProxyConfig() *ProxyConfig {
	return &ProxyConfig{User: &User{}, Headers: &http.Header{}, Method: "GET"}
}

```

> code.go
```go
package ipgw

var (
	codeMap = map[int]string{
		3: "与ipgw版本不兼容，请升级",
		4: "读取本地配置失败",
		5: "网络错误，请检查网络",

		11: "未指定密码",
		12: "无已保存账号，请指定账号密码",
		13: "账号或密码错误",
		14: "Cookie已失效",
		15: "账户被禁或服务未授权",
		16: "登陆失败，请重试",
	}
)

func CodeText(code int) string {
	return codeMap[code]
}

```
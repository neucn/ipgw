package api

import (
	"encoding/base64"
	"io/ioutil"
	"ipgw/base"
	"ipgw/ctx"
	. "ipgw/lib"
	"net/http"
	"strings"
)

// 发送请求，错误返回 0 -3
func sendRequest(c *ctx.Ctx, r *http.Request) {
	resp, err := c.Client.Do(r)
	if err != nil {
		Fatal(globalNetError)
	}
	c.Response = resp
}

// 载入一网通账号和密码，错误返回 0 -2
func LoadUser(c *ctx.Ctx) {
	// 准备读取
	path, err := GetPath(base.SavePath)
	if err != nil {
		Fatal(globalFailLoad)
	}

	// 读取
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		Fatal(globalFailLoad)
	}
	content := string(bytes)

	// 分割
	lines := strings.Split(content, LineDelimiter)
	if len(lines) < 2 {
		Fatal(globalFailLoad)
	}

	// 载入用户信息部分
	user := strings.Split(lines[0], PartDelimiter)
	if len(user) < 3 {
		Fatal(globalFailLoad)
	}

	// [b64(username), b64(password), CAS Cookie]
	username, err := base64.StdEncoding.DecodeString(user[0])
	c.User.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(user[1])
	c.User.Password = string(password)
}

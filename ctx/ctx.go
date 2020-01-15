// 存放上下文相关的结构体、变量与函数

package ctx

import (
	"bufio"
	"encoding/base64"
	"io/ioutil"
	. "ipgw/base"
	. "ipgw/lib"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Ctx struct {
	User     *User
	Net      *Net
	Client   *http.Client
	Response *http.Response
	Option   *Option
}

// b64(username):b64(password):CAS Cookie
// SID:ipgw Cookie

// 从配置文件中解析出用户配置
func (i *Ctx) Load() {
	InfoL(infoLoading)
	// 准备读取
	path, err := GetConfigPath(SavePath)
	if err != nil {
		FatalL(fatalGetPath)
	}

	// 读取
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		FatalL(fatalLoadInfo)
	}
	content := string(bytes)

	// 初次使用
	if len(content) < 1 {
		content = "::\n:"
	}

	// 分割
	lines := strings.Split(content, LineDelimiter)
	if len(lines) < 2 {
		FatalL(fatalInfoFormat)
	}

	// 载入用户信息部分
	user := strings.Split(lines[0], PartDelimiter)
	if len(user) < 3 {
		FatalL(fatalInfoFormat)
	}

	// [b64(username), b64(password), CAS Cookie]
	username, err := base64.StdEncoding.DecodeString(user[0])
	i.User.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(user[1])
	i.User.Password = string(password)

	i.User.SetCookie(user[2])

	// 载入网关信息部分
	net := strings.Split(lines[1], PartDelimiter)
	if len(net) < 2 {
		FatalL(fatalInfoFormat)
	}

	// [SID, ipgw Cookie]
	i.Net.SID = net[0]

	i.Net.SetCookie(net[1])
}

// 向配置文件中写入用户配置
func (i *Ctx) SaveAll() {
	InfoL(infoSaving)
	// 准备存储
	path, err := GetConfigPath(SavePath)
	if err != nil {
		FatalL(fatalGetPath)
	}

	// 打开
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		FatalL(fatalOpenFile)
	}

	// 写入
	username := base64.StdEncoding.EncodeToString([]byte(i.User.Username))
	password := base64.StdEncoding.EncodeToString([]byte(i.User.Password))
	casCookie := i.User.Cookie.Value

	sid := i.Net.SID
	netCookie := i.Net.Cookie.Value

	// 如果保存账号
	_, err = f.WriteString(username + PartDelimiter + password + PartDelimiter + casCookie + LineDelimiter +
		sid + PartDelimiter + netCookie)

	if err != nil {
		FatalL(fatalSaveInfo)
	}

	// 输出成功提示
	Info(infoSaved)
}

func (i *Ctx) SaveSession() {
	// 静默式，不需要输出
	path, err := GetConfigPath(SavePath)
	if err != nil {
		os.Exit(2)
	}

	// 读取
	bytes, _ := ioutil.ReadFile(path)
	content := string(bytes)

	// 分割
	lines := strings.Split(content, LineDelimiter)

	// 可能格式有误，静默兼容
	if len(lines) < 2 {
		lines = []string{"::", ":"}
	}

	// 保存一网通Cookie
	user := strings.Split(lines[0], PartDelimiter)
	if len(user) < 3 {
		os.Exit(2)
	}

	user[2] = i.User.Cookie.Value

	// 保存网关SID与Cookie
	net := strings.Split(lines[1], PartDelimiter)
	if len(net) < 2 {
		os.Exit(2)
	}

	net[1] = i.Net.Cookie.Value

	// 保存如文件
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		_ = f
	}()
	w := bufio.NewWriter(f)

	_, _ = w.WriteString(strings.Join(user, PartDelimiter) + LineDelimiter +
		strings.Join(net, PartDelimiter))
	_ = w.Flush()
}

func (i *Ctx) ExtractNetCookie() {
	cookie := i.Client.Jar.Cookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	})
	if len(cookie) == 0 {
		ErrorL(errGetCookie)
	} else {
		i.Net.Cookie = cookie[0]
		if !i.Option.Mute && FullView {
			InfoF(successGetCookie, i.Net.Cookie.Value)
		}
	}
}

func (i *Ctx) ExtractUserCookie() {
	cookie := i.Client.Jar.Cookies(&url.URL{
		Scheme: "https",
		Host:   "pass.neu.edu.cn",
		Path:   "/tpass/",
	})
	for _, cas := range cookie {
		if cas.Name == "CASTGC" {
			i.User.Cookie = cas
			break
		}
	}
}

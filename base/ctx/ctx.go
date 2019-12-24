package ctx

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/file"
	"ipgw/base/info"
	"os"
	"strings"
)

type Ctx struct {
	User *info.UserInfo
	Net  *info.NetInfo
	UA   string
}

// todo 暂时只考虑单用户，格式如下

// cookie:sid:ip\n
// b64(username):b64(password)

// 从配置文件中解析出用户配置
func (i *Ctx) Load() {
	fmt.Println(tipLoadInfo)
	// 准备读取
	path, err := file.GetPath(cfg.SavePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, errorGetPath)
		os.Exit(2)
	}

	// 读取
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, errorLoadInfo)
		os.Exit(2)
	}
	content := string(bytes)

	// 分割
	lines := strings.Split(content, file.LineDelimiter)
	if len(lines) < 3 {
		fmt.Fprintln(os.Stderr, errorInfoFormat)
		os.Exit(2)
	}

	// 载入用户信息部分
	user := strings.Split(lines[0], file.PartDelimiter)
	if len(user) < 2 {
		fmt.Fprintln(os.Stderr, errorInfoFormat)
		os.Exit(2)
	}

	username, err := base64.StdEncoding.DecodeString(user[0])
	i.User.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(user[1])
	i.User.Password = string(password)

	// 载入session部分
	session := strings.Split(lines[1], file.PartDelimiter)
	if len(session) < 3 {
		fmt.Fprintln(os.Stderr, errorInfoFormat)
		os.Exit(2)
	}

	cookie := session[0]
	i.User.SetCookie(cookie)

	sid := session[1]
	i.Net.SID = sid

	ip := session[2]
	i.Net.IP = ip

	cas := lines[2]
	i.User.SetCAS(cas)
}

// 向配置文件中写入用户配置
func (i *Ctx) SaveAll() {
	fmt.Println(tipSaveInfo)
	// 准备存储
	path, err := file.GetPath(cfg.SavePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, errorGetPath)
		os.Exit(2)
	}

	// 打开
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	defer func() {
		_ = f.Close()
	}()
	w := bufio.NewWriter(f)

	if err != nil {
		fmt.Fprintln(os.Stderr, errorOpenFile)
		os.Exit(2)
	}

	// 写入
	username := base64.StdEncoding.EncodeToString([]byte(i.User.Username))
	password := base64.StdEncoding.EncodeToString([]byte(i.User.Password))
	cookie := i.User.Cookie.Value
	cas := i.User.CAS.Value
	sid := i.Net.SID
	ip := i.Net.IP

	// 如果保存账号
	_, err = w.WriteString(username + file.PartDelimiter + password + file.LineDelimiter +
		cookie + file.PartDelimiter + sid + file.PartDelimiter + ip + file.LineDelimiter +
		cas)

	if err != nil {
		fmt.Fprintln(os.Stderr, errorSaveInfo)
		os.Exit(2)
	}
	_ = w.Flush()

	// 输出成功提示
	fmt.Println(successSaveInfo)
}

func (i *Ctx) SaveSession() {
	// 静默式，不需要输出
	path, err := file.GetPath(cfg.SavePath)
	if err != nil {
		os.Exit(2)
	}

	// 读取
	bytes, _ := ioutil.ReadFile(path)
	content := string(bytes)

	// 分割
	lines := strings.Split(content, file.LineDelimiter)
	// todo 后期多用户模式用下面的↓
	/*for ; len(lines) < 2; {
		lines = append(lines, "::")
	}*/
	if len(lines) < 3 {
		lines = []string{":", "::", ""}
	}

	session := strings.Split(lines[1], file.PartDelimiter)
	if len(session) < 3 {
		os.Exit(2)
	}

	session[0] = i.User.Cookie.Value
	session[1] = i.Net.SID
	session[2] = i.Net.IP

	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		_ = f
	}()
	w := bufio.NewWriter(f)

	_, _ = w.WriteString(lines[0] + file.LineDelimiter +
		strings.Join(session, file.PartDelimiter) + file.LineDelimiter +
		i.User.CAS.Value)
	_ = w.Flush()
}

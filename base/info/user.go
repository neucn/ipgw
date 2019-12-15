package info

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"ipgw/base/file"
	"net/http"
	"os"
	"strings"
)

type UserInfo struct {
	Username string
	Password string
	Cookie   *http.Cookie
}

// 格式化输出用户的所有信息
func (i *UserInfo) Print() {
	fmt.Println(
		"===========Info===========" + "\n" +
			"Username: \t" + i.Username + "\n" +
			//"IP: \t\t" + i.IP + "\n" +
			//"SID: \t\t" + i.SID + "\n" +
			//"Balance: \t" + fmt.Sprintf("%.2f", i.Balance) + " RMB\n" +
			//"Used: \t\t" + getUsedFlux(i.Used) + "\n" +
			//"Time: \t\t" + getUsedTime(i.Time) + "\n" +
			"==========================")
}

// todo 先实现单用户版本，格式如下

// b64(username):b64(password):cookie\n

// 从配置文件中解析出用户配置
func (i *UserInfo) Load(path string) {
	fmt.Println(tipLoadInfo)
	// 准备读取
	path, err := file.GetPath(path)
	if err != nil {
		fmt.Fprint(os.Stderr, errorGetPath)
		os.Exit(2)
	}

	// 读取
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Fprint(os.Stderr, errorLoadInfo)
		os.Exit(2)
	}
	content := string(bytes)

	// 分割
	split := strings.Split(content, file.Delimiter)
	if len(split) != 3 {
		fmt.Fprint(os.Stderr, errorInfoFormat)
		os.Exit(2)
	}

	username, err := base64.StdEncoding.DecodeString(split[0])
	i.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(split[1])
	i.Password = string(password)

	cookie := split[2]
	i.Cookie = &http.Cookie{
		Name:   "session_for%3Asrun_cas_php",
		Value:  string(cookie),
		Domain: "ipgw.neu.edu.cn",
	}
}

// 向配置文件中写入用户配置
func (i *UserInfo) SaveAll(path string) {
	fmt.Println(tipSaveInfo)

	// 准备存储
	path, err := file.GetPath(path)
	if err != nil {
		fmt.Fprint(os.Stderr, errorGetPath)
		os.Exit(2)
	}

	// 打开
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	defer func() {
		_ = f.Close()
	}()
	w := bufio.NewWriter(f)

	if err != nil {
		fmt.Fprint(os.Stderr, errorOpenFile)
		os.Exit(2)
	}

	// 写入
	username := base64.StdEncoding.EncodeToString([]byte(i.Username))
	password := base64.StdEncoding.EncodeToString([]byte(i.Password))
	cookie := i.Cookie.Value

	// 如果保存账户
	_, err = w.WriteString(username + file.Delimiter + password + file.Delimiter + cookie)

	if err != nil {
		fmt.Fprint(os.Stderr, errorSaveInfo)
		os.Exit(2)
	}
	_ = w.Flush()

	// 输出成功提示
	fmt.Println(successSaveInfo)
}

func (i *UserInfo) SaveC(path string) {
	// 静默式，不需要输出
	path, err := file.GetPath(path)
	if err != nil {
		os.Exit(2)
	}

	// 读取
	bytes, _ := ioutil.ReadFile(path)
	content := string(bytes)

	// 分割
	split := strings.Split(content, file.Delimiter)

	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		_ = f
	}()
	w := bufio.NewWriter(f)

	_, _ = w.WriteString(split[0] + file.Delimiter + split[1] + file.Delimiter + i.Cookie.Value)
	_ = w.Flush()
}

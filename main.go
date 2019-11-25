package main

import (
	"flag"
	"fmt"
	"ipgw/lib"
)

var (
	userInfo = &lib.UserInfo{}

	h bool // help
	v bool // version
	s bool // save
	o bool // logout
	// TODO t bool // toggle (logout & login)

	u, p string // user, password
	i    string // configuration file path

	k string // kick out
)

// 初始化所有flag
func init() {
	flag.BoolVar(&h, "h", false, "show help info")
	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&s, "s", false, "save username and password after login successfully")
	flag.BoolVar(&o, "o", false, "log out")

	flag.StringVar(&u, "u", "", "`username`")
	flag.StringVar(&p, "p", "", "`password`")
	flag.StringVar(&i, "i", "", "`path` to configuration file, default is %USER_PROFILE%/.ipgw")
	flag.StringVar(&k, "k", "", "kick out the specific device")

	flag.Usage = usage
}

// 解析flag并转发到interactive中处理
func main() {
	flag.Parse()
	// -h
	if h {
		flag.Usage()
		return
	}

	// -v
	if v {
		version()
		return
	}

	// 处理主体逻辑
	handler()

	lib.Pause()
}

// 提供帮助命令输出 -h
func usage() {
	fmt.Println(lib.Title + "\n" +
		lib.Version + "\n" +
		"Usage: ipgw [-s] [-u username] [-p password] [-i configPath]\n" +
		"Options:")
	flag.PrintDefaults()
}

// 提供版本命令输出 -v
func version() {
	fmt.Println(lib.Title + "\n" +
		lib.Version)
}

func handler() {
	var err error

	// 若指定了k，则直接kick out
	if k != "" {
		err = lib.IKickOut(k)
		return
	}

	// 有指定u，则读取u p
	// 无指定u，则使用i读取配置

	// 若有o，则将登录操作改为登出
	// TODO 若无o，若有t，则先登陆，选择下线，再登陆
	// 需要改变一下项目结构

	if u != "" {
		// 使用-u -p 指定用户名与密码
		if o {
			err = lib.ILogoutUP(userInfo, u, p)
		} else {
			err = lib.ILoginUP(userInfo, u, p)

		}
	} else {
		// 默认使用配置文件，-i 指定配置文件路径
		if o {
			err = lib.ILogoutINI(userInfo, i)
		} else {
			err = lib.ILoginINI(userInfo, i)
		}
	}

	if err == nil {
		// 若激活了-s ，则写入到配置文件
		if s {
			err := lib.SaveBaseInfo(userInfo, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

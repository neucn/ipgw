package main

import (
	"flag"
	"fmt"
	"ipgw/lib"
)

var (
	userInfo = &lib.UserInfo{}

	h bool
	v bool
	s bool

	u, p string
	i    string
)

func init() {
	flag.BoolVar(&h, "h", false, "show help info")

	flag.BoolVar(&v, "v", false, "show version and exit")
	flag.BoolVar(&s, "s", false, "save username and password after login successfully")

	flag.StringVar(&u, "u", "", "`username`")
	flag.StringVar(&p, "p", "", "`password`")
	flag.StringVar(&i, "i", "", "`path` to configuration file, default is %USER_PROFILE%/.ipgw")

	// 改变默认的 Usage，flag包中的Usage 其实是一个函数类型。这里是覆盖默认函数实现，具体见后面Usage部分的分析
	flag.Usage = usage
}
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

	var err error

	if u != "" {
		// 使用-u -p 指定用户名与密码
		err = lib.LoginUP(userInfo, u, p)
	} else {
		// 默认使用配置文件，-i 指定配置文件路径
		err = lib.LoginINI(userInfo, i)
	}

	if err == nil {
		// 输出登陆信息
		userInfo.Print()

		// 若激活了-s ，则写入到配置文件
		if s {
			err := lib.SaveBaseInfo(userInfo, i)
			if err != nil {
				fmt.Println(err)
			}
		}
	}

	lib.Pause()
}

func usage() {
	fmt.Println("IPGW Tool\n" +
		"version: 0.0.1\n" +
		"Usage: ipgw [-s] [-u username] [-p password] [-i configPath]\n" +
		"Options:")
	flag.PrintDefaults()
}

func version() {
	fmt.Println("IPGW Tool\n" +
		"version: 0.0.1")
}

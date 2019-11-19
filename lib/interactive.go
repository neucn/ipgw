package lib

import (
	"errors"
	"fmt"
	"runtime"
)

func LoginUP(userInfo *UserInfo, username, password string) error {
	userInfo.Username = username
	userInfo.Password = password

	return login(userInfo)
}

func LoginINI(userInfo *UserInfo, path string) error {
	err := LoadBaseInfo(userInfo, path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return login(userInfo)
}

func login(userInfo *UserInfo) error {
	fmt.Println(userInfo.Username + ", 开始登陆...")
	err := Login(userInfo)

	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorLoginFail)
	}

	err = GetAccountInfo(userInfo)
	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorFetchInfoFail)
	}

	fmt.Println(userInfo.Username + ", 登陆成功！")
	return nil
}

func Pause() {
	if runtime.GOOS == "windows" {
		fmt.Println("请输入回车退出...")
		_, _ = fmt.Scanln()
	}
}

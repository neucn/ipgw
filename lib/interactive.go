// 提供了命令行交互到网关操作的转发与封装

package lib

import (
	"errors"
	"fmt"
	"runtime"
)

// 指定用户名、密码登陆
func ILoginUP(userInfo *UserInfo, username, password string) error {
	userInfo.Username = username
	userInfo.Password = password

	return login(userInfo)
}

// 使用配置文件登陆
func ILoginINI(userInfo *UserInfo, path string) error {
	err := LoadBaseInfo(userInfo, path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return login(userInfo)
}

// 指定用户名、密码登出
func ILogoutUP(userInfo *UserInfo, username, password string) error {
	userInfo.Username = username
	userInfo.Password = password

	return logout(userInfo)
}

// 使用配置文件登出
func ILogoutINI(userInfo *UserInfo, path string) error {
	err := LoadBaseInfo(userInfo, path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return logout(userInfo)
}

func IKickOut(sid string) error {
	err := KickOut(sid)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf(InfoKickOutSuccess, sid)
	return nil
}

// 包装Login，并加入GetAccountInfo
func login(userInfo *UserInfo) error {
	fmt.Printf(InfoBeginLogin, userInfo.Username)
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

	fmt.Printf(InfoLoginSuccess, userInfo.Username)

	// 输出用户信息
	userInfo.Print()
	return nil
}

// 包装Logout
func logout(userInfo *UserInfo) error {
	fmt.Printf(InfoBeginLogout, userInfo.Username)
	// 先登陆
	err := Login(userInfo)
	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorLoginFail)
	}

	// 再登出
	err = Logout(userInfo)

	if err != nil {
		fmt.Println(err)
		return errors.New(ErrorLogoutFail)
	}

	fmt.Printf(InfoLogoutSuccess, userInfo.Username)
	return nil
}

// 提供pause
func Pause() {
	if runtime.GOOS == "windows" {
		fmt.Println(InfoPause)
		_, _ = fmt.Scanln()
	}
}

package main

import (
	"fmt"
	"ipgw/lib"
)

func main() {
	userInfo := &lib.UserInfo{}
	err := lib.LoadBaseInfo(userInfo, "")
	if err != nil {
		fmt.Println("%+v", err)
		return
	}
	err = lib.Login(userInfo)

	if err != nil {
		fmt.Println("%+v", err)
		return
	}
	fmt.Println("%v", userInfo)
}

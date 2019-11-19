package lib

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type UserInfo struct {
	Username string
	Password string
	IP       string
	SID      string
	Balance  float64
	Used     int
	Time     int
}

func (i *UserInfo) Print() {
	fmt.Println(
		"===========Info===========" + "\n" +
			"Username: \t" + i.Username + "\n" +
			"IP: \t\t" + i.IP + "\n" +
			"SID: \t\t" + i.SID + "\n" +
			"Balance: \t" + fmt.Sprintf("%.2f", i.Balance) + " RMB\n" +
			"Used: \t\t" + getUsedFlux(i.Used) + "\n" +
			"Time: \t\t" + getUsedTime(i.Time) + "\n" +
			"==========================")
}

func LoadBaseInfo(info *UserInfo, path string) error {
	fmt.Println("读取用户配置中...")
	// 准备读取
	path, err := getPath(path)
	if err != nil {
		return errors.New("在获取环境时出错")
	}

	// 读取
	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return errors.New("配置文件读取失败")
	}
	bytes, _ := ioutil.ReadAll(f)
	content := string(bytes)

	// 分割
	split := strings.Split(content, ";")
	if len(split) != 2 {
		return errors.New("配置文件格式不正确")
	}

	username, err := base64.StdEncoding.DecodeString(split[0])
	info.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(split[1])
	info.Password = string(password)

	return nil
}

func SaveBaseInfo(info *UserInfo, path string) error {
	fmt.Println("存储用户配置中...")
	// 准备存储
	path, err := getPath(path)
	if err != nil {
		return errors.New("在获取环境时出错")
	}

	// 打开
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		_ = f.Close()
	}()
	w := bufio.NewWriter(f)

	if err != nil {
		return errors.New("配置文件打开失败")
	}

	username := base64.StdEncoding.EncodeToString([]byte(info.Username))
	password := base64.StdEncoding.EncodeToString([]byte(info.Password))

	_, err = w.WriteString(username + ";" + password)
	if err != nil {
		return errors.New("在写入配置时出错")
	}
	_ = w.Flush()

	fmt.Println("存储用户配置成功")
	return nil
}

func getPath(path string) (string, error) {
	if path == "" {
		homeDir, err := Home()
		if err != nil {
			return "", err
		}
		path = homeDir + string(os.PathSeparator) + ".ipgw"
	}

	MustExist(path)
	return path, nil
}

func getUsedFlux(flux int) string {
	if flux > 1000*1000 {
		return fmt.Sprintf("%.2f M", float64(flux)/(1000*1000))
	}
	if flux > 1000 {
		return fmt.Sprintf("%.2f K", float64(flux)/1000)
	}
	return fmt.Sprintf("%.2f b", flux)
}

func getUsedTime(time int) string {
	h := time / 3600
	m := (time % 3600) / 60
	s := time % 3600 % 60

	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

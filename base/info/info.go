package info

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"ipgw/base/file"
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

// 格式化输出用户的所有信息
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

// 解析已用流量数
func getUsedFlux(flux int) string {
	if flux > 1000*1000 {
		return fmt.Sprintf("%.2f M", float64(flux)/(1000*1000))
	}
	if flux > 1000 {
		return fmt.Sprintf("%.2f K", float64(flux)/1000)
	}
	return fmt.Sprintf("%d b", flux)
}

// 解析已使用时长
func getUsedTime(time int) string {
	h := time / 3600
	m := (time % 3600) / 60
	s := time % 3600 % 60

	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

// 从配置文件中解析出用户配置
func LoadBaseInfo(info *UserInfo, path string) error {
	fmt.Println(InfoLoadUserInfo)
	// 准备读取
	path, err := getPath(path)
	if err != nil {
		return errors.New(ErrorGetPath)
	}

	// 读取
	f, err := os.Open(path)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		return errors.New(ErrorLoadUserInfo)
	}
	bytes, _ := ioutil.ReadAll(f)
	content := string(bytes)

	// 分割
	split := strings.Split(content, file.Delimiter)
	if len(split) != 2 {
		return errors.New(ErrorUserInfoFormat)
	}

	username, err := base64.StdEncoding.DecodeString(split[0])
	info.Username = string(username)

	password, err := base64.StdEncoding.DecodeString(split[1])
	info.Password = string(password)

	return nil
}

// 向配置文件中写入用户配置
func SaveBaseInfo(info *UserInfo, path string) error {
	fmt.Println(InfoSaveUserInfo)
	// 准备存储
	path, err := getPath(path)
	if err != nil {
		return errors.New(ErrorGetPath)
	}

	// 打开
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		_ = f.Close()
	}()
	w := bufio.NewWriter(f)

	if err != nil {
		return errors.New(ErrorOpenUserInfo)
	}

	// 写入
	username := base64.StdEncoding.EncodeToString([]byte(info.Username))
	password := base64.StdEncoding.EncodeToString([]byte(info.Password))

	_, err = w.WriteString(username + file.Delimiter + password)
	if err != nil {
		return errors.New(ErrorSaveUserInfo)
	}
	_ = w.Flush()

	// 输出成功提示
	fmt.Println(InfoSaveUserInfoSuccess)
	return nil
}

// 获取配置文件路径
func getPath(path string) (string, error) {
	// 若路径为空则使用默认路径
	if path == "" {
		homeDir, err := file.Home()
		if err != nil {
			return "", err
		}
		path = homeDir + string(os.PathSeparator) + ".ipgw"
	}

	// 确保路径存在
	file.MustExist(path)
	return path, nil
}

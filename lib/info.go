package lib

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type UserInfo struct {
	Username string
	Password string
	IP       string
	SID      string
	Balance  float32
	Used     float32
}

func LoadBaseInfo(info *UserInfo, path string) error {
	// 准备读取
	path, err := getPath(path)
	if err != nil {
		return errors.New("在获取环境时出错")
	}

	// 读取
	f, err := os.Open(path)
	defer func() {
		f.Close()
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
	// 准备读取
	path, err := getPath(path)
	if err != nil {
		return errors.New("在获取环境时出错")
	}

	// 读取
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	defer func() {
		f.Close()
	}()
	w := bufio.NewWriter(f)

	if err != nil {
		return errors.New("配置文件读取失败")
	}

	username := base64.StdEncoding.EncodeToString([]byte(info.Username))
	password := base64.StdEncoding.EncodeToString([]byte(info.Password))

	_, err = w.WriteString(username + ";" + password)
	if err != nil {
		return errors.New("在写入配置时出错")
	}
	_ = w.Flush()
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

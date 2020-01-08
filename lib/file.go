// 存放文件操作的相关函数

package lib

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

var (
	LineDelimiter = "\n"
	PartDelimiter = ":"
)

// 获取配置文件路径
func GetPath(path string) (string, error) {
	// 若路径为空则使用默认路径
	if path == "" {
		homeDir, err := home()
		if err != nil {
			return "", err
		}
		path = homeDir + string(os.PathSeparator) + ".ipgw"
	}

	// 确保路径存在
	mustExist(path)
	return path, nil
}

// 获取用户目录路径
func home() (string, error) {
	usr, err := user.Current()
	if nil == err {
		return usr.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

// 确保路径存在
func mustExist(path string) {
	file, err := os.Open(path)
	defer func() { _ = file.Close() }()
	if err != nil && os.IsNotExist(err) {
		file, _ = os.Create(path)
	}
}

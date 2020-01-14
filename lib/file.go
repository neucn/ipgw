// 存放文件操作的相关函数

package lib

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	. "ipgw/base"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	LineDelimiter = "\n"
	PartDelimiter = ":"
)

// 获取配置文件路径
func GetConfigPath(path string) (string, error) {
	// 若路径为空则使用默认路径
	if path == "" {
		homeDir, err := home()
		if err != nil {
			return "", err
		}
		path = homeDir + string(os.PathSeparator) + ".ipgw"
	}

	// 确保路径存在
	FileMustExist(path)
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
func FileMustExist(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_, _ = os.Create(path)
	}
}

// 确保路径存在
func DirMustExist(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		_ = os.Mkdir(path, os.ModePerm)
	}

}

// 获取ipgw实际路径与所在的目录(末尾有`/`)
func GetRealPathAndDir() (path, dir string) {
	// 获取到当前执行路径
	p, err := os.Executable()
	if err != nil {
		FatalF(errEnvReason, err)
	}
	// 当前运行的版本的路径
	path, _ = filepath.Abs(p)
	// 当前运行的版本的所在目录
	dir = filepath.Dir(path) + string(os.PathSeparator)
	return
}

func Unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		fpath := filepath.Join(destDir, f.Name)
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			inFile, err := f.Open()
			if err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, inFile)

			// 关闭资源
			inFile.Close()
			_ = outFile.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取 ipgwTools/tools.json 所在路径
func GetToolsConfigPath() (path string) {
	// 不使用filepath.Join
	// 确保文件夹和文件都存在
	path = GetToolsDir()

	path += "tools.json"
	FileMustExist(path)
	return
}

// 获取 ipgwTools 所在路径，末尾有`/`
func GetToolsDir() (path string) {
	// 获取程序所在目录
	_, dir := GetRealPathAndDir()
	// 不使用filepath.Join
	path = dir + "ipgwTool" + string(os.PathSeparator)
	// 路径不存在则新建
	DirMustExist(path)
	return
}

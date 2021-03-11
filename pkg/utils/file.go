package utils

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if nil == err {
		return usr.HomeDir, nil
	}
	if "windows" == runtime.GOOS {
		return homeWindows()
	}
	return homeUnix()
}

func homeUnix() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

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

func FileMustExist(path string) (err error) {
	if !IsFileExist(path) {
		_, err = os.Create(path)
	}
	return
}

func IsFileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func GetExecutablePathAndDir() (path, dir string, err error) {
	p, err := os.Executable()
	if err != nil {
		return "", "", fmt.Errorf("fail to get executable path: %v", err)
	}
	path, _ = filepath.Abs(p)
	dir = filepath.Dir(path)
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

			inFile.Close()
			_ = outFile.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}

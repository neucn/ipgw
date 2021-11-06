package handler

import (
	"errors"
	"fmt"
	"github.com/neucn/ipgw"
	"github.com/neucn/ipgw/pkg/console"
	"github.com/neucn/ipgw/pkg/utils"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"
)

type downloader struct {
	io.Reader
	total   int64
	current int64
}

type UpdateHandler struct {
	client *http.Client
}

func (d *downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)
	d.current += int64(n)
	console.InfoF("\rdownloading %.2f%%", float64(d.current*10000/d.total)/100)
	return
}

func NewUpdateHandler() *UpdateHandler {
	return &UpdateHandler{
		client: &http.Client{Timeout: 90 * time.Second},
	}
}

// CheckLatestVersion returns true when there is a newer version
func (u *UpdateHandler) CheckLatestVersion() (bool, error) {
	resp, err := u.client.Get(fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", ipgw.Repo))
	if err != nil {
		return false, fmt.Errorf("fail to check latest version:\n\t%v", err)
	}
	body := utils.ReadBody(resp)
	latestVersion, _ := utils.MatchSingle(regexp.MustCompile(`"tag_name": *"(.+?)"`), body)
	return utils.CompareVersion(utils.ParseVersion(latestVersion), utils.ParseVersion(ipgw.Version)), nil
}

// download returns downloaded path
func (u *UpdateHandler) download(url string) (string, error) {
	resp, err := u.client.Get(url)
	if err != nil {
		return "", err
	}
	// not found
	if resp.StatusCode == http.StatusNotFound {
		return "", errors.New("release not found")
	}
	raw := resp.Body
	defer raw.Close()

	// create a temp file
	f, err := os.CreateTemp(os.TempDir(), "ipgw.release.*")
	if err != nil {
		return "", err
	}

	_, err = io.Copy(f, &downloader{
		Reader: raw,
		total:  resp.ContentLength,
	})
	console.InfoL()
	if err != nil {
		return "", err
	}

	_ = f.Close()
	return f.Name(), nil
}

func (u *UpdateHandler) Update() error {
	url := fmt.Sprintf("https://github.com/%s/releases/latest/download/ipgw-%s-%s.zip", ipgw.Repo, runtime.GOOS, runtime.GOARCH)
	downloaded, err := u.download(url)
	if err != nil {
		return err
	}
	_ = os.Chmod(downloaded, 0777)

	// get executable path
	path, dir, err := utils.GetExecutablePathAndDir()
	if err != nil {
		return err
	}
	// rename current running executable
	backupPath := filepath.Join(dir, "ipgw."+ipgw.Version)
	if err = os.Rename(path, backupPath); err != nil {
		return err
	}

	if err = utils.Unzip(downloaded, dir); err != nil {
		_ = os.Rename(backupPath, path)
		return err
	}

	return nil
}

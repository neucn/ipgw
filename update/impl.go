package update

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"ipgw/base/cfg"
	"ipgw/base/ctx"
	"ipgw/share"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"time"
)

func checkVersion() *ver {
	fmt.Println(querying)
	i := &ver{Update: false}
	client := ctx.GetClient()
	// 优先neu.ee
	resp, err := client.Get(cfg.ReleasePath + "/info.json")
	if err == nil {
		res, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if cfg.FullView {
			fmt.Printf(getResponse, string(res))
		}

		_ = json.Unmarshal(res, &i)

		if len(i.Latest) < 1 {
			fmt.Fprintln(os.Stderr, failQuery)
			os.Exit(2)
		}

		if !isNewer(i.Latest) {
			fmt.Println(alreadyLatest)
			return i
		}
		i.Update = true
		fmt.Printf(latestVersion, i.Latest)
		return i
	} else {
		if cfg.FullView {
			fmt.Fprintln(os.Stderr, useGithub)
		}
	}

	// neu.ee超时之后使用github。

	resp, err = client.Get("https://api.github.com/repos/imyown/ipgw/releases/latest")
	if err != nil {
		fmt.Fprintln(os.Stderr, errNet)
		os.Exit(2)
	}

	body := share.ReadBody(resp)

	tagExp := regexp.MustCompile(`"tag_name":"(.+?)"`)
	tags := tagExp.FindAllStringSubmatch(body, -1)

	if len(tags) == 0 {
		fmt.Fprintln(os.Stderr, failQuery)
		os.Exit(2)
	}

	tag := tags[0][1]
	if !isNewer(tag) {
		fmt.Println(alreadyLatest)
		return i
	}

	i.Update = true
	fmt.Printf(latestVersion, tag)
	return i
}

func printChangelog(i *ver) {
	fmt.Println()
	fmt.Println(changelog)
	// 版本排序
	var tmp = make([]string, 0)
	for k, _ := range i.Changelog {
		tmp = append(tmp, k)
	}

	sort.Sort(sort.Reverse(sort.StringSlice(tmp)))
	for _, v := range tmp {
		if v == cfg.Version {
			break
		}
		fmt.Printf(changelogTitle, v)
		for _, l := range i.Changelog[v] {
			fmt.Printf(changelogContent, l)
		}
	}
	fmt.Println()
}

func getReleaseUrl(i *ver) string {
	u := cfg.ReleasePath
	// 检查arch
	_, ok := i.Arch[runtime.GOARCH]
	if !ok {
		fmt.Fprintln(os.Stderr, failArchNotSupported)
		os.Exit(2)
	}

	// 检查版本
	if len(i.Latest) < 1 {
		fmt.Fprintln(os.Stderr, failGetReleasePath)
		os.Exit(2)
	}
	u += "/" + i.Latest

	// 检查os
	s, ok := i.OS[runtime.GOOS]
	if !ok {
		fmt.Fprintln(os.Stderr, failOSNotSupported)
		os.Exit(2)
	}
	u += "/" + s

	// 获取文件名
	n, ok := i.Name[runtime.GOOS]
	if !ok {
		fmt.Fprintln(os.Stderr, failOSNotSupported)
		os.Exit(2)
	}
	u += "/" + n
	return u
}

func download(u string, dir string) {
	client := &http.Client{}
	client.Timeout = 60 * time.Second
	resp, err := client.Get(u)
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errReason, err)
		}
		fmt.Fprintln(os.Stderr, errNet)
		os.Exit(2)
	}
	if resp.StatusCode == http.StatusNotFound {
		fmt.Fprintln(os.Stderr, wrongUrl)
		os.Exit(2)
	}

	raw := resp.Body
	defer raw.Close()

	f, err := os.Create(dir + "ipgw.download")
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errReason, err)
		}
		fmt.Fprintln(os.Stderr, failCreate)
		os.Exit(2)
	}

	d := &downloader{
		Reader: raw,
		Total:  resp.ContentLength,
	}
	io.Copy(f, d)
	f.Chmod(0777)
	f.Close()
	fmt.Println()
}

func update(i *ver) {
	path, err := os.Executable()
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errReason, err)
		}
		fmt.Fprintln(os.Stderr, errRunEnv)
		os.Exit(2)
	}
	old, _ := filepath.Abs(path)
	dir := filepath.Dir(old) + string(os.PathSeparator)

	// 下载
	download(getReleaseUrl(i), dir)

	fmt.Println(updating)

	if cfg.FullView {
		fmt.Println(removing)
	}
	err = os.Rename(old, dir+"ipgw.old")
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errReason, err)
		}
		fmt.Println(failUpdate)
		os.Exit(2)
	}

	if cfg.FullView {
		fmt.Println(covering)
	}

	err = os.Rename(dir+"ipgw.download", dir+i.Name[runtime.GOOS])
	if err != nil {
		if cfg.FullView {
			fmt.Fprintf(os.Stderr, errReason, err)
		}
		fmt.Println(failUpdate)
		os.Exit(2)
	}
	fmt.Println(successUpdate)
}

func isNewer(v string) bool {
	var fetched, local []string
	if len(cfg.Version) < 1 {
		return true
	}
	fetched = strings.Split(v[1:], ".")
	local = strings.Split(cfg.Version[1:], ".")
	for i := 0; i < 3; i++ {
		if fetched[i] > local[i] {
			return true
		}
	}
	return false
}

package lib

import (
	"io"
	"net/http"
	"os"
	"runtime"
	"sort"
	"time"
)

type Ver struct {
	Update    bool
	Latest    string              `json:"latest"`
	API       string              `json:"api"`
	Changelog map[string][]string `json:"changelog"`
	OS        map[string]string   `json:"os"`
	Arch      map[string]string   `json:"arch"`
	Name      map[string]string   `json:"name"`
}

type downloader struct {
	io.Reader
	Total   int64
	Current int64
}

func (d *downloader) Read(p []byte) (n int, err error) {
	n, err = d.Reader.Read(p)

	d.Current += int64(n)
	InfoF("\r下载进度 %.2f%%", float64(d.Current*10000/d.Total)/100)

	return
}

func Download(u string, dir, tmpFile string) {
	// 获取client, 不适用ctx中的client
	client := &http.Client{}
	client.Timeout = 60 * time.Second

	// 发送请求, 不适用global中的SendRequest
	resp, err := client.Get(u)
	if err != nil {
		FatalF(failConnect, err)
	}

	// 若404
	if resp.StatusCode == http.StatusNotFound {
		Fatal(wrongUrl)
	}

	raw := resp.Body
	defer raw.Close()

	// 新建临时文件
	f, err := os.Create(dir + tmpFile)
	if err != nil {
		FatalF(failCreate, err)
	}

	// 实例化下载工具对象
	d := &downloader{
		Reader: raw,
		Total:  resp.ContentLength,
	}

	// 简单地使用Copy函数，downloader对象中增加了hook因此能获取到copy的进度
	_, err = io.Copy(f, d)
	if err != nil {
		// 若下载失败
		Fatal(failDownload)
	}

	// 修改权限, 这个函数里有下载文件的file对象，修改比较方便，因此就不分离出去了
	err = f.Chmod(0777)
	if err != nil {
		FatalF(failChmod, err)
	}
	_ = f.Close()
	InfoLine()
}

// 拼接获取下载路径
// u为release的目录，如neu.ee/release
func GetDownloadUrl(u string, i *Ver) string {
	// 检查本机arch是否存在于预编译版本arch中
	_, ok := i.Arch[runtime.GOARCH]
	if !ok {
		// 若不存在
		Fatal(failArchNotSupported)
	}

	// 版本号接入Release地址
	u += "/" + i.Latest

	// 检查os是否存在于预编译版本os中
	s, ok := i.OS[runtime.GOOS]
	if !ok {
		// 若不存在
		Fatal(failOSNotSupported)
	}

	// os接入Release地址
	u += "/" + s

	// 获取对应os下的预编译文件名
	n, ok := i.Name[runtime.GOOS]
	if !ok {
		// 若不存在
		Fatal(failOSNotSupported)
	}

	// 文件名接入Release地址
	u += "/" + n

	// 拼接完成
	return u
}

// 输出更新日志，v为当前版本
func PrintChangelog(local string, i *Ver) {
	InfoLine()
	InfoLine(changelog)
	// 版本排序，由于json反序列化之后map无序，因此需要排序
	var tmp = make([]string, 0)
	for k, _ := range i.Changelog {
		tmp = append(tmp, k)
	}

	// 使得版本号从大到小
	sort.Sort(sort.Reverse(sort.StringSlice(tmp)))

	// 输出更新日志
	for _, v := range tmp {
		if v == local {
			// 若遍历到当前版本，则跳出循环
			break
		}
		InfoF(changelogTitle, v)
		for _, l := range i.Changelog[v] {
			InfoF(changelogContent, l)
		}
	}
	InfoLine()
}

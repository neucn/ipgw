package lib

import (
	"encoding/json"
	"io"
	"io/ioutil"
	. "ipgw/base"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 传入Release或某工具的根路径（不含最后一个`/`)，解析其中的`info.json`为Ver对象
func ParseVer(u string) (v *Ver) {
	client := &http.Client{Timeout: 3 * time.Second}
	// 使用/info.json获取版本信息
	resp, err := client.Get(u + "/info.json")
	if err != nil {
		FatalL(errNetwork)
	}
	// 因为json.Unmarshal，所以这里不适用ReadBody
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	v = &Ver{Update: false}
	_ = json.Unmarshal(res, &v)

	// 若解析失败
	if len(v.Latest) < 1 {
		FatalL(failQuery)
	}
	return
}

// 下载u到dest指定的路径
func Download(u, dest string) {
	// 获取client, 不适用ctx中的client
	client := &http.Client{Timeout: 60 * time.Second}

	// 发送请求, 不适用global中的SendRequest
	resp, err := client.Get(u)
	if err != nil {
		FatalF(failConnect, err)
	}

	// 若404
	if resp.StatusCode == http.StatusNotFound {
		FatalL(wrongUrl)
	}

	raw := resp.Body
	defer raw.Close()

	// 新建临时文件
	f, err := os.Create(dest)
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
	// 换行
	InfoL()
	if err != nil {
		// 若下载失败
		FatalL(failDownload)
	}

	_ = f.Close()
}

// 拼接获取下载路径
// u为release的目录，如neu.ee/release
func GetDownloadUrl(u string, i *Ver) string {
	// 检查本机arch是否存在于预编译版本arch中
	if _, ok := i.Arch[runtime.GOARCH]; !ok {
		// 若不存在
		FatalL(failArchNotSupported)
	}

	// 版本号接入Release地址
	u += "/" + i.Latest

	// 检查os是否存在于预编译版本os中
	s, ok := i.OS[runtime.GOOS]
	if !ok {
		// 若不存在
		FatalL(failOSNotSupported)
	}

	// os接入Release地址
	u += "/" + s

	// 获取对应os下的预编译文件名
	n, ok := i.Name[runtime.GOOS]
	if !ok {
		// 若不存在
		FatalL(failOSNotSupported)
	}

	// 文件名接入Release地址
	u += "/" + n

	// 拼接完成
	return u
}

// 拼接获取工具下载路径
// u为tool的根目录，如neu.ee/tools/teemo
func GetToolDownloadUrl(u, toolName string, i *Ver) string {
	// 检查本机arch是否存在于预编译版本arch中
	if _, ok := i.Arch[runtime.GOARCH]; !ok {
		// 若不存在
		FatalL(failArchNotSupported)
	}

	// 版本号接入Release地址
	u += "/" + i.Latest

	// 检查os是否存在于预编译版本os中
	s, ok := i.OS[runtime.GOOS]
	if !ok {
		// 若不存在
		FatalL(failOSNotSupported)
	}

	// os接入Release地址
	u += "/" + s

	// 文件名接入Release地址
	u += "/" + toolName + ".zip"

	// 拼接完成
	return u
}

// 输出更新日志，v为当前版本
func PrintChangelog(local string, i *Ver) {
	InfoL()
	InfoL(changelog)
	// 版本排序，由于json反序列化之后map无序，因此需要排序
	var tmp = make([]string, 0)
	for k := range i.Changelog {
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
	InfoL()
}

// 简单判断是否兼容
func IsAPICompatible(v string) bool {
	return API == v
}

// 获取Tools列表
func GetOnlineToolList() *Tools {
	// 实例化一个客户端
	client := &http.Client{Timeout: 3 * time.Second}
	// 请求工具列表
	resp, err := client.Get(ToolReleasePath + "/tools.json")
	if err != nil {
		FatalL(errNetwork)
	}
	// 因为需要用到[]byte所以不ReadBody
	bytes, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	// 实例化Tools
	toolList := &Tools{}
	// 解析
	toolList.Parse(bytes)

	return toolList
}

// 判断a是否比b版本新
func IsNewer(a, b string) bool {
	if len(b) < 1 {
		return true
	}
	aList := strings.Split(a[1:], ".")
	bList := strings.Split(b[1:], ".")
	var tmpA, tmpB int
	for i := 0; i < 3; i++ {
		tmpA, _ = strconv.Atoi(aList[i])
		tmpB, _ = strconv.Atoi(bList[i])
		if tmpA > tmpB {
			return true
		}
	}
	return false
}

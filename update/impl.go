package update

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"ipgw/base"
	"ipgw/ctx"
	. "ipgw/lib"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"
)

// 获取版本信息
func checkVersion() (i *ver) {
	// Todo 暂时删去对Github Release API的使用。原因如下：
	//  1. 速度有待商榷；
	//  2.就算通过Github获取了新版本信息，由于没有对应的下载实现，等于无用

	InfoLine(querying)

	// 实例化对象
	i = &ver{Update: false}
	client := ctx.NewClient()

	// 优先使用neu.ee/release/info.json获取版本信息
	resp, err := client.Get(base.ReleasePath + "/info.json")
	if err != nil {
		Fatal(errNet)
	}
	// 因为json.Unmarshal，所以这里不适用ReadBody
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(res, &i)

	// 若解析失败
	if len(i.Latest) < 1 {
		Fatal(failQuery)
	}

	// 判断是否是更新版本
	if !isNewer(i.Latest) {
		// 已是最新
		InfoLine(alreadyLatest)
		return i
	}
	// 不是最新
	i.Update = true
	InfoF(latestVersion, i.Latest)
	return i
}

func printChangelog(i *ver) {
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
		if v == base.Version {
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

func getReleaseUrl(i *ver) string {
	u := base.ReleasePath
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

func download(u string, dir string) {
	// 获取client, 不适用ctx中的client
	client := &http.Client{}
	client.Timeout = 60 * time.Second

	// 发送请求, 不适用global中的SendRequest
	resp, err := client.Get(u)
	fatalHandler(err, errNet)

	// 若404
	if resp.StatusCode == http.StatusNotFound {
		Fatal(wrongUrl)
	}

	raw := resp.Body
	defer raw.Close()

	// 新建临时文件
	f, err := os.Create(dir + "ipgw.download")
	fatalHandler(err, failCreate)

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
		Fatal(failChmod)
	}
	_ = f.Close()
	InfoLine()
}

func update(i *ver) {
	// 获取到当前执行路径
	path, err := os.Executable()
	if err != nil {
		if ctx.FullView {
			ErrorF(errReason, err)
		}
		Fatal(errRunEnv)
	}
	// 当前运行的版本的路径
	old, _ := filepath.Abs(path)
	// 当前运行的版本的所在目录
	dir := filepath.Dir(old) + string(os.PathSeparator)

	// 下载
	download(getReleaseUrl(i), dir)

	InfoLine(updating)

	if ctx.FullView {
		InfoLine(removing)
	}

	// 将当前运行的版本改名为ipgw.old
	err = os.Rename(old, dir+"ipgw.old")
	fatalHandler(err, failUpdate)

	if ctx.FullView {
		InfoLine(covering)
	}

	// 将下载的新版本改名为正确名字
	err = os.Rename(dir+"ipgw.download", dir+i.Name[runtime.GOOS])
	fatalHandler(err, failUpdate)

	InfoLine(successUpdate)
}

// 判断一个版本是否比当前版本新
func isNewer(v string) bool {
	var fetched, local []string
	if len(base.Version) < 1 {
		return true
	}
	fetched = strings.Split(v[1:], ".")
	local = strings.Split(base.Version[1:], ".")
	for i := 0; i < 3; i++ {
		if fetched[i] > local[i] {
			return true
		}
	}
	return false
}

// 处理err
func fatalHandler(err error, fatalText string) {
	if err != nil {
		if ctx.FullView {
			ErrorF(errReason, err)
		}
		Fatal(fatalText)
	}
}

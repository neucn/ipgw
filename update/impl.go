package update

import (
	"encoding/json"
	"io/ioutil"
	"ipgw/base"
	"ipgw/ctx"
	. "ipgw/lib"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// 获取版本信息
func checkVersion() (v *Ver) {
	// Todo 暂时删去对Github Release API的使用。原因如下：
	//  1. 速度有待商榷；
	//  2.就算通过Github获取了新版本信息，由于没有对应的下载实现，等于无用

	InfoLine(querying)

	// 实例化对象
	v = &Ver{Update: false}
	client := ctx.NewClient()

	// 优先使用neu.ee/release/info.json获取版本信息
	resp, err := client.Get(base.ReleasePath + "/info.json")
	if err != nil {
		Fatal(errNet)
	}
	// 因为json.Unmarshal，所以这里不适用ReadBody
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	_ = json.Unmarshal(res, &v)

	// 若解析失败
	if len(v.Latest) < 1 {
		Fatal(failQuery)
	}

	// 判断是否是更新版本
	if !isNewer(v.Latest) {
		// 已是最新
		InfoLine(alreadyLatest)
		return v
	}
	// 不是最新
	v.Update = true
	InfoF(latestVersion, v.Latest)
	return v
}

func update(v *Ver) {
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
	Download(GetDownloadUrl(base.ReleasePath, v), dir, "ipgw.download")

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
	err = os.Rename(dir+"ipgw.download", dir+v.Name[runtime.GOOS])
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
	var tmpF, tmpL int
	for i := 0; i < 3; i++ {
		tmpF, _ = strconv.Atoi(fetched[i])
		tmpL, _ = strconv.Atoi(local[i])
		if tmpF > tmpL {
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

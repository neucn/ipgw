package update

import (
	. "ipgw/base"
	"ipgw/ctx"
	. "ipgw/lib"
	"os"
	"runtime"
)

// 获取版本信息
func checkVersion() (v *Ver) {
	// Todo 暂时删去对Github Release API的使用。原因如下：
	//  1. 速度有待商榷；
	//  2.就算通过Github获取了新版本信息，由于没有对应的下载实现，等于无用

	InfoL(querying)

	// 查询信息解析为Ver
	v = ParseVer(ReleasePath)

	// 判断是否是更新版本
	if !IsNewer(v.Latest, Version) {
		// 本地已是最新
		InfoL(alreadyLatest)
		return v
	}
	// 本地不是最新
	v.Update = true
	InfoF(latestVersion, v.Latest)
	return v
}

func update(v *Ver) {
	// 当前运行的版本的所在目录
	old, dir := GetRealPathAndDir()

	// 临时文件名
	tmpName := dir + "ipgw.download"

	// 下载
	Download(GetDownloadUrl(ReleasePath, v), tmpName)

	InfoL(updating)

	// 修改下载的文件权限
	_ = os.Chmod(tmpName, 0777)

	if ctx.FullView {
		InfoL(removing)
	}

	// 将当前运行的版本改名为ipgw.old
	err := os.Rename(old, dir+"ipgw.old")
	fatalHandler(err, failUpdate)

	if ctx.FullView {
		InfoL(covering)
	}

	// 将下载的新版本改名为正确名字
	err = os.Rename(tmpName, dir+v.Name[runtime.GOOS])
	fatalHandler(err, failUpdate)

	InfoL(successUpdate)
}

// 处理err
func fatalHandler(err error, fatalText string) {
	if err != nil {
		if ctx.FullView {
			ErrorF(errReason, err)
		}
		FatalL(fatalText)
	}
}

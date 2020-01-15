package lib

import (
	"encoding/json"
	"io"
	"io/ioutil"
	. "ipgw/base"
	"os"
)

type Tool struct {
	// Version字段里的内容来自Ver.Latest，为了方便存入tools.json因此在这里冗余
	Version     string
	API         string `json:"api"`
	Name        string `json:"name"`
	Path        string `json:"path"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

type Tools struct {
	List map[string]*Tool
}

// 从配置文件里载入
func (i *Tools) Load() {
	// 准备读取
	path := GetToolsConfigPath()

	// 读取
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		FatalL(fatalLoadToolInfo)
	}

	// 反序列化
	i.Parse(bytes)
}

// 写入配置文件
func (i *Tools) Save() {
	// 准备写入
	path := GetToolsConfigPath()
	// 序列化
	bytes, err := json.Marshal(i.List)
	if err != nil {
		FatalL(fatalSaveToolInfo)
	}

	// 打开
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)

	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		FatalL(fatalSaveToolInfo)
	}

	// 写入
	if _, err = f.Write(bytes); err != nil {
		FatalL(fatalSaveToolInfo)
	}
}

// 反序列化
func (i *Tools) Parse(text []byte) {
	_ = json.Unmarshal(text, &i.List)
	if i.List == nil {
		i.List = map[string]*Tool{}
	}
}

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

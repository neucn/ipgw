package update

import (
	"io"
	. "ipgw/lib"
)

type ver struct {
	Update    bool
	Latest    string              `json:"latest"`
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

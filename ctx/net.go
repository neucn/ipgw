package ctx

import (
	"fmt"
	. "ipgw/base"
	"net/http"
)

type Net struct {
	IP      string
	SID     string
	Balance float64
	Used    int
	Time    int
	// 网关的Cookie
	Cookie *http.Cookie
}

func (i *Net) Print() {
	InfoF(`
# 信息
   IP	%16s
   SID	%16s
   余额	%16s
   流量	%16s
   时长	%16s
`, i.IP, i.SID, getBalance(i.Balance), getUsedFlux(i.Used), getUsedTime(i.Time))
}

func (i *Net) SetCookie(c string) {
	i.Cookie = &http.Cookie{
		Name:   "session_for%3Asrun_cas_php",
		Value:  c,
		Domain: "ipgw.neu.edu.cn",
	}
}

// 解析已用流量数
func getUsedFlux(flux int) string {
	if flux > 1000*1000 {
		return fmt.Sprintf("%.2f M", float64(flux)/(1000*1000))
	}
	if flux > 1000 {
		return fmt.Sprintf("%.2f K", float64(flux)/1000)
	}
	return fmt.Sprintf("%d b", flux)
}

// 解析已使用时长
func getUsedTime(time int) string {
	h := time / 3600
	m := (time % 3600) / 60
	s := time % 3600 % 60

	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

func getBalance(balance float64) string {
	return fmt.Sprintf("%.2f R", balance)
}

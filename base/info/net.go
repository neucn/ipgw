package info

import "fmt"

type NetInfo struct {
	IP      string
	SID     string
	Balance float64
	Used    int
	Time    int
}

func (i *NetInfo) Print() {
	fmt.Printf(
		`=========信息=========
IP	%14s
SID	%14s
余额	%14s
流量	%14s
时长	%14s
`, i.IP, i.SID, getBalance(i.Balance), getUsedFlux(i.Used), getUsedTime(i.Time))
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

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
	fmt.Print(
		"\n========信息========\n" +
			"IP: \t" + i.IP + "\n" +
			"SID: \t" + i.SID + "\n" +
			"余额: \t" + fmt.Sprintf("%.2f", i.Balance) + " 元\n" +
			"流量: \t" + getUsedFlux(i.Used) + "\n" +
			"时长: \t" + getUsedTime(i.Time) + "\n")
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

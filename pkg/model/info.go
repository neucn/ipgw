package model

import (
	"fmt"
)

type Info struct {
	Username          string
	IP                string
	Traffic, UsedTime int
	Balance           float64
}

func (i *Info) FormattedTraffic() string {
	if i.Traffic > 1000*1000 {
		return fmt.Sprintf("%.2f M", float64(i.Traffic)/(1000*1000))
	}
	if i.Traffic > 1000 {
		return fmt.Sprintf("%.2f K", float64(i.Traffic)/1000)
	}
	return fmt.Sprintf("%d b", i.Traffic)
}

func (i *Info) FormattedUsedTime() string {
	time := i.UsedTime
	h := time / 3600
	m := (time % 3600) / 60
	s := time % 3600 % 60

	return fmt.Sprintf("%d:%02d:%02d", h, m, s)
}

func (i *Info) FormattedBalance() string {
	return fmt.Sprintf("%.2f R", i.Balance)
}

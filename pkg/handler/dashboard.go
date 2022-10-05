package handler

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/neucn/ipgw/pkg/model"
	"github.com/neucn/ipgw/pkg/utils"
	"github.com/neucn/neugo"
)

type DashboardHandler struct {
	client                      *http.Client
	cachedDashboardIndexContent string
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{
		client: neugo.NewSession(),
	}
}

func (d *DashboardHandler) Login(account *model.Account) error {
	password, err := account.GetPassword()
	if err != nil {
		return err
	}
	err = neugo.Use(d.client).WithAuth(account.Username, password).Login(neugo.CAS)
	if err != nil {
		return err
	}
	_, err = d.client.Get("https://ipgw.neu.edu.cn:8800/sso/neusoft/index") // 统一认证获取cookie
	if err != nil {
		return err
	}
	return nil
}

func (d *DashboardHandler) fetchDashboardIndexBody() (string, error) {
	resp, err := d.client.Get("https://ipgw.neu.edu.cn:8800/home")
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) fetchDashboardBillsBody(page int) (string, error) {
	resp, err := d.client.Get(fmt.Sprintf("https://ipgw.neu.edu.cn:8800/log/check-out?page=%d&per-page=10", page))
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) fetchDashboardRechargeBody(page int) (string, error) {
	resp, err := d.client.Get(fmt.Sprintf("https://ipgw.neu.edu.cn:8800/log/pay?page=%d&per-page=10", page))
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) fetchDashboardUsageLogBody(page int) (string, error) {
	resp, err := d.client.Get(fmt.Sprintf("https://ipgw.neu.edu.cn:8800/log/detail?page=%d&per-page=10", page))
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) getCachedDashboardIndexBody() (string, error) {
	if d.cachedDashboardIndexContent == "" {
		body, err := d.fetchDashboardIndexBody()
		if err != nil {
			return "", err
		}
		d.cachedDashboardIndexContent = body
	}
	return d.cachedDashboardIndexContent, nil
}

type Basic struct {
	ID   string
	Name string
}

func (d *DashboardHandler) GetBasic() (*Basic, error) {
	body, err := d.getCachedDashboardIndexBody()
	if err != nil {
		return nil, err
	}
	id, _ := utils.MatchSingle(regexp.MustCompile(`用户名</label>(.+?)</li>`), body)
	name, _ := utils.MatchSingle(regexp.MustCompile(`姓名</label>(.+?)</li>`), body)
	return &Basic{
		ID:   id,
		Name: name,
	}, nil
}

type Package struct {
	PackageCost  string
	UsedTraffic  string
	UsedDuration string
	Balance      string
	// Balance < 0
	Overdue bool
}

func (d *DashboardHandler) GetPackage() (*Package, error) {
	body, err := d.getCachedDashboardIndexBody()
	if err != nil {
		return nil, err
	}

	infos, _ := utils.MatchMultiple(regexp.MustCompile(`<td data-col-seq="3">(.+?)</td><td data-col-seq="4">(.+?)</td><td data-col-seq="6">(.+?)</td><td data-col-seq="7">(.+?)</td>`), body)
	if len(infos) < 1 {
		return nil, errors.New("fail to get package info")
	}
	info := infos[0]

	result := &Package{
		UsedTraffic:  info[1],
		UsedDuration: info[2],
		PackageCost:  info[3],
		Balance:      info[4],
	}

	b, _ := strconv.ParseFloat(info[4], 32)
	if b < 0 {
		result.Overdue = true
	}

	return result, nil
}

type Device struct {
	ID        int
	IP        string
	StartTime string
	Stage     string
	SID       string
}

func (d *DashboardHandler) GetDevice() ([]Device, error) {
	body, err := d.getCachedDashboardIndexBody()
	if err != nil {
		return []Device{}, err
	}
	ds, _ := utils.MatchMultiple(regexp.MustCompile(`<tr data-key="(\d+)"><td data-col-seq="0">\d+</td><td data-col-seq="1">(.+?)</td><td data-col-seq="3">(.+?)</td><td data-col-seq="7">(.+?)</td><td data-col-seq="9">.+?</td>`), body)
	result := make([]Device, len(ds))
	for i, device := range ds {
		result[i] = Device{i, device[2], device[3], device[4], device[1]}
	}
	return result, nil
}

type BillRecord struct {
	ID           string
	Cost         float64
	Traffic      string
	UsedDuration string
	Date         string
}

func (d *DashboardHandler) GetBill(page int) ([]BillRecord, error) {
	body, err := d.fetchDashboardBillsBody(page)
	if err != nil {
		return []BillRecord{}, nil
	}
	t, _ := utils.MatchSingle(regexp.MustCompile(`<title>(.+?)</title>`), body)
	if t != "结算清单" {
		return []BillRecord{}, errors.New("error occurs when parsing bill page")
	}

	bills, _ := utils.MatchMultiple(regexp.MustCompile(`<td data-col-seq="0">(\d+)</td><td data-col-seq="1">.+?</td><td data-col-seq="2">(.+?)</td><td data-col-seq="3">(.+?)</td><td style="display: none;" data-col-seq="6">.+?</td><td data-col-seq="7">(.+?)</td><td data-col-seq="10">(\d+)</td><td data-col-seq="12">(.+?)</td></tr>`), body)

	// b1   b2    b3       b4    b5
	// id  cost  traffic  time  date
	result := make([]BillRecord, len(bills))

	for i, b := range bills {
		// 总消费为固定费用+实时费用
		fixedCost, _ := strconv.ParseFloat(b[2], 32)
		realtimeCost, _ := strconv.ParseFloat(b[3], 32)
		result[i] = BillRecord{b[1], fixedCost + realtimeCost, b[4], b[5], b[6]}
	}
	return result, nil
}

type UsageRecord struct {
	StartTime    string
	EndTime      string
	IP           string
	Traffic      string
	UsedDuration string
}

func (d *DashboardHandler) GetUsageRecords(page int) ([]UsageRecord, error) {
	body, err := d.fetchDashboardUsageLogBody(page)
	if err != nil {
		return []UsageRecord{}, nil
	}
	t, _ := utils.MatchSingle(regexp.MustCompile(`<title>(.+?)</title>`), body)
	if t != "上网明细" {
		return []UsageRecord{}, errors.New("error occurs when parsing usage log page")
	}

	hs, _ := utils.MatchMultiple(regexp.MustCompile(`<td data-col-seq="0">.+?</td><td data-col-seq="1">(.+?)</td><td data-col-seq="2">(.+?)</td><td data-col-seq="5">(.+?)</td><td data-col-seq="10">.+?</td><td data-col-seq="12">.+?</td><td style="display: none;" data-col-seq="16">.+?</td><td data-col-seq="17">(.+?)</td><td style="display: none;" data-col-seq="18">.+?</td><td data-col-seq="19">(.+?)</td><td data-col-seq="20">.+?</td>`), body)
	result := make([]UsageRecord, len(hs))
	for i, h := range hs {
		result[i] = UsageRecord{h[1], h[2], h[3], h[4], h[5]}
	}
	return result, nil
}

type RechargeRecord struct {
	ID   string
	Cost string
	Time string
}

func (d *DashboardHandler) GetRecharge(page int) ([]RechargeRecord, error) {
	body, err := d.fetchDashboardRechargeBody(page)
	if err != nil {
		return []RechargeRecord{}, err
	}
	t, _ := utils.MatchSingle(regexp.MustCompile(`<title>(.+?)</title>`), body)
	if t != "缴费清单" {
		return []RechargeRecord{}, errors.New("error occurs when parsing recharge page")
	}

	rs, _ := utils.MatchMultiple(regexp.MustCompile(`<td data-col-seq="0">(.+?)</td><td data-col-seq="1">\d+</td><td data-col-seq="2">(.+?)</td><td data-col-seq="3">.+?</td><td data-col-seq="6">(.+?)</td></tr>`), body)

	result := make([]RechargeRecord, len(rs))
	// r1  r2     r3
	// id  cost  time
	for i, r := range rs {
		result[i] = RechargeRecord{r[1], r[2], r[3]}
	}
	return result, nil
}

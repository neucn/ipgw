package handler

import (
	"errors"
	"fmt"
	"github.com/neucn/ipgw/pkg/model"
	"github.com/neucn/ipgw/pkg/utils"
	"github.com/neucn/neugo"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	return nil
}

func (d *DashboardHandler) fetchDashboardIndexBody() (string, error) {
	resp, err := d.client.Get("http://ipgw.neu.edu.cn:8800/sso/default/neusoft")
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) fetchDashboardBillsBody(page int) (string, error) {
	resp, err := d.client.Get(fmt.Sprintf("http://ipgw.neu.edu.cn:8800/financial/checkout/list?page=%d&per-page=10", page))
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) fetchDashboardRechargeBody(page int) (string, error) {
	resp, err := d.client.Get(fmt.Sprintf("http://ipgw.neu.edu.cn:8800/financial/pay/list?page=%d&per-page=10", page))
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (d *DashboardHandler) fetchDashboardUsageLogBody(page int) (string, error) {
	resp, err := d.client.Get(fmt.Sprintf("http://ipgw.neu.edu.cn:8800/log/detail/index?page=%d&per-page=20", page))
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
	id, _ := utils.MatchSingle(regexp.MustCompile(`账号</label>\s+(\d+)\s+`), body)
	name, _ := utils.MatchSingle(regexp.MustCompile(`姓名</label>\s+(.+?)\s+`), body)
	return &Basic{
		ID:   id,
		Name: name,
	}, nil
}

type Package struct {
	// Unit is G
	PackageTraffic string
	PackageCost    string
	UsedTraffic    string
	UsedDuration   string
	UsedTimes      string
	Balance        string
	// Balance < 0
	Overdue bool
	// UsedTraffic > PackageTraffic
	ExcessPackageTraffic bool
}

func (d *DashboardHandler) GetPackage() (*Package, error) {
	body, err := d.getCachedDashboardIndexBody()
	if err != nil {
		return nil, err
	}

	infos, _ := utils.MatchMultiple(regexp.MustCompile(`<td>\W+.+?(\d+?)G下行流量(.+?)元?/.+?</td>\W+<td>\W+(.+?)\W+</td>\W+<td>(.+?)</td>\W+<td>(.+?)</td>\W+<td>.+?</td>\W+<td>(.+?)</td>\W+<td>.+?</td>`), body)
	if len(infos) < 1 {
		return nil, errors.New("fail to get package info")
	}
	info := infos[0]

	result := &Package{
		UsedTraffic:  info[3],
		UsedDuration: info[4],
		UsedTimes:    info[5],
		Balance:      info[6],
	}

	b, _ := strconv.ParseFloat(info[6], 32)
	if b < 0 {
		result.Overdue = true
	}

	result.PackageTraffic = info[1]

	if info[2] == "免费" {
		result.PackageCost = "0"
	} else {
		result.PackageCost = info[2]
	}

	if strings.HasSuffix(info[3], "G") {
		u, _ := strconv.ParseFloat(strings.TrimSuffix(info[6], "G"), 32)
		t, _ := strconv.ParseFloat(info[1], 32)
		if u > t {
			result.ExcessPackageTraffic = true
		}
	}

	return result, nil
}

type Device struct {
	ID        int
	IP        string
	StartTime string
	SID       string
}

func (d *DashboardHandler) GetDevice() ([]Device, error) {
	body, err := d.getCachedDashboardIndexBody()
	if err != nil {
		return []Device{}, err
	}
	ds, _ := utils.MatchMultiple(regexp.MustCompile(`<td>\d+</td>\W+?<td>(.+?)</td>\W+?<td>.+?</td>\W+?<td>(.+?)</td>\W+?<td>.+?</td>\W+?<td><a id="(\d+)".+?下线</a></td>`), body)
	result := make([]Device, len(ds))
	for i, device := range ds {
		result[i] = Device{i, device[1], device[2], device[3]}
	}
	return result, nil
}

type BillRecord struct {
	ID           string
	Cost         string
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

	bills, _ := utils.MatchMultiple(regexp.MustCompile(`<td>(\d+?)</td><td>\d+?</td><td>(\d+?)</td><td>.+?</td><td>.+?</td><td>(.+?)</td><td>(.+?)</td><td>.+?</td><td>(.+?)</td>`), body)

	// b1   b2    b3       b4    b5
	// id  cost  traffic  time  date
	result := make([]BillRecord, len(bills))
	for i, b := range bills {
		result[i] = BillRecord{b[1], b[2], b[3], b[4], strings.Split(b[5], " ")[0]}
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

	hs, _ := utils.MatchMultiple(regexp.MustCompile(`<td>\d+?</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>(.+?)</td><td>.+?</td></tr>`), body)
	result := make([]UsageRecord, len(hs))
	for i, h := range hs {
		result[i] = UsageRecord{h[1], h[2], h[3], h[4], h[5]}
	}
	return result, nil
}

type RechargeRecord struct {
	ID     string
	Cost   string
	Method string
	Time   string
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

	rs, _ := utils.MatchMultiple(regexp.MustCompile(`<td>(\d+?)</td><td>\d+?</td><td>(\d+?)</td><td>.+?</td><td>(.+?)</td><td>.+?</td><td>(.+?)</td><td>.+?</td>`), body)

	result := make([]RechargeRecord, len(rs))
	// r1  r2     r3      r4
	// id  cost  method  time
	for i, r := range rs {
		result[i] = RechargeRecord{r[1], r[2], r[3], r[4]}
	}
	return result, nil
}

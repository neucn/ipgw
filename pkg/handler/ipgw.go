package handler

import (
	"errors"
	"github.com/neucn/ipgw/pkg/model"
	"github.com/neucn/ipgw/pkg/utils"
	"github.com/neucn/neugo"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

type IpgwHandler struct {
	info   *model.Info
	client *http.Client
}

func (h *IpgwHandler) GetInfo() *model.Info {
	return h.info
}

func (h *IpgwHandler) GetClient() *http.Client {
	return h.client
}

func NewIpgwHandler() *IpgwHandler {
	return &IpgwHandler{
		info:   &model.Info{},
		client: neugo.NewSession(),
	}
}

// Login will call FetchBasicInfoWithoutLogin
func (h *IpgwHandler) Login(account *model.Account) error {
	if account.Cookie != "" {
		h.client.Jar.SetCookies(&url.URL{
			Scheme: "https",
			Host:   "ipgw.neu.edu.cn",
		}, []*http.Cookie{{
			Name:   "session_for%3Asrun_cas_php",
			Value:  account.Cookie,
			Domain: "ipgw.neu.edu.cn",
		}})
	} else {
		password, err := account.GetPassword()
		if err != nil {
			return err
		}
		err = neugo.Use(h.client).WithAuth(account.Username, password).Login(neugo.CAS)
		if err != nil {
			return err
		}
	}
	return h.FetchBasicInfo()
}

func (h *IpgwHandler) FetchUsageInfo() error {
	body, err := h.getRawUsageInfo()
	if err != nil {
		return err
	}
	items := strings.Split(body, ",")

	if len(items) != 6 {
		return errors.New("usage info is incomplete")
	}

	h.info.Traffic, _ = strconv.Atoi(items[0])
	h.info.UsedTime, _ = strconv.Atoi(items[1])
	h.info.Balance, _ = strconv.ParseFloat(items[2], 64)

	return nil
}

func (h *IpgwHandler) getRawIpgwPage() (string, error) {
	resp, err := h.client.Get("https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (h *IpgwHandler) getRawUsageInfo() (string, error) {
	req, _ := http.NewRequest("POST", "https://ipgw.neu.edu.cn/include/auth_action.php",
		strings.NewReader("action=get_online_info"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Referer", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (h *IpgwHandler) getRawOldIpgwPage() (string, error) {
	resp, err := h.client.Get("http://ipgw.neu.edu.cn/srun_portal_pc_succeed.php")
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func getUsernameAndIPFromOld(body string) (username, ip string) {
	username, _ = utils.MatchSingle(regexp.MustCompile(`name="username" value="(.+?)"`), body)
	ip, _ = utils.MatchSingle(regexp.MustCompile(`name="user_ip" value="(.+?)"`), body)
	return
}

func getUsernameAndIP(body string) (username, ip string) {
	username, _ = utils.MatchSingle(regexp.MustCompile(`id="user_name" style=".+?">(.+?)</span>`), body)
	ip, _ = utils.MatchSingle(regexp.MustCompile(`id="user_ip" style=".+?">(.+?)</span>`), body)
	return
}

func isOverdue(body string) (out bool) {
	return regexp.MustCompile(`余额不足月租`).MatchString(body)
}

func (h *IpgwHandler) FetchBasicInfo() error {
	body, err := h.getRawIpgwPage()
	if err != nil {
		return err
	}
	h.info.Username, h.info.IP = getUsernameAndIP(body)
	h.info.Overdue = isOverdue(body)
	return nil
}

func (h *IpgwHandler) FetchBasicInfoWithoutLogin() error {
	body, err := h.getRawOldIpgwPage()
	if err != nil {
		return err
	}
	h.info.Username, h.info.IP = getUsernameAndIPFromOld(body)
	return nil
}

func (h *IpgwHandler) Logout() error {
	req, _ := http.NewRequest("POST", "http://ipgw.neu.edu.cn/srun_portal_pc_succeed.php",
		strings.NewReader("action=auto_logout&username="+h.info.Username))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Referer", "http://ipgw.neu.edu.cn/srun_portal_pc_succeed.php")

	_, err := h.client.Do(req)
	return err
}

// IsLoggedIn will fetch basic info by calling FetchBasicInfoWithoutLogin
func (h *IpgwHandler) IsLoggedIn() bool {
	err := h.FetchBasicInfoWithoutLogin()
	return err == nil && h.info.Username != ""
}

func (h *IpgwHandler) IsConnected() bool {
	body, err := h.getRawOldIpgwPage()
	if err != nil && utils.IsNetworkError(err) {
		return false
	}
	_, ip := getUsernameAndIPFromOld(body)
	return ip != ""
}

func (h *IpgwHandler) Kick(sid string) (bool, error) {
	req, _ := http.NewRequest("POST", "https://ipgw.neu.edu.cn/srun_cas.php",
		strings.NewReader("action=dm&sid="+sid))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "https://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	resp, err := h.client.Do(req)
	if err != nil {
		return false, err
	}
	result := utils.ReadBody(resp)
	return result == "下线请求已发送", nil
}

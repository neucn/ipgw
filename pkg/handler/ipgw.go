package handler

import (
	"errors"
	"fmt"
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

func (h *IpgwHandler) Login(account *model.Account) error {
	var err error
	var body string
	if account.Cookie != "" {
		body, err = h.loginCookie(account.Cookie)
		return h.ParseBasicInfo(body, false)
	}

	var password string
	password, err = account.GetPassword()
	if err != nil {
		return err
	}
	if account.NonUnified {
		body, err = h.loginNonUnified(account.Username, password)
	} else {
		body, err = h.loginUnified(account.Username, password)
	}

	if err != nil {
		return err
	}

	return h.ParseBasicInfo(body, account.NonUnified)
}

func (h *IpgwHandler) loginNonUnified(username, password string) (string, error) {
	req, _ := http.NewRequest(http.MethodPost, "http://ipgw.neu.edu.cn/srun_portal_pc.php?ac_id=1&", strings.NewReader(
		fmt.Sprintf("action=login&ac_id=1&user_ip=&nas_ip=&user_mac=&url=&username=%s&password=%s&save_me=0",
			username, url.QueryEscape(password))))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (h *IpgwHandler) loginCookie(cookie string) (string, error) {
	h.client.Jar.SetCookies(&url.URL{
		Scheme: "http",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{{
		Name:   "session_for%3Asrun_cas_php",
		Value:  cookie,
		Domain: "ipgw.neu.edu.cn",
	}})
	return h.getRawIpgwPage()
}

func (h *IpgwHandler) loginUnified(username, password string) (string, error) {
	if err := neugo.Use(h.client).WithAuth(username, password).Login(neugo.CAS); err != nil {
		return "", err
	}
	return h.getRawIpgwPage()
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
	resp, err := h.client.Get("http://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (h *IpgwHandler) getRawUsageInfo() (string, error) {
	req, _ := http.NewRequest("POST", "http://ipgw.neu.edu.cn/include/auth_action.php",
		strings.NewReader("action=get_online_info"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
	req.Header.Add("Referer", "http://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

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
	username, _ = utils.MatchSingle(regexp.MustCompile(`id="user_name".*?>(.+?)</span>`), body)
	ip, _ = utils.MatchSingle(regexp.MustCompile(`id="user_ip".*?>(.+?)</span>`), body)
	return
}

func isOverdue(body string) (out bool) {
	return regexp.MustCompile(`余额不足月租`).MatchString(body)
}

func (h *IpgwHandler) ParseBasicInfo(body string, nonUnified bool) error {
	if nonUnified {
		// the result of non-unified method is embedded from the old ipgw page so use getUsernameAndIPFromOld to parse
		h.info.Username, h.info.IP = getUsernameAndIPFromOld(body)
	} else {
		h.info.Username, h.info.IP = getUsernameAndIP(body)
	}
	h.info.Overdue = isOverdue(body)
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

func (h *IpgwHandler) IsConnectedAndLoggedIn() (connected bool, loggedIn bool) {
	body, err := h.getRawOldIpgwPage()
	if err != nil && utils.IsNetworkError(err) {
		return false, false
	}
	h.info.Username, h.info.IP = getUsernameAndIPFromOld(body)
	return h.info.IP != "", h.info.Username != ""
}

func (h *IpgwHandler) Kick(sid string) (bool, error) {
	req, _ := http.NewRequest("POST", "http://ipgw.neu.edu.cn/srun_cas.php",
		strings.NewReader("action=dm&sid="+sid))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", "http://ipgw.neu.edu.cn/srun_cas.php?ac_id=1")

	resp, err := h.client.Do(req)
	if err != nil {
		return false, err
	}
	result := utils.ReadBody(resp)
	return result == "下线请求已发送", nil
}

package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/neucn/ipgw/pkg/model"
	"github.com/neucn/ipgw/pkg/utils"
	"github.com/neucn/neugo"
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
	var (
		password string
		err      error
	)
	if account.Cookie != "" {
		_, err = h.loginCookie(account.Cookie) // 通过cookie登录
	} else {
		password, err = account.GetPassword()
		if err != nil {
			return err
		}
		_, err = h.login(account.Username, password) // 通过用户名、密码登录
	}

	if err != nil {
		return err
	}

	return h.ParseBasicInfo() // 解析信息
}

func (h *IpgwHandler) loginCookie(cookie string) (string, error) {
	h.client.Jar.SetCookies(&url.URL{
		Scheme: "https",
		Host:   "ipgw.neu.edu.cn",
	}, []*http.Cookie{{
		Name:   "session_for%3Asrun_cas_php",
		Value:  cookie,
		Domain: "ipgw.neu.edu.cn",
	}})
	return h.requestLoginApi()
}

func (h *IpgwHandler) login(username, password string) (string, error) {
	if err := neugo.Use(h.client).WithAuth(username, password).Login(neugo.CAS); err != nil {
		return "", err
	}
	return h.requestLoginApi()
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

func (h *IpgwHandler) requestLoginApi() (string, error) {
	// 统一认证拿到ticket
	resp, err := h.client.Get("https://pass.neu.edu.cn/tpass/login?service=http://ipgw.neu.edu.cn/srun_portal_sso?ac_id=1")
	if err != nil {
		return "", err
	}
	// 使用ticket调用api登录
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/v1"+resp.Request.URL.RequestURI(), nil)
	resp, _ = h.client.Do(req)
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

func (h *IpgwHandler) getJsonIpgwData() (string, error) {
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/cgi-bin/rad_user_info", nil)
	req.Header.Set("Accept", "application/json;")
	resp, err := h.client.Do(req)
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func getUsernameAndIPFromJson(body string) (username, ip string) {
	data := make(map[string]interface{})
	json.Unmarshal([]byte(body), &data)
	if data["error"].(string) != "ok" {
		return "", data["client_ip"].(string)
	}
	return data["user_name"].(string), data["online_ip"].(string)
}

func isOverdue(body string) (out bool) {
	return regexp.MustCompile(`余额不足月租`).MatchString(body)
}

func (h *IpgwHandler) ParseBasicInfo() error {
	body, _ := h.getJsonIpgwData()
	h.info.Username, h.info.IP = getUsernameAndIPFromJson(body)
	h.info.Overdue = isOverdue(body)
	return nil
}

func (h *IpgwHandler) Logout() error {
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/cgi-bin/srun_portal?action=logout&username="+h.info.Username, nil)
	req.Header.Add("Referer", "http://ipgw.neu.edu.cn/srun_portal_success?ac_id=1")
	_, err := h.client.Do(req)
	return err
}

func (h *IpgwHandler) IsConnectedAndLoggedIn() (connected bool, loggedIn bool) {
	// 调用ipgw信息api
	body, err := h.getJsonIpgwData()
	if err != nil && utils.IsNetworkError(err) {
		return false, false
	}
	h.info.Username, h.info.IP = getUsernameAndIPFromJson(body)
	return h.info.IP != "", h.info.Username != ""
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

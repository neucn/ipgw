package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/neucn/ipgw/pkg/model"
	"github.com/neucn/ipgw/pkg/utils"
	"github.com/neucn/neugo"
)

var (
	once sync.Once
)

type IpgwHandler struct {
	info    *model.Info
	client  *http.Client
	oriInfo map[string]interface{}
}

func (h *IpgwHandler) GetInfo() *model.Info {
	return h.info
}

func (h *IpgwHandler) GetClient() *http.Client {
	return h.client
}

func NewIpgwHandler() *IpgwHandler {
	return &IpgwHandler{
		info:    &model.Info{},
		client:  neugo.NewSession(),
		oriInfo: make(map[string]interface{}),
	}
}

func (h *IpgwHandler) Login(account *model.Account) error {
	var (
		password string
		body     string
		err      error
	)
	if account.Cookie != "" {
		body, err = h.loginCookie(account.Cookie) // 通过cookie登录
	} else {
		password, err = account.GetPassword()
		if err != nil {
			return err
		}
		body, err = h.login(account.Username, password) // 通过用户名、密码登录
	}

	if err != nil {
		return err
	}

	if strings.Contains(body, "Arrearage users") {
		return fmt.Errorf("overdue")
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

func (h *IpgwHandler) NEUAuth(username, password string) error {
	// 仅登录统一认证
	return neugo.Use(h.client).WithAuth(username, password).Login(neugo.CAS)
}

func (h *IpgwHandler) login(username, password string) (string, error) {
	if err := h.NEUAuth(username, password); err != nil {
		return "", err
	}
	return h.requestLoginApi()
}

func (h *IpgwHandler) FetchUsageInfo() error {
	err := h.getJsonIpgwData()
	if err != nil {
		return err
	}

	h.info.Traffic = int(h.oriInfo["sum_bytes"].(float64))
	h.info.UsedTime = int(h.oriInfo["sum_seconds"].(float64))
	h.info.Balance, _ = h.oriInfo["user_balance"].(float64)

	return nil
}

func (h *IpgwHandler) requestLoginApi() (string, error) {
	// 获取当前网络下对应网关url的query参数
	resp, err := h.client.Get("https://ipgw.neu.edu.cn/")
	if err != nil {
		return "", err
	}
	// 统一认证拿到ticket
	resp, err = h.client.Get("https://pass.neu.edu.cn/tpass/login?service=http://ipgw.neu.edu.cn/srun_portal_sso?" + resp.Request.URL.RawQuery)
	if err != nil {
		return "", err
	}
	// 使用ticket调用api登录
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/v1"+resp.Request.URL.RequestURI(), nil)
	resp, err = h.client.Do(req)
	if err != nil {
		return "", err
	}
	return utils.ReadBody(resp), nil
}

func (h *IpgwHandler) getJsonIpgwData() error {
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/cgi-bin/rad_user_info", nil)
	req.Header.Set("Accept", "application/json;")
	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	oriInfoBytes, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	return json.Unmarshal(oriInfoBytes, &h.oriInfo)
}

func getUsernameAndIPFromJson(data map[string]interface{}) (username, ip string) {
	errorMsg, ok := data["error"].(string)
	if !ok {
		return "", ""
	}
	if errorMsg != "ok" {
		return "", data["client_ip"].(string)
	}
	return data["user_name"].(string), data["online_ip"].(string)
}

func (h *IpgwHandler) ParseBasicInfo() error {
	h.getJsonIpgwData()
	h.info.Username, h.info.IP = getUsernameAndIPFromJson(h.oriInfo)
	return nil
}

func (h *IpgwHandler) Logout() error {
	req, _ := http.NewRequest("GET", "https://ipgw.neu.edu.cn/cgi-bin/srun_portal?action=logout&username="+h.info.Username, nil)
	req.Header.Add("Referer", "https://ipgw.neu.edu.cn/srun_portal_success?ac_id=1")
	_, err := h.client.Do(req)
	return err
}

func (h *IpgwHandler) IsConnectedAndLoggedIn() (connected bool, loggedIn bool) {
	// 调用ipgw信息api
	err := h.getJsonIpgwData()
	if err != nil && utils.IsNetworkError(err) {
		return false, false
	}
	h.info.Username, h.info.IP = getUsernameAndIPFromJson(h.oriInfo)
	return h.info.IP != "", h.info.Username != ""
}

func (h *IpgwHandler) Kick(sid string) (bool, error) {
	once.Do(func() {
		h.client.Get("https://ipgw.neu.edu.cn:8800/sso/neusoft/index")
	})
	// 请求主页
	resp, err := h.client.Get("https://ipgw.neu.edu.cn:8800/home")
	if err != nil {
		return false, err
	}
	body := utils.ReadBody(resp)
	// 获取csrf-token
	token, _ := utils.MatchSingle(regexp.MustCompile(`<meta name="csrf-token" content="(.+?)">`), body)

	req, _ := http.NewRequest("POST", "https://ipgw.neu.edu.cn:8800/home/delete?id="+sid, strings.NewReader("_csrf-8800="+token))
	req.Header.Set("Referer", "https://ipgw.neu.edu.cn:8800/home/index")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err = h.client.Do(req)
	if err != nil {
		return false, err
	}
	result := utils.ReadBody(resp)
	return strings.Contains(result, "下线请求已发出"), nil
}

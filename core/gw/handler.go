package gw

import (
	"regexp"
)

// 若id为空，则没有重复登陆
// 若id存在而sid为空，目前不知道什么情况下会出现
// id与sid都存在，重复登陆
func IsLoginRepeatedly(body string) (id, sid string) {
	idExp := regexp.MustCompile(`aaa\n(\d+?)ccc`)
	idList := idExp.FindAllStringSubmatch(body, -1)
	if len(idList) < 1 {
		return "", ""
	}

	id = idList[0][1]

	sidExp := regexp.MustCompile(`btn-dark" href="javascript\(0\);" onclick="do_drop\('(\d+?)'\);`)
	sidList := sidExp.FindAllStringSubmatch(body, -1)
	if len(sidList) < 1 {
		return id, ""
	}

	return id, sidList[0][1]
}

// 判断是否欠费
func IsOverdue(body string) (out bool) {
	outExp := regexp.MustCompile(`余额不足月租`)
	return outExp.MatchString(body)
}

// 获取学号(用于使用Cookie登陆)
func GetID(body string) (id string) {
	usernameExp := regexp.MustCompile(`user_name" style="float:right;color: #894324;">(.+?)</span>`)
	username := usernameExp.FindAllStringSubmatch(body, -1)

	if len(username) == 0 {
		return ""
	}
	return username[0][1]
}

// 获取SID和IP
func GetSIDAndIP(body string) (sid, ip string) {
	// 匹配IP
	ipExp := regexp.MustCompile(`get_online_info\('(.+?)'\)`)
	ips := ipExp.FindAllStringSubmatch(body, -1)

	if len(ips) != 0 {
		ip = ips[0][1]
	}

	// 匹配SID
	sidExp := regexp.MustCompile(`background:lightgreen[\w\W]+?onclick="do_drop\('(\d+)'\)`)
	sids := sidExp.FindAllStringSubmatch(body, -1)
	if len(sids) != 0 {
		sid = sids[0][1]
	}

	return
}

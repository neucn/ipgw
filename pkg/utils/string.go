package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"regexp"
)

var errorNoMatched = errors.New("no matched")

func MatchMultiple(re *regexp.Regexp, content string) ([][]string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return nil, errorNoMatched
	}
	return matched, nil
}

func MatchSingle(re *regexp.Regexp, content string) (string, error) {
	matched := re.FindAllStringSubmatch(content, -1)
	if len(matched) < 1 {
		return "", errorNoMatched
	}
	return matched[0][1], nil
}

func ReadBody(resp *http.Response) (body string) {
	res, _ := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	return string(res)
}

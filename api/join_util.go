package api

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
)


func JoinCookieList(oldCookie string, newCookie []*http.Cookie) string {
	var build strings.Builder
	if len(oldCookie) > 0 {
		build.WriteString(oldCookie)
		if len(newCookie) > 0 {
			build.WriteString("; ")
		}
	}
	for _, cookie := range newCookie {
		if cookie.Value == "" {
			continue
		}
		item := fmt.Sprintf("%v=%v; ", cookie.Name, cookie.Value)
		build.WriteString(item)
	}
	cookieList := build.String()
	if len(cookieList) > 2 {
		return cookieList[0 : len(cookieList)-2]
	}
	return cookieList
}

func CreateFormReader(data map[string]interface{}) io.Reader {
	form := url.Values{}
	for k, v := range data {
		val := fmt.Sprintf("%v", v)
		form.Add(k, val)
	}
	return strings.NewReader(form.Encode())
}

// 随机生成未被占用的端口号
func PickUnusedPort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	if err := l.Close(); err != nil {
		return 0, err
	}
	return port, nil
}

func MapToJsonStr(param map[string]interface{}) string {
	marshal, err := json.Marshal(param)
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}

func UrlAppendParam(param map[string]interface{}) string {
	if param == nil {
		return ""
	}
	urlParams := ""
	for k, v := range param {
		if v != "" {
			urlParams = fmt.Sprintf("%v%v=%v&", urlParams, k, v)
		}
	}
	if urlParams != "" {
		return "?" + urlParams[0:len(urlParams)-1]
	}
	return urlParams
}

func UrlAppendParam2(param map[string]string) string {
	urlParams := ""
	for k, v := range param {
		if v != "" {
			urlParams = fmt.Sprintf("%v%v=%v&", urlParams, k, v)
		}
	}
	if urlParams != "" {
		return "?" + urlParams[0:len(urlParams)-1]
	}
	return urlParams
}

func RandomBoundary() string {
	var buf [30]byte
	_, err := io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", buf[:])
}

func RandomUserAgent() string {
	userAgent := "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3100.0 Safari/537.36"
	return userAgent
}

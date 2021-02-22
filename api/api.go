package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Method string
type Params map[string]interface{}

type Session struct {
	client    *http.Client
	Header    http.Header
	RespCode  int
	RespData  []byte
	Cookie    []*http.Cookie
	notHeader bool
}

type Result []byte

func (session *Session) DefaultClient() {
	session.client = http.DefaultClient
}

func (session *Session) ClientProxy(proxy string) {
	u := url.URL{}
	if proxy == "" {
		session.DefaultClient()
	}
	if !strings.HasPrefix(proxy, "http") {
		session.DefaultClient()
	}
	urlProxy, _ := u.Parse(proxy)
	c := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(urlProxy),
		},
		Timeout: time.Duration(15) * time.Second,
	}
	session.client = &c
}

func (session *Session) SetHeader(hdr http.Header) {
	session.Header = hdr
}

func (session *Session) AddHeader(maps map[string]string) {
	header := session.Header
	if header == nil {
		header = http.Header{}
	}
	for k, v := range maps {
		header.Set(k, v)
	}
	session.Header = header
}

func (session *Session) SetCookie(cookie string) {
	hdr := http.Header{}
	if session.Header != nil {
		hdr = session.Header
	}
	hdr.Set("cookie", cookie)
	session.Header = hdr
}

func (session *Session) Get(path string, params Params) error {
	return session.Api(path, http.MethodGet, params)
}

func (session *Session) Post(path string, params Params) error {
	return session.Api(path, http.MethodPost, params)
}

func (session *Session) PostForUrl(path string, params Params) error {
	session.notHeader = true
	return session.Api(path, http.MethodPost, params)
}

func (session *Session) PostForJson(path string, params Params) error {
	header := session.Header
	if header == nil {
		header = http.Header{}
	}
	header.Set("Content-Type", "application/json")
	session.SetHeader(header)
	return session.Api(path, http.MethodPost, params)
}

func (session *Session) Api(path string, method Method, params Params) error {
	return session.graph(path, method, params)
}

func (session *Session) graph(path string, method Method, params Params) error {

	if params == nil {
		params = Params{}
	}
	if method == http.MethodGet {
		path = fmt.Sprintf("%s%s", path, UrlAppendParam(params))
		err := session.sendGetRequest(path)
		if err != nil {
			return err
		}
	} else if method == http.MethodPost {
		err := session.sendPostRequest(path, params)
		if err != nil {
			return err
		}
	}

	return nil
}

func (session *Session) sendGetRequest(uri string) error {
	log.Printf("请求的接口为 %s\n", uri)
	//log.Printf("请求头为 %v\n", session.Header)
	parsedURL, err := url.Parse(uri)
	req := &http.Request{
		Method:     http.MethodGet,
		URL:        parsedURL,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     session.Header,
	}

	response, data, err := session.sendRequest(req)

	if err != nil {
		return err
	}
	session.RespData = data
	session.Cookie = response.Cookies()
	return err
}
func (session *Session) sendPostRequest(uri string, params Params) error {

	var rc io.Reader
	log.Printf("请求的接口为 %s\n", uri)
	if session.Header == nil {
		session.Header = http.Header{}
	}
	contentType := session.Header.Get("Content-Type")

	if session.notHeader {
		uri = fmt.Sprintf("%s%s", uri, UrlAppendParam(params))
	} else if strings.Contains(contentType, "json") {
		jsonParams, err := json.Marshal(params)
		if err != nil {
			return fmt.Errorf("post params json encode error： %v", err)
		}
		rc = bytes.NewReader(jsonParams)
	} else {
		if contentType == "" {
			contentType = "application/x-www-form-urlencoded"
			session.Header.Set("Content-Type", contentType)
		}
		rc = createFormReader(params)
	}
	request, err := http.NewRequest(http.MethodPost, uri, rc)
	if err != nil {
		log.Println("req error", err.Error())
		return err
	}
	request.Header = session.Header

	//marshal, _ := json.Marshal(session.Header)
	//log.Printf("请求类型为 %v\n", contentType)
	//log.Printf("请求头为 %v\n", string(marshal))

	response, data, err := session.sendRequest(request)

	if err != nil {
		return err
	}
	session.RespData = data
	session.Cookie = response.Cookies()
	session.RespCode = response.StatusCode

	//cookie, _ := json.Marshal(response.Cookies())
	//log.Println("返回cookie:", string(cookie))
	//log.Println("====================")
	if response.StatusCode != 200 {
		return fmt.Errorf("请求失败：%s", string(data))
	}
	return nil
}

func (session *Session) sendRequest(request *http.Request) (response *http.Response, data []byte, err error) {

	if session.client == nil {
		response, err = http.DefaultClient.Do(request)
	} else {
		response, err = session.client.Do(request)
	}

	if err != nil {
		err = fmt.Errorf("发送http请求失败: %v", err)
		return
	}

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, response.Body)
	_ = response.Body.Close()

	if err != nil {
		log.Println("错误:", err.Error())
		err = fmt.Errorf("http response error : %v", err)
	}

	data = buf.Bytes()
	//log.Printf("[返回]: code:%s , %s\n", response.Status, string(data))
	return
}

func createFormReader(params Params) io.Reader {
	form := url.Values{}
	for k, v := range params {
		form.Add(k, fmt.Sprintf("%v", v))
	}
	log.Fatalf("请求的接口：%s", form.Encode())
	return strings.NewReader(form.Encode())
}

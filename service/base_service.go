package service

import (
	"encoding/json"
	"fmt"
	"log"
	"seckill-request/api"
)

var BaseUrl = "http://127.0.0.1:7087"
var ApiRegister = "/user/register"
var ApiLogin = "/user/login"
var ApiAddress = "/auth/user/address"
var ApiCharge = "/auth/fund/charge"
var ApiGoods = "/goods/"
var ApiBuy = "/auth/order/buy"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Register(userName, passWord string) {
	session := &api.Session{}
	session.DefaultClient()
	err := session.PostForJson(fmt.Sprintf("%s%s", BaseUrl, ApiRegister), api.Params{
		"user_name": userName,
		"pass_word": passWord,
	})
	if err != nil {
		log.Fatalln("注册请求失败", err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil {
		log.Fatalln("注册失败", err.Error())
	}
	if resp.Code == 0 {
		log.Println("注册成功：", userName)
	} else {
		log.Fatalf("注册失败:[返回] %s", string(session.RespData))
	}
}

func Login(userName, passWord string) string {
	session := &api.Session{}
	session.DefaultClient()
	err := session.PostForJson(fmt.Sprintf("%s%s", BaseUrl, ApiLogin), api.Params{
		"user_name": userName,
		"pass_word": passWord,
	})
	if err != nil {
		log.Fatalln("登录请求失败", err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil {
		log.Fatalln("登录失败", err.Error())
	}
	if resp.Code == 0 {
		log.Println("登录成功：", userName)
	} else {
		log.Fatalf("登录失败:[返回] %s", string(session.RespData))
	}
	token := resp.Data.(map[string]interface{})["token"]

	return token.(string)
}

func AddAddress(province, city, detail, token string) int64 {
	session := &api.Session{}
	session.DefaultClient()
	session.AddHeader(map[string]string{
		"token": token,
	})
	err := session.PostForJson(fmt.Sprintf("%s%s", BaseUrl, ApiAddress), api.Params{
		"province": province,
		"city":     city,
		"detail":   detail,
	})
	if err != nil {
		log.Fatalln("收货地址请求失败", err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil {
		log.Fatalln("收货地址失败", err.Error())
	}
	if resp.Code == 0 {
		log.Printf("收货地址成功：%s%s%s\n", province, city, detail)
	} else {
		log.Fatalf("收货地址失败:[返回] %s", string(session.RespData))
	}

	return 1
}

func Charge(amount int, token string) {
	session := &api.Session{}
	session.DefaultClient()
	session.AddHeader(map[string]string{
		"token": token,
	})
	err := session.PostForJson(fmt.Sprintf("%s%s", BaseUrl, ApiCharge), api.Params{
		"amount": amount,
		"source": 1,
	})
	if err != nil {
		log.Fatalln("充值请求失败", err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil {
		log.Fatalln("充值失败", err.Error())
	}
	if resp.Code == 0 {
		log.Printf("充值成功：%v\n", amount)
	} else {
		log.Fatalf("充值失败:[返回] %s", string(session.RespData))
	}
}

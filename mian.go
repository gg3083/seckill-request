package main

import (
	"encoding/json"
	"fmt"
	"log"
	"seckill-request/api"
	"time"
)

var BaseUrl = "http://127.0.0.1:7087"
var ApiRegister = "/user/register"
var ApiLogin = "/user/login"
var ApiAddress = "/auth/user/address"
var ApiCharge = "/auth/user/charge"

type Response struct {
	Code   int         `json:"code"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

type UserInfo struct {
	PkId     int64  `json:"pk_id"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
	Token      string `json:"token"`
	CreateTime string `json:"create_time"`
}

func main() {

	//注册-->登录 --> 添加收获地址/充值
	//初始化商品
	//查询商品 --秒杀
	userName := "a6"
	passWord := "123456"
	Register(userName, passWord)
	token := Login(userName, passWord)
	log.Printf("当前用户token为:%s\n", token)
	go AddAddress("湖北","武汉", "江夏区",token)
	go Charge(100,token)
	time.Sleep(10* time.Second)
}

func Register(userName, passWord string)  {
	session := &api.Session{}
	session.DefaultClient()
	err := session.PostForJson(fmt.Sprintf("%s%s", BaseUrl, ApiRegister), api.Params{
		"user_name": userName,
		"pass_word": passWord,
	})
	if err != nil  {
		log.Fatalln("注册请求失败",err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil{
		log.Fatalln("注册失败",err.Error())
	}
	if resp.Code == 0 {
		log.Println("注册成功：",userName)
	}else {
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
	if err != nil  {
		log.Fatalln("登录请求失败",err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil{
		log.Fatalln("登录失败",err.Error())
	}
	if resp.Code == 0 {
		log.Println("登录成功：",userName)
	}else {
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
		"city": city,
		"detail": detail,
	})
	if err != nil  {
		log.Fatalln("收货地址请求失败",err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil{
		log.Fatalln("收货地址失败",err.Error())
	}
	if resp.Code == 0 {
		log.Printf("收货地址成功：%s%s%s\n", province, city, detail)
	}else {
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
	if err != nil  {
		log.Fatalln("充值请求失败",err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil{
		log.Fatalln("充值失败",err.Error())
	}
	if resp.Code == 0 {
		log.Printf("充值成功：%v\n", amount)
	}else {
		log.Fatalf("充值失败:[返回] %s", string(session.RespData))
	}
}


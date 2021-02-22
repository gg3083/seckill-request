package service

import (
	"encoding/json"
	"fmt"
	"log"
	"seckill-request/api"
)

func QueryGoods(id int64) string {
	session := &api.Session{}
	session.DefaultClient()

	err := session.Get(fmt.Sprintf("%s%s%v", BaseUrl, ApiGoods, id), api.Params{})
	if err != nil {
		log.Fatalln("获取商品详情失败", err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil {
		log.Fatalln("获取商品详情失败", err.Error())
	}
	if resp.Code == 0 {
		log.Printf("获取商品详情成功：名称为 { %v }\n", resp.Data.(map[string]interface{})["goods_name"])
		res, _ := json.Marshal(resp.Data)
		return string(res)
	} else {
		log.Fatalf("获取商品详情失败:[返回] %s", string(session.RespData))
	}
	return ""
}

func Buy(id, amount int64, num int, token string) {
	session := &api.Session{}
	session.DefaultClient()
	session.AddHeader(map[string]string{
		"token": token,
	})
	err := session.PostForJson(fmt.Sprintf("%s%s", BaseUrl, ApiCharge), api.Params{
		"fk_good_id":  id,
		"price":       amount,
		"num":         num,
		"total_price": amount * int64(num),
		"source":      "秒杀",
		"address": api.Params{
			"pk_id":    "1",
			"province": "湖北",
			"city":     "武汉",
			"detail":   "江夏",
		},
	})
	if err != nil {
		log.Fatalln("购买请求失败", err.Error())
	}
	var resp Response
	if err := json.Unmarshal(session.RespData, &resp); err != nil {
		log.Fatalln("购买失败", err.Error())
	}
	if resp.Code == 0 {
		log.Printf("购买成功：%v\n", amount)
	} else {
		log.Fatalf("购买失败:[返回] %s", string(session.RespData))
	}
}

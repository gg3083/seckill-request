package main

import (
	"encoding/json"
	"log"
	"seckill-request/service"
	"time"
)

type UserInfo struct {
	PkId       int64  `json:"pk_id"`
	UserName   string `json:"user_name"`
	Password   string `json:"password"`
	Token      string `json:"token"`
	CreateTime string `json:"create_time"`
}

type Goods struct {
	PkId        int64  `json:"pk_id"`
	GoodsName   string `json:"goods_name"`
	Price       int64  `json:"price"`
	Stock       int    `json:"stock"`
	SaleNum     int    `json:"sale_num"`
	IsSeckill   int    `json:"is_seckill"`
	SeckillTime int64  `json:"seckill_time"`
}

func main() {

	//注册-->登录 --> 添加收获地址/充值
	userName := "a1"
	passWord := "123456"
	//service.Register(userName, passWord)
	token := service.Login(userName, passWord)
	log.Printf("当前用户token为:%s\n", token)
	//go service.AddAddress("湖北","武汉", "江夏区",token)
	//go service.Charge(100,token)
	//初始化商品
	//查询商品 -- 秒杀
	var goodsId int64 = 1613986640107874600
	goodsJson := service.QueryGoods(goodsId)
	var goods Goods

	if err := json.Unmarshal([]byte(goodsJson), &goods); err != nil {
		log.Fatalf(err.Error())
	}
	log.Println(goodsJson)
	stock := goods.Stock - goods.SaleNum
	if stock <= 0 {
		log.Fatalf("当前商品已被抢购完: %v", stock)
		return
	}
	log.Printf("当前剩余的商品数量为: %v\n", stock)
	//开始下单
	//
	service.Buy(goodsId, goods.Price, 1, token)
	time.Sleep(10 * time.Second)
}

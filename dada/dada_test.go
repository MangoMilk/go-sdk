package dada

import (
	"fmt"
	"testing"
)

var (
	dada *Dada
	// 测试数据
	sourceID = "73753"
	shopNo   = "11047059"
)

func setup() {
	appKey := ""
	appSecret := ""
	dada = NewDada(appKey, appSecret, sourceID, TestEnv)
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func TestGetCity(t *testing.T) {
	ddRes, res,err := dada.GetCity()
	if err!=nil{
		fmt.Println(err)
	}else{
		if ddRes.Code != 0 || ddRes.Status != "success" {
			fmt.Println(ddRes)
		}
		fmt.Println(res)
		fmt.Println(res[0].CityCode, res[0].CityName)
	}
}

func TestAddOrder(t *testing.T) {
	req := AddOrderReq{
		ShopNo:          shopNo,                            // 是	门店编号，门店创建后可在门店列表和单页查看
		OriginID:        "1621506287000",                   // 是	第三方订单ID
		CityCode:        "020",                             // 是	订单所在城市的code（查看各城市对应的code值）
		CargoPrice:      10.00,                             // 是	订单金额（单位：元）
		IsPrepay:        0,                                 // 是	是否需要垫付 1:是 0:否 (垫付订单金额，非运费)
		ReceiverName:    "vincent",                         // 是	收货人姓名
		ReceiverAddress: "广州市天河区花城广场",                      // 是	收货人地址
		ReceiverLat:     23.123641,                         // 是	收货人地址纬度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器）
		ReceiverLng:     113.345769,                        // 是	收货人地址经度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器)
		Callback:        "http://127.0.0.1/v1/dada/notify", // 是	回调URL（查看回调说明）
		CargoWeight:     0.5,                               // 是	订单重量（单位：Kg）
		ReceiverPhone:   "11111111111",                     // 是	收货人手机号（手机号和座机号必填一项）
	}

	ddRes, res,err := dada.AddOrder(&req)
	if err!=nil{
		fmt.Println(err)
	}else {
		if ddRes.Code != 0 || ddRes.Status != "success" {
			fmt.Println(ddRes)
		}

		fmt.Println(res)
	}

}

func TestQueryDeliverFee(t *testing.T) {
	req := QueryDeliverFeeReq{
		ShopNo:          shopNo,          // 是	门店编号，门店创建后可在门店列表和单页查看
		OriginID:        "1621578228000", // 是	第三方订单ID
		CityCode:        "020",           // 是	订单所在城市的code（查看各城市对应的code值）
		CargoPrice:      10.00,           // 是	订单金额（单位：元）
		IsPrepay:        0,               // 是	是否需要垫付 1:是 0:否 (垫付订单金额，非运费)
		ReceiverName:    "vincent",       // 是	收货人姓名
		ReceiverAddress: "广州市天河区花城广场",    // 是	收货人地址
		//ReceiverLat:     23.123641,                         // 是	收货人地址纬度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器）
		//ReceiverLng:     113.345769,                        // 是	收货人地址经度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器)
		Callback:      "http://127.0.0.1/v1/dada/notify", // 是	回调URL（查看回调说明）
		CargoWeight:   0.5,                               // 是	订单重量（单位：Kg）
		ReceiverPhone: "11111111111",                     // 是	收货人手机号（手机号和座机号必填一项）
	}

	ddRes, res,err := dada.QueryDeliverFee(&req)
	if err!=nil{
		fmt.Println(err)
	}else {
		if ddRes.Code != 0 || ddRes.Status != "success" {
			fmt.Println(ddRes)
		}

		fmt.Println(res)
	}
}

func TestQueryOrder(t *testing.T) {
	ddRes, res,err := dada.QueryOrder("1621506287000")
	if err!=nil{
		fmt.Println(err)
	}else {
		if ddRes.Code != 0 || ddRes.Status != "success" {
			fmt.Println(ddRes)
		}
		fmt.Println(res)
	}
}

func TestCancelOrder(t *testing.T) {
	req := CancelOrderReq{
		OrderID:        "1621506287000",
		CancelReasonID: CancelReasonNotNeed,
	}

	ddRes, res,err := dada.CancelOrder(&req)
	if err!=nil{
		fmt.Println(err)
	}else {
		if ddRes.Code != 0 || ddRes.Status != "success" {
			fmt.Println(ddRes)
		}

		fmt.Println(res)
	}
}

func TestAddShop(t *testing.T) {
	//StationName    string  `json:"station_name"`    //	是	门店名称
	//Business       int64   `json:"business"`        //	是	业务类型(食品小吃-1,饮料-2,鲜花绿植-3,文印票务-8,便利店-9,水果生鲜-13,同城电商-19, 医药-20,蛋糕-21,酒品-24,小商品市场-25,服装-26,汽修零配-27,数码家电-28,小龙虾-29,个人-50,火锅-51,个护美妆-53、母婴-55,家居家纺-57,手机-59,家装-61,其他-5)
	//CityName       string  `json:"city_name"`       //	是	城市名称(如,上海)
	//AreaName       string  `json:"area_name"`       //	是	区域名称(如,浦东新区)
	//StationAddress string  `json:"station_address"` //	是	门店地址
	//Lng            float64 `json:"lng"`             //	是	门店经度
	//Lat            float64 `json:"lat"`             //	是	门店纬度
	//ContactName    string  `json:"contact_name"`    //	是	联系人姓名
	//Phone          string  `json:"phone"`           //	是	联系人电话

	//OriginShopID   string  `json:"origin_shop_id"`  //	否	门店编码,可自定义,但必须唯一;若不填写,则系统自动生成
	//IdCard         string  `json:"id_card"`         //	否	联系人身份证
	//Username       string  `json:"username"`        //	否	达达商家app账号(若不需要登陆app,则不用设置)
	//Password       string  `json:"password"`        //	否	达达商家app密码(若不需要登陆app,则不用设置)

	// 19 53 5
	req := AddShopReq{
		shops: []Shop{
			{
				StationName:    "永旺梦乐城（番禺店）",
				Business:       19,
				CityName:       "广州",
				AreaName:       "番禺区",
				StationAddress: "广东省广州市番禺区亚运大道1号",
				Lng:            0,
				Lat:            0,
				ContactName:    "Y先生",
				Phone:          "13412345678",
				OriginShopID:   "10000001",
			},
		},
	}

	ddRes, res,err := dada.AddShop(&req)
	if err!=nil{
		fmt.Println(err)
	}else {
		if ddRes.Code != 0 || ddRes.Status != "success" {
			fmt.Println(ddRes)
		}

		fmt.Println(res)
		//fmt.Println(res.FailedList)
		//fmt.Println(res.SuccessList)
	}
}

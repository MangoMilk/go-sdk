package amap

import (
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/net"
)

/*
 高德地图
*/

const (
	geoUrl = "https://restapi.amap.com/v3/geocode/geo"
)

type Amap struct {
	Key string
}

func NewAmap(key string) *Amap {
	return &Amap{
		Key: key,
	}
}

type getGeoRes struct {
	Status   string `json:"status"`
	Info     string `json:"info"`
	Infocode string `json:"infocode"`
	Count    string `json:"count"`
	Geocodes []geo  `json:"geocodes"`
}

type geo struct {
	FormattedAddress string `json:"formatted_address"` // 结构化地址信息, 省份＋城市＋区县＋城镇＋乡村＋街道＋门牌号码
	Country          string `json:"country"`           // 国家, 默认返回中国
	Province         string `json:"province"`          // 地址所在的省份名, 北京
	City             string `json:"city"`              // 地址所在的城市名, 北京
	Citycode         string `json:"citycode"`          // 城市编码, 010
	District         string `json:"district"`          // 地址所在的区, 朝阳区
	//township: [ ],
	//	neighborhood: {
	//	name: [ ],
	//	type: [ ]
	//},
	//	building: {
	//	name: [ ],
	//	type: [ ]
	//},
	//	street: [ ],// slice or string？
	//	number: [ ],
	Adcode   string `json:"adcode"`   // 区域编码, 110101
	Location string `json:"location"` // 坐标点, 经度，纬度=>"113.321124,23.119252"
	Level    string `json:"level"`    // 匹配级别
}

func (a *Amap) GetGeo(address string, city string) (*getGeoRes,error) {
	query := fmt.Sprintf("?key=%s&address=%s&city=%s", a.Key, address, city)

	data,err := net.HttpGet(geoUrl+query, nil)
	if err!=nil{
		return nil,err
	}

	var amapRes getGeoRes

	if jsonErr:=json.Unmarshal(data,&amapRes);jsonErr!=nil {
		return nil,jsonErr
	}

	return &amapRes,nil
}

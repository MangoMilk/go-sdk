package dada

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/encode"
	"github.com/MangoMilk/go-kit/encrypt"
	"github.com/MangoMilk/go-kit/net"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 在测试环境，使用统一商户和门店进行发单。
// 其中，商户id(source_id)：73753，门店编号：11047059。
const (
	scheme = "http://"

	//测试域名
	testDomain = "newopen.qa.imdada.cn"
	//线上域名
	onlineDomain = "newopen.imdada.cn"

	getCityUrl           = "/api/cityCode/list"
	addOrderUrl          = "/api/order/addOrder"
	reAddOrderUrl        = "/api/order/reAddOrder"
	queryDeliverFeeUrl   = "/api/order/queryDeliverFee"
	addAfterQueryUrl     = "/api/order/addAfterQuery"
	queryOrderUrl        = "/api/order/status/query"
	cancelOrderUrl       = "/api/order/formalCancel"
	confirmOrderGoodsUrl = "/api/order/confirm/goods"
	addMerchantUrl       = "/merchantApi/merchant/add"
	addShopUrl           = "/api/shop/add"
	updateShopUrl        = "/api/shop/update"
	confirmMessageUrl    = "/api/message/confirm"
)

type CancelFrom int64

const (
	CancelFromTransporter = CancelFrom(1)
	CancelFromSupplier    = CancelFrom(2)
	CancelFromSystem      = CancelFrom(3)
)

type Env string

const (
	TestEnv   = Env("test")
	OnlineEnv = Env("online")
)

type Dada struct {
	AppKey     string
	AppSecret  string
	HttpHeader map[string]string
	SourceID   string
	Env        Env
}

func NewDada(appKey string, appSecret string, sourceID string, env Env) *Dada {
	return &Dada{
		AppKey:     appKey,
		AppSecret:  appSecret,
		HttpHeader: map[string]string{"Content-Type": "application/json"},
		SourceID:   sourceID,
		Env:        env,
	}
}

func (dd *Dada) SetHttpHeader(k string, v string) {
	dd.HttpHeader[k] = v
}

func (dd *Dada) genUrl(api string) (url string) {
	if dd.Env == TestEnv {
		url = scheme + testDomain + api
	} else {
		url = scheme + onlineDomain + api
	}

	return
}

func (dd *Dada) genSign(body interface{}, appSecret string) string {
	var data = make(map[string]interface{})
	refVal := reflect.ValueOf(body)
	for i := 0; i < refVal.NumField(); i++ {
		data[refVal.Type().Field(i).Tag.Get("json")] = refVal.Field(i).String()
	}

	var dataKeys []string
	for k, v := range data {
		if v != "" && k != "signature" {
			dataKeys = append(dataKeys, k)
		}
	}
	sort.Strings(dataKeys)

	var signStr = ""
	for _, v := range dataKeys {
		if v != "" {
			val := data[v].(string)
			signStr += (v + val)
		}
	}

	h := md5.New()
	s := appSecret + signStr + appSecret
	signByte := []byte(s)
	h.Write(signByte)

	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

type baseReq struct {
	AppKey    string `json:"app_key"`   // 应用Key，对应开发者账号中的app_key
	Signature string `json:"signature"` // 签名Hash值
	Timestamp string `json:"timestamp"` // 时间戳,单位秒，即unix-timestamp
	Format    string `json:"format"`    // 请求格式，暂时只支持json
	V         string `json:"v"`         // API版本
	SourceID  string `json:"source_id"` // 商户编号（创建商户账号分配的编号）
	Body      string `json:"body"`      // 业务参数，JSON字符串
}

func (dd *Dada) genBaseReq(req interface{}) (*baseReq, error) {
	dataByte, jsonErr := json.Marshal(req)
	if jsonErr != nil {
		return nil, jsonErr
	}

	r := baseReq{
		AppKey:    dd.AppKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Format:    "json",
		V:         "1.0",
		SourceID:  dd.SourceID,
		Body:      string(dataByte),
	}

	r.Signature = dd.genSign(r, dd.AppSecret)

	return &r, nil
}

type baseRes struct {
	Status    string      `json:"status"`    // 响应状态，成功为"success"，失败为"fail"
	Code      int64       `json:"code"`      // 响应返回码，参考接口返回码
	Msg       string      `json:"msg"`       // 响应描述
	Result    interface{} `json:"result"`    // 响应结果，JSON对象，详见具体的接口描述
	ErrorCode int64       `json:"errorCode"` // 错误编码，与code一致
}

// ==================== 城市信息 ====================
type getCityRes struct {
	CityName string `map:"cityName"`
	CityCode string `map:"cityCode"`
}

func (dd *Dada) GetCity() (*baseRes, []*getCityRes, error) {
	ddReq, reqErr := dd.genBaseReq("")
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(getCityUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	result := ddRes.Result.([]interface{})
	arrLen := len(result)
	var res = make([]*getCityRes, arrLen)
	for i := 0; i < arrLen; i++ {
		var unit getCityRes
		encode.Map2struct(result[i], &unit)
		res[i] = &unit
	}

	return &ddRes, res, nil
}

// ==================== 新增订单 ====================
type IsFinishCodeNeeded int64

const (
	IsFinishCodeNeededYes = IsFinishCodeNeeded(1)
	IsFinishCodeNeededNo  = IsFinishCodeNeeded(0)
)

type AddOrderReq struct {
	ShopNo          string  `json:"shop_no"`          // 是	门店编号，门店创建后可在门店列表和单页查看
	OriginID        string  `json:"origin_id"`        // 是	第三方订单ID
	CityCode        string  `json:"city_code"`        // 是	订单所在城市的code（查看各城市对应的code值）
	CargoPrice      float64 `json:"cargo_price"`      // 是	订单金额（单位：元）
	IsPrepay        int64   `json:"is_prepay"`        // 是	是否需要垫付 1:是 0:否 (垫付订单金额，非运费)
	ReceiverName    string  `json:"receiver_name"`    // 是	收货人姓名
	ReceiverAddress string  `json:"receiver_address"` // 是	收货人地址
	ReceiverLat     float64 `json:"receiver_lat"`     // 是	收货人地址纬度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器）
	ReceiverLng     float64 `json:"receiver_lng"`     // 是	收货人地址经度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器)
	Callback        string  `json:"callback"`         // 是	回调URL（查看回调说明）
	CargoWeight     float64 `json:"cargo_weight"`     // 是	订单重量（单位：Kg）
	ReceiverPhone   string  `json:"receiver_phone"`   // 是	收货人手机号（手机号和座机号必填一项）
	ReceiverTel     string  `json:"receiver_tel"`     // 否	收货人座机号（手机号和座机号必填一项）
	Tips            float64 `json:"tips"`             // 否	小费（单位：元，精确小数点后一位）
	Info            string  `json:"info"`             // 否	订单备注
	CargoType       int64   `json:"cargo_type"`       // 否	订单商品类型：食品小吃-1,饮料-2,鲜花绿植-3,文印票务-8,便利店-9,水果生鲜-13,同城电商-19, 医药-20,蛋糕-21,酒品-24,小商品市场-25,服装-26,汽修零配-27,数码家电-28,小龙虾-29,个人-50,火锅-51,个护美妆-53、母婴-55,家居家纺-57,手机-59,家装-61,其他-5
	CargoNum        int64   `json:"cargo_num"`        // 否	订单商品数量
	InvoiceTitle    string  `json:"invoice_title"`    // 否	发票抬头
	OriginMark      string  `json:"origin_mark"`      // 否	订单来源标示（只支持字母，最大长度为10）
	OriginMarkNo    string  `json:"origin_mark_no"`
	/* 否	订单来源编号，最大长度为30，该字段可以显示在骑士APP订单详情页面，示例：
	origin_mark_no:"#京东到家#1"
	达达骑士APP看到的是：#京东到家#1
	*/
	IsUseInsurance int64 `json:"is_use_insurance"`
	/*否
	是否使用保价费（0：不使用保价，1：使用保价； 同时，请确保填写了订单金额（cargo_price））
	商品保价费(当商品出现损坏，可获取一定金额的赔付)
	保费=配送物品实际价值*费率（5‰），配送物品价值及最高赔付不超过10000元， 最高保费为50元（物品价格最小单位为100元，不足100元部分按100元认定，保价费向上取整数， 如：物品声明价值为201元，保价费为300元*5‰=1.5元，取整数为2元。）
	若您选择不保价，若物品出现丢失或损毁，最高可获得平台30元优惠券。 （优惠券直接存入用户账户中）。
	*/
	IsFinishCodeNeeded IsFinishCodeNeeded `json:"is_finish_code_needed"` // 否	收货码（0：不需要；1：需要。收货码的作用是：骑手必须输入收货码才能完成订单妥投）
	DelayPublishTime   int64              `json:"delay_publish_time"`    // 否	预约发单时间（预约时间unix时间戳(10位),精确到分;整分钟为间隔，并且需要至少提前5分钟预约，可以支持未来3天内的订单发预约单。）
	IsDirectDelivery   int64              `json:"is_direct_delivery"`    // 否	是否选择直拿直送（0：不需要；1：需要。选择直拿直送后，同一时间骑士只能配送此订单至完成，同时，也会相应的增加配送费用）
	ProductList        []Product          `json:"product_list"`          // 否	订单商品明细
	PickUpPos          string             `json:"pick_up_pos"`           // 否	货架信息,该字段可在骑士APP订单备注中展示
}

type Product struct {
	SkuName      string  `json:"sku_name"`       // 是	商品名称，限制长度128
	SrcProductNo string  `json:"src_product_no"` // 是	商品编码，限制长度64
	Count        float64 `json:"count"`          // 是	商品数量，精确到小数点后两位
	Unit         string  `json:"unit"`           // 否	商品单位，默认：件
}

type addOrderRes struct {
	Distance     float64 `map:"distance"`     // 是 配送距离(单位：米)
	Fee          float64 `map:"fee"`          // 是 实际运费(单位：元)，运费减去优惠券费用
	DeliverFee   float64 `map:"deliverFee"`   // 是 运费(单位：元)
	CouponFee    float64 `map:"couponFee"`    // 否 优惠券费用(单位：元)
	Tips         float64 `map:"tips"`         // 否 小费(单位：元)
	InsuranceFee float64 `map:"insuranceFee"` // 否 保价费(单位：元)
}

func (dd *Dada) AddOrder(req *AddOrderReq) (*baseRes, *addOrderRes, error) {

	ddReq, reqErr := dd.genBaseReq(req)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(addOrderUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	var res addOrderRes
	encode.Map2struct(ddRes.Result, &res)

	return &ddRes, &res, nil
}

// ==================== 订单重发 ====================
type ReAddOrderReq struct {
	ShopNo          string  `json:"shop_no"`          // 是	门店编号，门店创建后可在门店列表和单页查看
	OriginID        string  `json:"origin_id"`        // 是	第三方订单ID
	CityCode        string  `json:"city_code"`        // 是	订单所在城市的code（查看各城市对应的code值）
	CargoPrice      float64 `json:"cargo_price"`      // 是	订单金额（单位：元）
	IsPrepay        int64   `json:"is_prepay"`        // 是	是否需要垫付 1:是 0:否 (垫付订单金额，非运费)
	ReceiverName    string  `json:"receiver_name"`    // 是	收货人姓名
	ReceiverAddress string  `json:"receiver_address"` // 是	收货人地址
	ReceiverLat     float64 `json:"receiver_lat"`     // 是	收货人地址纬度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器）
	ReceiverLng     float64 `json:"receiver_lng"`     // 是	收货人地址经度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器)
	Callback        string  `json:"callback"`         // 是	回调URL（查看回调说明）
	CargoWeight     float64 `json:"cargo_weight"`     // 是	订单重量（单位：Kg）
	ReceiverPhone   string  `json:"receiver_phone"`   // 是	收货人手机号（手机号和座机号必填一项）
	ReceiverTel     string  `json:"receiver_tel"`     // 否	收货人座机号（手机号和座机号必填一项）
	Tips            float64 `json:"tips"`             // 否	小费（单位：元，精确小数点后一位）
	Info            string  `json:"info"`             // 否	订单备注
	CargoType       int64   `json:"cargo_type"`       // 否	订单商品类型：食品小吃-1,饮料-2,鲜花绿植-3,文印票务-8,便利店-9,水果生鲜-13,同城电商-19, 医药-20,蛋糕-21,酒品-24,小商品市场-25,服装-26,汽修零配-27,数码家电-28,小龙虾-29,个人-50,火锅-51,个护美妆-53、母婴-55,家居家纺-57,手机-59,家装-61,其他-5
	CargoNum        int64   `json:"cargo_num"`        // 否	订单商品数量
	InvoiceTitle    string  `json:"invoice_title"`    // 否	发票抬头
	OriginMark      string  `json:"origin_mark"`      // 否	订单来源标示（只支持字母，最大长度为10）
	OriginMarkNo    string  `json:"origin_mark_no"`
	/* 否	订单来源编号，最大长度为30，该字段可以显示在骑士APP订单详情页面，示例：
	origin_mark_no:"#京东到家#1"
	达达骑士APP看到的是：#京东到家#1
	*/
	IsUseInsurance int64 `json:"is_use_insurance"`
	/*否
	是否使用保价费（0：不使用保价，1：使用保价； 同时，请确保填写了订单金额（cargo_price））
	商品保价费(当商品出现损坏，可获取一定金额的赔付)
	保费=配送物品实际价值*费率（5‰），配送物品价值及最高赔付不超过10000元， 最高保费为50元（物品价格最小单位为100元，不足100元部分按100元认定，保价费向上取整数， 如：物品声明价值为201元，保价费为300元*5‰=1.5元，取整数为2元。）
	若您选择不保价，若物品出现丢失或损毁，最高可获得平台30元优惠券。 （优惠券直接存入用户账户中）。
	*/
	IsFinishCodeNeeded IsFinishCodeNeeded `json:"is_finish_code_needed"` // 否	收货码（0：不需要；1：需要。收货码的作用是：骑手必须输入收货码才能完成订单妥投）
	DelayPublishTime   int64              `json:"delay_publish_time"`    // 否	预约发单时间（预约时间unix时间戳(10位),精确到分;整分钟为间隔，并且需要至少提前5分钟预约，可以支持未来3天内的订单发预约单。）
	IsDirectDelivery   int64              `json:"is_direct_delivery"`    // 否	是否选择直拿直送（0：不需要；1：需要。选择直拿直送后，同一时间骑士只能配送此订单至完成，同时，也会相应的增加配送费用）
	ProductList        []Product          `json:"product_list"`          // 否	订单商品明细
	PickUpPos          string             `json:"pick_up_pos"`           // 否	货架信息,该字段可在骑士APP订单备注中展示
}

type reAddOrderRes struct {
	Distance     float64 `map:"distance"`     // 是 配送距离(单位：米)
	Fee          float64 `map:"fee"`          // 是 实际运费(单位：元)，运费减去优惠券费用
	DeliverFee   float64 `map:"deliverFee"`   // 是 运费(单位：元)
	CouponFee    float64 `map:"couponFee"`    // 否 优惠券费用(单位：元)
	Tips         float64 `map:"tips"`         // 否 小费(单位：元)
	InsuranceFee float64 `map:"insuranceFee"` // 否 保价费(单位：元)
}

func (dd *Dada) ReAddOrder(req *ReAddOrderReq) (*baseRes, *reAddOrderRes, error) {

	ddReq, reqErr := dd.genBaseReq(req)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(reAddOrderUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	var res reAddOrderRes
	encode.Map2struct(ddRes.Result, &res)

	return &ddRes, &res, nil
}

// ==================== 预发布订单 ====================
type QueryDeliverFeeReq struct {
	ShopNo          string  `json:"shop_no"`          // 是	门店编号，门店创建后可在门店列表和单页查看
	OriginID        string  `json:"origin_id"`        // 是	第三方订单ID
	CityCode        string  `json:"city_code"`        // 是	订单所在城市的code（查看各城市对应的code值）
	CargoPrice      float64 `json:"cargo_price"`      // 是	订单金额（单位：元）
	IsPrepay        int64   `json:"is_prepay"`        // 是	是否需要垫付 1:是 0:否 (垫付订单金额，非运费)
	ReceiverName    string  `json:"receiver_name"`    // 是	收货人姓名
	ReceiverAddress string  `json:"receiver_address"` // 是	收货人地址
	ReceiverLat     float64 `json:"receiver_lat"`     // 否	收货人地址纬度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器）
	ReceiverLng     float64 `json:"receiver_lng"`     // 否	收货人地址经度（高德坐标系，若是其他地图经纬度需要转化成高德地图经纬度，高德地图坐标拾取器)
	Callback        string  `json:"callback"`         // 是	回调URL（查看回调说明）
	CargoWeight     float64 `json:"cargo_weight"`     // 是	订单重量（单位：Kg）
	ReceiverPhone   string  `json:"receiver_phone"`   // 是	收货人手机号（手机号和座机号必填一项）
	ReceiverTel     string  `json:"receiver_tel"`     // 否	收货人座机号（手机号和座机号必填一项）
	Tips            float64 `json:"tips"`             // 否	小费（单位：元，精确小数点后一位）
	Info            string  `json:"info"`             // 否	订单备注
	CargoType       int64   `json:"cargo_type"`       // 否	订单商品类型：食品小吃-1,饮料-2,鲜花绿植-3,文印票务-8,便利店-9,水果生鲜-13,同城电商-19, 医药-20,蛋糕-21,酒品-24,小商品市场-25,服装-26,汽修零配-27,数码家电-28,小龙虾-29,个人-50,火锅-51,个护美妆-53、母婴-55,家居家纺-57,手机-59,家装-61,其他-5
	CargoNum        int64   `json:"cargo_num"`        // 否	订单商品数量
	InvoiceTitle    string  `json:"invoice_title"`    // 否	发票抬头
	OriginMark      string  `json:"origin_mark"`      // 否	订单来源标示（只支持字母，最大长度为10）
	OriginMarkNo    string  `json:"origin_mark_no"`
	/* 否	订单来源编号，最大长度为30，该字段可以显示在骑士APP订单详情页面，示例：
	origin_mark_no:"#京东到家#1"
	达达骑士APP看到的是：#京东到家#1
	*/
	IsUseInsurance int64 `json:"is_use_insurance"`
	/*否
	是否使用保价费（0：不使用保价，1：使用保价； 同时，请确保填写了订单金额（cargo_price））
	商品保价费(当商品出现损坏，可获取一定金额的赔付)
	保费=配送物品实际价值*费率（5‰），配送物品价值及最高赔付不超过10000元， 最高保费为50元（物品价格最小单位为100元，不足100元部分按100元认定，保价费向上取整数， 如：物品声明价值为201元，保价费为300元*5‰=1.5元，取整数为2元。）
	若您选择不保价，若物品出现丢失或损毁，最高可获得平台30元优惠券。 （优惠券直接存入用户账户中）。
	*/
	IsFinishCodeNeeded IsFinishCodeNeeded `json:"is_finish_code_needed"` // 否	收货码（0：不需要；1：需要。收货码的作用是：骑手必须输入收货码才能完成订单妥投）
	DelayPublishTime   int64              `json:"delay_publish_time"`    // 否	预约发单时间（预约时间unix时间戳(10位),精确到分;整分钟为间隔，并且需要至少提前5分钟预约，可以支持未来3天内的订单发预约单。）
	IsDirectDelivery   int64              `json:"is_direct_delivery"`    // 否	是否选择直拿直送（0：不需要；1：需要。选择直拿直送后，同一时间骑士只能配送此订单至完成，同时，也会相应的增加配送费用）
	ProductList        []Product          `json:"product_list"`          // 否	订单商品明细
	PickUpPos          string             `json:"pick_up_pos"`
}

type queryDeliverFeeRes struct {
	Distance     float64 `map:"distance"`     // 是 配送距离(单位：米)
	Fee          float64 `map:"fee"`          // 是 实际运费(单位：元)，运费减去优惠券费用
	DeliverFee   float64 `map:"deliverFee"`   // 是 运费(单位：元)
	CouponFee    float64 `map:"couponFee"`    // 否 优惠券费用(单位：元)
	Tips         float64 `map:"tips"`         // 否 小费(单位：元)
	InsuranceFee float64 `map:"insuranceFee"` // 否 保价费(单位：元)
	DeliveryNo   string  `map:"deliveryNo"`   // 是	平台订单号，有效期3分钟
}

func (dd *Dada) QueryDeliverFee(req *QueryDeliverFeeReq) (*baseRes, *queryDeliverFeeRes, error) {
	ddReq, reqErr := dd.genBaseReq(req)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(queryDeliverFeeUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	var res queryDeliverFeeRes
	encode.Map2struct(ddRes.Result, &res)

	return &ddRes, &res, nil
}

// ==================== 预发布订单后直接下单 ====================
type AddAfterQueryReq struct {
	DeliveryNo string `json:"deliveryNo"` // 是	平台订单号
}

func (dd *Dada) AddAfterQuery(deliveryNo string) (*baseRes, error) {
	ddReq, reqErr := dd.genBaseReq(&AddAfterQueryReq{DeliveryNo: deliveryNo})
	if reqErr != nil {
		return nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(addAfterQueryUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, reqErr
	}

	return &ddRes, nil
}

// ==================== 订单详情查询（一分钟更新一次） ====================
type OrderStatusCode float64

const (
	OrderStatusCodeWaitToReceive     = OrderStatusCode(1)
	OrderStatusCodeWaitToTake        = OrderStatusCode(2)
	OrderStatusCodeDelivering        = OrderStatusCode(3)
	OrderStatusCodeDone              = OrderStatusCode(4)
	OrderStatusCodeCanceled          = OrderStatusCode(5)
	OrderStatusCodeAssign            = OrderStatusCode(8)
	OrderStatusCodeUnusualBacking    = OrderStatusCode(9)
	OrderStatusCodeUnusualBackDone   = OrderStatusCode(10)
	OrderStatusCodeTransporterArrive = OrderStatusCode(100)
	OrderStatusCodeAddOrderFail      = OrderStatusCode(1000)
)

type QueryOrderReq struct {
	OrderID string `json:"order_id"` // 是	第三方订单ID
}

type queryOrderRes struct {
	OrderId          string  `map:"orderId"`          //	第三方订单编号
	StatusCode       float64 `map:"statusCode"`       // 	订单状态(待接单＝1,待取货＝2,配送中＝3,已完成＝4,已取消＝5, 指派单=8,妥投异常之物品返回中=9, 妥投异常之物品返回完成=10, 骑士到店=100,创建达达运单失败=1000 可参考文末的状态说明）
	StatusMsg        string  `map:"statusMsg"`        // 	订单状态
	TransporterName  string  `map:"transporterName"`  // 	骑手姓名
	TransporterPhone string  `map:"transporterPhone"` //	骑手电话
	TransporterLng   string  `map:"transporterLng"`   //	骑手经度
	TransporterLat   string  `map:"transporterLat"`   //	骑手纬度
	DeliveryFee      float64 `map:"deliveryFee"`      //	配送费,单位为元
	Tips             float64 `map:"tips"`             //	小费,单位为元
	CouponFee        float64 `map:"couponFee"`        //	优惠券费用,单位为元
	InsuranceFee     float64 `map:"insuranceFee"`     //	保价费,单位为元
	ActualFee        float64 `map:"actualFee"`        //	实际支付费用,单位为元
	Distance         float64 `map:"distance"`         //	配送距离,单位为米
	CreateTime       string  `map:"createTime"`       //	发单时间
	AcceptTime       string  `map:"acceptTime"`       //	接单时间,若未接单,则为空
	FetchTime        string  `map:"fetchTime"`        //	取货时间,若未取货,则为空
	FinishTime       string  `map:"finishTime"`       //	送达时间,若未送达,则为空
	CancelTime       string  `map:"cancelTime"`       //	取消时间,若未取消,则为空
	OrderFinishCode  string  `map:"orderFinishCode"`  //	收货码
	DeductFee        float64 `map:"deductFee"`        //	违约金
	ReceiptUrl       string  `map:"receiptUrl"`       //
	SupplierName     string  `map:"supplierName"`     //	店铺名
	SupplierAddress  string  `map:"supplierAddress"`  //	店铺地址
	SupplierPhone    string  `map:"supplierPhone"`    //	店铺联系手机
	SupplierLat      string  `map:"supplierLat"`      //	店铺纬度
	SupplierLng      string  `map:"supplierLng"`      //	店铺经度
}

func (dd *Dada) QueryOrder(orderID string) (*baseRes, *queryOrderRes, error) {
	ddReq, reqErr := dd.genBaseReq(&QueryOrderReq{OrderID: orderID})
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(queryOrderUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	var res queryOrderRes
	encode.Map2struct(ddRes.Result, &res)

	return &ddRes, &res, nil
}

// ==================== 取消订单 ====================
type CancelReasonID int64

const (
	CancelReasonNoTransporterReceive = CancelReasonID(1)     //没有配送员接单
	CancelReasonNoTransporterTake    = CancelReasonID(2)     //配送员没来取货
	CancelReasonTransporterBadManner = CancelReasonID(3)     //配送员态度太差
	CancelReasonUserCancel           = CancelReasonID(4)     //顾客取消订单
	CancelReasonWrongOrderInfo       = CancelReasonID(5)     //订单填写错误
	CancelReasonTransporterCancel    = CancelReasonID(34)    //配送员让我取消此单
	CancelReasonTransporterUnwilling = CancelReasonID(35)    //配送员不愿上门取货
	CancelReasonNotNeed              = CancelReasonID(36)    //我不需要配送了
	CancelReasonTransporterNotDone   = CancelReasonID(37)    //配送员以各种理由表示无法完成订单
	CancelReasonOther                = CancelReasonID(10000) //其他
)

type CancelOrderReq struct {
	OrderID        string         `json:"order_id"`         // 是	第三方订单ID
	CancelReasonID CancelReasonID `json:"cancel_reason_id"` // 是	取消原因ID
	CancelReason   string         `json:"cancel_reason"`    // 是	取消原因(当取消原因ID为其他时，此字段必填)
}

type cancelOrderRes struct {
	DeductFee float64 `map:"deduct_fee"` //	违约金
}

func (dd *Dada) CancelOrder(req *CancelOrderReq) (*baseRes, *cancelOrderRes, error) {
	ddReq, reqErr := dd.genBaseReq(req)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(cancelOrderUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	var res cancelOrderRes
	encode.Map2struct(ddRes.Result, &res)

	return &ddRes, &res, nil
}

// ==================== 注册商户 ====================
type AddMerchantReq struct {
	Mobile            string `json:"mobile"`             // 是	注册商户手机号,用于登陆商户后台
	CityName          string `json:"city_name"`          // 是	商户城市名称(如,上海)
	EnterpriseName    string `json:"enterprise_name"`    // 是	企业全称
	EnterpriseAddress string `json:"enterprise_address"` // 是	企业地址
	ContactName       string `json:"contact_name"`       // 是	联系人姓名
	ContactPhone      string `json:"contact_phone"`      // 是	联系人电话
	Email             string `json:"email"`              // 是	邮箱地址
}

func (dd *Dada) AddMerchant(req *AddMerchantReq) (*baseRes, string, error) {
	ddReq, reqErr := dd.genBaseReq(req)
	if reqErr != nil {
		return nil, "", reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(addMerchantUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, "", reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, "", reqErr
	}

	return &ddRes, strconv.FormatInt(ddRes.Result.(int64), 10), nil
}

// ==================== 创建门店 ====================
type AddShopReq struct {
	Shops []Shop
}

type Shop struct {
	StationName    string  `json:"station_name"`    //	是	门店名称
	Business       int64   `json:"business"`        //	是	业务类型(食品小吃-1,饮料-2,鲜花绿植-3,文印票务-8,便利店-9,水果生鲜-13,同城电商-19, 医药-20,蛋糕-21,酒品-24,小商品市场-25,服装-26,汽修零配-27,数码家电-28,小龙虾-29,个人-50,火锅-51,个护美妆-53、母婴-55,家居家纺-57,手机-59,家装-61,其他-5)
	CityName       string  `json:"city_name"`       //	是	城市名称(如,上海)
	AreaName       string  `json:"area_name"`       //	是	区域名称(如,浦东新区)
	StationAddress string  `json:"station_address"` //	是	门店地址
	Lng            float64 `json:"lng"`             //	是	门店经度
	Lat            float64 `json:"lat"`             //	是	门店纬度
	ContactName    string  `json:"contact_name"`    //	是	联系人姓名
	Phone          string  `json:"phone"`           //	是	联系人电话
	OriginShopID   string  `json:"origin_shop_id"`  //	否	门店编码,可自定义,但必须唯一;若不填写,则系统自动生成
	IdCard         string  `json:"id_card"`         //	否	联系人身份证
	Username       string  `json:"username"`        //	否	达达商家app账号(若不需要登陆app,则不用设置)
	Password       string  `json:"password"`        //	否	达达商家app密码(若不需要登陆app,则不用设置)
}

type addShopRes struct {
	Success     float64       `map:"success"`
	SuccessList []successItem `map:"successList"`
	FailedList  []failedItem  `map:"failedList"`
}

type successItem struct {
	Phone          string  `map:"phone"`
	Business       float64 `map:"business"`
	Lng            float64 `map:"lng"`
	Lat            float64 `map:"lat"`
	StationName    string  `map:"stationName"`
	OriginShopId   string  `map:"originShopId"`
	ContactName    string  `map:"contactName"`
	StationAddress string  `map:"stationAddress"`
	CityName       string  `map:"cityName"`
	AreaName       string  `map:"areaName"`
}

type failedItem struct {
	ShopNo   string `map:"shopNo"`
	Msg      string `map:"msg"`
	ShopName string `map:"shopName"`
}

func (dd *Dada) AddShop(req *AddShopReq) (*baseRes, *addShopRes, error) {
	ddReq, reqErr := dd.genBaseReq(req.Shops)
	if reqErr != nil {
		return nil, nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(addShopUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, nil, reqErr
	}

	var res addShopRes

	if ddRes.Result != nil {
		res.Success = ddRes.Result.(map[string]interface{})["success"].(float64)

		successList := ddRes.Result.(map[string]interface{})["successList"]
		if successList != nil {
			successRes := make([]successItem, len(successList.([]interface{})))
			for i := 0; i < len(successList.([]interface{})); i++ {
				encode.Map2struct(successList.([]interface{})[i], &successRes[i])
			}
			res.SuccessList = successRes
		}

		failedList := ddRes.Result.(map[string]interface{})["failedList"]
		if failedList != nil {
			failRes := make([]failedItem, len(failedList.([]interface{})))
			for i := 0; i < len(failedList.([]interface{})); i++ {
				encode.Map2struct(failedList.([]interface{})[i], &failRes[i])
			}
			res.FailedList = failRes
		}
	}

	return &ddRes, &res, nil
}

// ==================== 更新门店 ====================
type UpdateShopReq struct {
	OriginShopID   string  `json:"origin_shop_id"`  //	是	门店编码
	NewShopID      string  `json:"new_shop_id"`     //	否	新的门店编码
	StationName    string  `json:"station_name"`    //	否	门店名称
	Business       int64   `json:"business"`        //	否	业务类型(食品小吃-1,饮料-2,鲜花绿植-3,文印票务-8,便利店-9,水果生鲜-13,同城电商-19, 医药-20,蛋糕-21,酒品-24,小商品市场-25,服装-26,汽修零配-27,数码家电-28,小龙虾-29,个人-50,火锅-51,个护美妆-53、母婴-55,家居家纺-57,手机-59,家装-61,其他-5)
	CityName       string  `json:"city_name"`       //	否	城市名称(如,上海)
	AreaName       string  `json:"area_name"`       //	否	区域名称(如,浦东新区)
	StationAddress string  `json:"station_address"` //	否	门店地址
	Lng            float64 `json:"lng"`             //	否	门店经度
	Lat            float64 `json:"lat"`             //	否	门店纬度
	ContactName    string  `json:"contact_name"`    //	否	联系人姓名
	Phone          string  `json:"phone"`           //	否	联系人电话
	Status         int64   `json:"status"`          //	否	门店状态（1-门店激活，0-门店下线）
}

func (dd *Dada) UpdateShop(req *UpdateShopReq) (*baseRes, error) {
	ddReq, reqErr := dd.genBaseReq(req)
	if reqErr != nil {
		return nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(updateShopUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, reqErr
	}

	return &ddRes, nil
}

// ==================== 订单状态通知 ====================
type OrderNotifyReq struct {
	ClientID     string          `json:"client_id" validate:"required"`    //	是	返回达达运单号，默认为空
	OrderID      string          `json:"order_id" validate:"required"`     //	是	添加订单接口中的origin_id值
	OrderStatus  OrderStatusCode `json:"order_status" validate:"required"` //	是	订单状态(待接单＝1,待取货＝2,配送中＝3,已完成＝4,已取消＝5, 指派单=8,妥投异常之物品返回中=9, 妥投异常之物品返回完成=10, 骑士到店=100,创建达达运单失败=1000 可参考文末的状态说明）
	CancelReason string          `json:"cancel_reason"`                    //	是	订单取消原因,其他状态下默认值为空字符串
	CancelFrom   CancelFrom      `json:"cancel_from"`                      //	是	订单取消原因来源(1:达达配送员取消；2:商家主动取消；3:系统或客服取消；0:默认值)
	UpdateTime   string          `json:"update_time" validate:"required"`  //	是	更新时间，时间戳除了创建达达运单失败=1000的精确毫秒，其他时间戳精确到秒
	Signature    string          `json:"signature" validate:"required"`    //	是	对client_id, order_id, update_time的值进行字符串升序排列，再连接字符串，取md5值
	DmID         int64           `json:"dm_id"`                            //	否	达达配送员id，接单以后会传
	DmName       string          `json:"dm_name"`                          //	否	配送员姓名，接单以后会传
	DmMobile     string          `json:"dm_mobile"`                        //	否	配送员手机号，接单以后会传
	FinishCode   string          `json:"finish_code"`                      //	否	收货码
}

func CheckNotifySign(clientID, orderID, updateTime, sign string) error {
	args := []string{clientID, orderID, updateTime}
	sort.Strings(args)
	str := strings.Join(args, "")

	md5Res, err := encrypt.MD5(str)
	if err != nil {
		return err
	}

	if md5Res != sign {
		return fmt.Errorf("check signature fail")
	}

	return nil
}

// ==================== 达达消息通知 ====================
type NotifyMessageType int64

const (
	NotifyMessageTypeTransporterCancel = NotifyMessageType(1)
)

type NotifyReq struct {
	MessageBody string            `json:"messageBody"`
	MessageType NotifyMessageType `json:"messageType"`
}

type NotifyReturnMsg string

const (
	NotifySuccessReturnMsg = NotifyReturnMsg("ok")
	NotifyFailReturnMsg    = NotifyReturnMsg("fail")
)

type NotifyRes struct {
	Status NotifyReturnMsg `json:"status"`
}

// ==================== 骑手取消通知 ====================
type TransporterNotifyReq struct {
	OrderId      string `json:"orderId" validate:"required"` // 是 商家第三方订单号
	DadaOrderId  int64  `json:"dadaOrderId"`                 // 否 达达订单号
	CancelReason string `json:"cancelReason"`                // 是 骑士取消原因
}

// ==================== 通知确认 ====================
type IsConfirmTransporterCancel int64

const (
	IsConfirmTransporterCancelYes = IsConfirmTransporterCancel(1)
	IsConfirmTransporterCancelNo  = IsConfirmTransporterCancel(0)
)

type NotifyConfirmReq struct {
	OrderId     string                     `json:"orderId"`     // 是 商家第三方订单号
	DadaOrderId int64                      `json:"dadaOrderId"` // 否 达达订单号
	IsConfirm   IsConfirmTransporterCancel `json:"isConfirm"`   // 是 0:不同意，1:表示同意
}

func (dd *Dada) ConfirmMessage(req *NotifyConfirmReq) (*baseRes, error) {
	reqByte, jsonMarshalErr := json.Marshal(req)
	if jsonMarshalErr != nil {
		return nil, jsonMarshalErr
	}

	messageReq := &NotifyReq{
		MessageBody: string(reqByte),
		MessageType: NotifyMessageTypeTransporterCancel,
	}

	ddReq, reqErr := dd.genBaseReq(messageReq)
	if reqErr != nil {
		return nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(confirmMessageUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, reqErr
	}

	return &ddRes, nil
}

// ==================== 确认妥投异常之物品返回完成 ====================
type ConfirmOrderGoodsReq struct {
	OrderID string `json:"order_id"` // 是	第三方订单ID
}

func (dd *Dada) ConfirmOrderGoods(orderID string) (*baseRes, error) {
	ddReq, reqErr := dd.genBaseReq(&ConfirmOrderGoodsReq{OrderID: orderID})
	if reqErr != nil {
		return nil, reqErr
	}

	data, httpErr := net.HttpPost(dd.genUrl(confirmOrderGoodsUrl), *ddReq, dd.HttpHeader, false, "", "")
	if httpErr != nil {
		return nil, reqErr
	}

	var ddRes baseRes
	if jsonErr := json.Unmarshal(data, &ddRes); jsonErr != nil {
		return nil, reqErr
	}

	return &ddRes, nil
}

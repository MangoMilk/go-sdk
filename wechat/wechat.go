package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"github.com/MangoMilk/go-kit/net"
	"reflect"
	"sort"
	"strings"
)

const (
	unifiedOrderUrl = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	jsCode2SessionUrl = "https://api.weixin.qq.com/sns/jscode2session"
)

type Wechat struct {
	AppID string
	AppSecret string
}

func NewWechat(appID,appSecret string) *Wechat {
	return &Wechat{
		AppID: appID,
		AppSecret: appSecret,
	}
}

func (wx *Wechat) GenSign(body interface{},apiKey string) string {
	var data = make(map[string]interface{})
	refVal := reflect.ValueOf(body)
	for i := 0; i < refVal.NumField(); i++ {
		xmlKey:=refVal.Type().Field(i).Tag.Get("xml")
		switch xmlKey {
		case "xml":
			break
		case "appid":
			data[xmlKey] = wx.AppID
			break
		default:
			data[xmlKey] = refVal.Field(i).String()
		}
	}

	// 生成签名（sign）和请求参数（xml）
	var dataKeys []string
	for k, v := range data {
		if v != "" {
			dataKeys = append(dataKeys, k)
		}
	}
	sort.Strings(dataKeys)

	// joint data
	signStr  := ""
	for _, key := range dataKeys {
		signStr += key + "=" + data[key].(string) + "&"
	}

	h := md5.New()
	signByte := []byte(strings.Trim(signStr, "&") + "&key=" + apiKey)
	h.Write(signByte)

	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

// ==================== 统一下单 ====================
type UnifiedOrderReq struct {
	XMLName  	xml.Name `xml:"xml"`
	AppID          string `xml:"appid"`  // 是，微信分配的小程序ID
	MchID          string `xml:"mch_id"` // 是，微信分配的小程序ID
	DeviceInfo     string `xml:"device_info"`// 否，自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传"WEB"
	NonceStr       string `xml:"nonce_str"` // 是，随机字符串，长度要求在32位以内。推荐随机数生成算法
	Sign           string `xml:"sign"`      // 是，通过签名算法计算得出的签名值，详见签名生成算法
	SignType       string `xml:"sign_type"`// 否，通过签名算法计算得出的签名值，详见签名生成算法
	Body           string `xml:"body"` // 是，通过签名算法计算得出的签名值，详见签名生成算法
	Detail         string `xml:"detail"`// 否，商品详细描述，对于使用单品优惠的商户，该字段必须按照规范上传，详见“单品优惠参数说明”
	Attach         string `xml:"attach"` // 否，附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	OutTradeNo     string `xml:"out_trade_no"` // 是，商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*且在同一个商户号下唯一。详见商户订单号
	FeeType        string `xml:"fee_type"` // 否，符合ISO 4217标准的三位字母代码，默认人民币：CNY，详细列表请参见货币类型
	TotalFee       string `xml:"total_fee"`        // 是，订单总金额，单位为分，详见支付金额
	SpbillCreateIP string `xml:"spbill_create_ip"` // 是，支持IPV4和IPV6两种格式的IP地址。调用微信支付API的机器IP
	TimeStart      string `xml:"time_start"`       // 否，订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire     string `xml:"time_expire"`
	/* 否，订单失效时间，
		格式为yyyyMMddHHmmss，
		如2009年12月27日9点10分10秒表示为20091227091010。
		订单失效时间是针对订单号而言的，
		由于在请求支付的时候有一个必传参数prepay_id只有两小时的有效期，
		所以在重入时间超过2小时的时候需要重新请求下单接口获取新的prepay_id。
		其他详见时间规则 建议：最短失效时间间隔大于1分钟
	 */
	GoodsTag       string `xml:"goods_tag"` // 否，订单优惠标记，使用代金券或立减优惠功能时需要的参数，说明详见代金券或立减优惠
	NotifyUrl      string `xml:"notify_url"` // 是，异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。公网域名必须为https，如果是走专线接入，使用专线NAT IP或者私有回调域名可使用http。
	TradeType      string `xml:"trade_type"` // 是，小程序取值如下：JSAPI，详细说明见参数规定
	ProductID      string `xml:"product_id"` // 否，trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	LimitPay       string `xml:"limit_pay"` // 否，上传此参数no_credit--可限制用户不能使用信用卡支付
	OpenID         string `xml:"openid"` // 否，trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。openid如何获取，可参考【获取openid】。
	Receipt        string `xml:"receipt"` // 否，Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	ProfitSharing  string `xml:"profit_sharing"`
	/* 否，Y-是，需要分账
		N-否，不分账
		字母要求大写，不传默认不分账
	 */
	SceneInfo      string `xml:"scene_info"`
	/* 否，该字段常用于线下活动时的场景信息上报，
		支持上报实际门店信息，商户也可以按需求自己上报相关信息。
		该字段为JSON对象数据，
		对象格式为
			{
				"store_info": {
					"id": "门店ID",
					"name": "名称",
					"area_code": "编码",
					"address": "地址"
				}
			} ，
		字段详细说明请点击行前的+展开
	 */
}

type unifiedOrderRes struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	DeviceInfo string `xml:"device_info"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	TradeType  string `xml:"trade_type"`
	PrepayID   string `xml:"prepay_id"`
	CodeUrl    string `xml:"code_url"`
}

func (wx *Wechat)UnifiedOrder(req * UnifiedOrderReq) (*unifiedOrderRes,error){
	req.AppID = wx.AppID

	xmlParamsByte,xmlErr:=xml.Marshal(req)
	if xmlErr != nil {
		return nil,xmlErr
	}

	res ,httpErr:= net.HttpPost(unifiedOrderUrl, string(xmlParamsByte), nil, false,"","")
	if httpErr != nil {
		return nil,httpErr
	}

	var data unifiedOrderRes
	if xmlErr := xml.Unmarshal(res, &data);xmlErr != nil {
		return nil,xmlErr
	}

	return &data,nil
}

// ==================== 授权 ====================
type jsCode2SessionRes struct {
	ErrCode    float64 `json:"errcode"`
	ErrMsg     string  `json:"errmsg"`
	SessionKey string  `json:"session_key"`
	OpenID     string  `json:"openid"`
	UnionID    string  `json:"unionid"`
}



func (wx *Wechat)JsCode2Session(code string) (*jsCode2SessionRes,error) {
	queryParam := "?appid=" + wx.AppID + "&secret=" + wx.AppSecret + "&js_code=" + code + "&grant_type=authorization_code"
	res,httpErr:=net.HttpGet(jsCode2SessionUrl+queryParam, nil)
	if httpErr != nil {
		return nil,httpErr
	}

	var data jsCode2SessionRes

	if jsonErr := json.Unmarshal(res, &data);jsonErr != nil {
		return nil,jsonErr
	}

	return &data,nil
}

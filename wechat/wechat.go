package wechat

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/MangoMilk/go-kit/encode"
	"github.com/MangoMilk/go-kit/encrypt"
	"github.com/MangoMilk/go-kit/net"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

const (
	unifiedOrderUrl   = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	jsCode2SessionUrl = "https://api.weixin.qq.com/sns/jscode2session"
	refundUrl         = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	accessTokenUrl    = "https://api.weixin.qq.com/cgi-bin/token"
	wxACodeUnLimitUrl = "https://api.weixin.qq.com/wxa/getwxacodeunlimit"
)

type Wechat struct {
	AppID        string
	AppSecret    string
	ZeroValueMap map[string]interface{} // use for gen sign
}

func NewWechat(appID, appSecret string) *Wechat {
	return &Wechat{
		AppID:        appID,
		AppSecret:    appSecret,
		ZeroValueMap: make(map[string]interface{}),
	}
}

// ==================== 零值处理 ====================
func (wx *Wechat) ZeroValueProcess(res []byte) {
	wx.ZeroValueMap = make(map[string]interface{})
	re, _ := regexp.Compile(`\<(.+)\>0`)
	matchRes := re.FindAllStringSubmatch(string(res), -1)
	if len(matchRes) > 0 {
		for _, match := range matchRes {
			wx.ZeroValueMap[match[1]] = 1
		}
	}
}

// ==================== 签名 ====================
type SignType string

const (
	SignTypeMD5        = SignType("MD5")
	SignTypeHmacSha256 = SignType("HMAC-SHA256")
)

func (wx *Wechat) GenSign(body interface{}, apiKey string) string {
	var data = make(map[string]reflect.Value)
	refVal := reflect.ValueOf(body)
	for i := 0; i < refVal.NumField(); i++ {
		xmlKey := refVal.Type().Field(i).Tag.Get("xml")
		switch xmlKey {
		case "xml", "sign":
			break
		case "appid":
			data[xmlKey] = reflect.ValueOf(wx.AppID)
			break
		default:
			data[xmlKey] = refVal.Field(i)
		}
	}

	// 生成签名（sign）和请求参数（xml）
	var dataKeys []string
	for k, v := range data {
		_, inZeroValueMap := wx.ZeroValueMap[k]
		if !v.IsZero() || inZeroValueMap {
			dataKeys = append(dataKeys, k)
		}
	}

	sort.Strings(dataKeys)

	// joint data
	signStr := ""
	for _, key := range dataKeys {
		signStr += fmt.Sprintf("%v=%v&", key, data[key])
	}

	h := md5.New()
	signByte := []byte(strings.Trim(signStr, "&") + "&key=" + apiKey)
	h.Write(signByte)

	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

type PaySignData struct {
	AppID     string   `xml:"appId"`
	TimeStamp string   `xml:"timeStamp"`
	NonceStr  string   `xml:"nonceStr"`
	Package   string   `xml:"package"`
	SignType  SignType `xml:"signType"`
}

func (wx *Wechat) GenPaySignPackage(prepayID string) string {
	return "prepay_id=" + prepayID
}

// ==================== 统一下单 ====================
type TradeType string

const (
	TradeTypeJsapi  = TradeType("JSAPI")  // JSAPI/小程序
	TradeTypeNative = TradeType("NATIVE") // Native
	TradeTypeApp    = TradeType("APP")    // app
	TradeTypeMWeb   = TradeType("MWEB")   // H5
)

type UnifiedOrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppID          string   `xml:"appid"`            // 是，微信分配的小程序ID
	MchID          string   `xml:"mch_id"`           // 是，微信支付分配的商户号
	DeviceInfo     string   `xml:"device_info"`      // 否，自定义参数，可以为终端设备号(门店号或收银设备ID)，PC网页或公众号内支付可以传"WEB"
	NonceStr       string   `xml:"nonce_str"`        // 是，随机字符串，长度要求在32位以内。推荐随机数生成算法
	Sign           string   `xml:"sign"`             // 是，通过签名算法计算得出的签名值，详见签名生成算法
	SignType       SignType `xml:"sign_type"`        // 否，通过签名算法计算得出的签名值，详见签名生成算法
	Body           string   `xml:"body"`             // 是，通过签名算法计算得出的签名值，详见签名生成算法
	Detail         string   `xml:"detail"`           // 否，商品详细描述，对于使用单品优惠的商户，该字段必须按照规范上传，详见“单品优惠参数说明”
	Attach         string   `xml:"attach"`           // 否，附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。
	OutTradeNo     string   `xml:"out_trade_no"`     // 是，商户系统内部订单号，要求32个字符内，只能是数字、大小写字母_-|*且在同一个商户号下唯一。详见商户订单号
	FeeType        string   `xml:"fee_type"`         // 否，符合ISO 4217标准的三位字母代码，默认人民币：CNY，详细列表请参见货币类型
	TotalFee       string   `xml:"total_fee"`        // 是，订单总金额，单位为分，详见支付金额
	SpbillCreateIP string   `xml:"spbill_create_ip"` // 是，支持IPV4和IPV6两种格式的IP地址。调用微信支付API的机器IP
	TimeStart      string   `xml:"time_start"`       // 否，订单生成时间，格式为yyyyMMddHHmmss，如2009年12月25日9点10分10秒表示为20091225091010。其他详见时间规则
	TimeExpire     string   `xml:"time_expire"`
	/* 否，订单失效时间，
	格式为yyyyMMddHHmmss，
	如2009年12月27日9点10分10秒表示为20091227091010。
	订单失效时间是针对订单号而言的，
	由于在请求支付的时候有一个必传参数prepay_id只有两小时的有效期，
	所以在重入时间超过2小时的时候需要重新请求下单接口获取新的prepay_id。
	其他详见时间规则 建议：最短失效时间间隔大于1分钟
	*/
	GoodsTag      string    `xml:"goods_tag"`  // 否，订单优惠标记，使用代金券或立减优惠功能时需要的参数，说明详见代金券或立减优惠
	NotifyUrl     string    `xml:"notify_url"` // 是，异步接收微信支付结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。公网域名必须为https，如果是走专线接入，使用专线NAT IP或者私有回调域名可使用http。
	TradeType     TradeType `xml:"trade_type"` // 是，小程序取值如下：JSAPI，详细说明见参数规定
	ProductID     string    `xml:"product_id"` // 否，trade_type=NATIVE时，此参数必传。此参数为二维码中包含的商品ID，商户自行定义。
	LimitPay      string    `xml:"limit_pay"`  // 否，上传此参数no_credit--可限制用户不能使用信用卡支付
	OpenID        string    `xml:"openid"`     // 否，trade_type=JSAPI，此参数必传，用户在商户appid下的唯一标识。openid如何获取，可参考【获取openid】。
	Receipt       string    `xml:"receipt"`    // 否，Y，传入Y时，支付成功消息和支付详情页将出现开票入口。需要在微信支付商户平台或微信公众平台开通电子发票功能，传此字段才可生效
	ProfitSharing string    `xml:"profit_sharing"`
	/* 否，Y-是，需要分账
	N-否，不分账
	字母要求大写，不传默认不分账
	*/
	SceneInfo string `xml:"scene_info"`
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
	ReturnCode string    `xml:"return_code"`
	ReturnMsg  string    `xml:"return_msg"`
	AppID      string    `xml:"appid"`
	MchID      string    `xml:"mch_id"`
	DeviceInfo string    `xml:"device_info"`
	NonceStr   string    `xml:"nonce_str"`
	Sign       string    `xml:"sign"`
	ResultCode string    `xml:"result_code"`
	ErrCode    string    `xml:"err_code"`
	ErrCodeDes string    `xml:"err_code_des"`
	TradeType  TradeType `xml:"trade_type"`
	PrepayID   string    `xml:"prepay_id"`
	CodeUrl    string    `xml:"code_url"`
}

func (wx *Wechat) UnifiedOrder(req *UnifiedOrderReq) (*unifiedOrderRes, error) {
	req.AppID = wx.AppID

	xmlParamsByte, xmlErr := xml.Marshal(req)
	if xmlErr != nil {
		return nil, xmlErr
	}

	res, httpErr := net.HttpPost(unifiedOrderUrl, string(xmlParamsByte), nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	var data unifiedOrderRes
	if xmlErr := xml.Unmarshal(res, &data); xmlErr != nil {
		return nil, xmlErr
	}

	return &data, nil
}

// ==================== 支付通知 ====================
type PaymentNotifyReq struct {
	ReturnCode         string    `xml:"return_code" validate:"required"`
	ReturnMsg          string    `xml:"return_msg"`
	AppID              string    `xml:"appid" validate:"required"`
	MchID              string    `xml:"mch_id" validate:"required"`
	DeviceInfo         string    `xml:"device_info"`
	NonceStr           string    `xml:"nonce_str" validate:"required"`
	Sign               string    `xml:"sign" validate:"required"`
	SignType           SignType  `xml:"sign_type"`
	ResultCode         string    `xml:"result_code" validate:"required"`
	ErrCode            string    `xml:"err_code"`
	ErrCodeDes         string    `xml:"err_code_des"`
	OpenID             string    `xml:"openid" validate:"required"`
	IsSubscribe        string    `xml:"is_subscribe" validate:"required"`
	TradeType          TradeType `xml:"trade_type" validate:"required"`
	BankType           string    `xml:"bank_type" validate:"required"`
	TotalFee           int64     `xml:"total_fee" validate:"required"`
	SettlementTotalFee int64     `xml:"settlement_total_fee"`
	FeeType            string    `xml:"fee_type"`
	CashFee            int64     `xml:"cash_fee"`
	CashFeeType        string    `xml:"cash_fee_type"`
	CouponFee          int64     `xml:"coupon_fee"`
	CouponCount        int64     `xml:"coupon_count"`
	TransactionID      string    `xml:"transaction_id" validate:"required"`
	OutTradeNo         string    `xml:"out_trade_no"  validate:"required"`
	Attach             string    `xml:"attach"`
	TimeEnd            string    `xml:"time_end" validate:"required"`
}

type PaymentNotifyCode string

const (
	PaymentNotifySuccessReturnCode = PaymentNotifyCode("SUCCESS")
	PaymentNotifySuccessReturnMsg  = "OK"

	PaymentNotifyFailReturnCode = PaymentNotifyCode("FAIL")
)

type PaymentNotifyRes struct {
	XMLName    xml.Name          `xml:"xml"`
	ReturnCode PaymentNotifyCode `xml:"return_code"`
	ReturnMsg  string            `xml:"return_msg"`
}

// ==================== 授权 ====================
type jsCode2SessionRes struct {
	ErrCode    float64 `json:"errcode"`
	ErrMsg     string  `json:"errmsg"`
	SessionKey string  `json:"session_key"`
	OpenID     string  `json:"openid"`
	UnionID    string  `json:"unionid"`
}

func (wx *Wechat) JsCode2Session(code string) (*jsCode2SessionRes, error) {
	queryParam := "?appid=" + wx.AppID + "&secret=" + wx.AppSecret + "&js_code=" + code + "&grant_type=authorization_code"
	res, httpErr := net.HttpGet(jsCode2SessionUrl+queryParam, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	var data jsCode2SessionRes

	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 退款 ====================
type RefundReq struct {
	XMLName  xml.Name `xml:"xml"`
	AppID    string   `xml:"appid"`     // 是，微信分配的小程序ID
	MchID    string   `xml:"mch_id"`    // 是，微信支付分配的商户号
	NonceStr string   `xml:"nonce_str"` // 是，随机字符串，长度要求在32位以内。推荐随机数生成算法
	Sign     string   `xml:"sign"`      // 是，通过签名算法计算得出的签名值，详见签名生成算法
	SignType SignType `xml:"sign_type"` // 否，通过签名算法计算得出的签名值，详见签名生成算法
	//TransactionID string   `xml:"transaction_id"` // 微信生成的订单号，在支付通知中有返回
	OutTradeNo string `xml:"out_trade_no"`
	/* 是，商户系统内部订单号，要求32个字符内，
	只能是数字、大小写字母_-|*@ ，且在同一个商户号下唯一。
	transaction_id、out_trade_no二选一，
	如果同时存在优先级：transaction_id > out_trade_no
	*/
	OutRefundNo   string `xml:"out_refund_no"`   // 是，商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@ ，同一退款单号多次请求只退一笔。
	TotalFee      int64  `xml:"total_fee"`       // 是，订单总金额，单位为分，只能为整数，详见支付金额
	RefundFee     int64  `xml:"refund_fee"`      // 是，退款总金额，订单总金额，单位为分，只能为整数，详见支付金额
	RefundFeeType string `xml:"refund_fee_type"` // 否，货币类型，符合ISO 4217标准的三位字母代码，默认人民币：CNY，其他值列表详见货币类型
	RefundDesc    string `xml:"refund_desc"`
	/* 否，若商户传入，会在下发给用户的退款消息中体现退款原因
	注意：若订单退款金额≤1元，且属于部分退款，则不会在退款消息中体现退款原因
	*/
	RefundAccount string `xml:"refund_account"`
	/* 否，仅针对老资金流商户使用
	REFUND_SOURCE_UNSETTLED_FUNDS---未结算资金退款（默认使用未结算资金退款）
	REFUND_SOURCE_RECHARGE_FUNDS---可用余额退款
	*/
	NotifyUrl string `xml:"notify_url"`
	/* 否，
	异步接收微信支付退款结果通知的回调地址，通知URL必须为外网可访问的url，不允许带参数
	公网域名必须为https，如果是走专线接入，使用专线NAT IP或者私有回调域名可使用http。
	如果参数中传了notify_url，则商户平台上配置的回调地址将不会生效。
	*/
}

type refundRes struct {
	ReturnCode PaymentNotifyCode `xml:"return_code"` //SUCCESS/FAIL
	ReturnMsg  string            `xml:"return_msg"`
	/* 返回信息，如非空，为错误原因
	签名失败
	参数格式校验错误
	*/
	ResultCode string `xml:"result_code"`
	/*业务结果
	SUCCESS/FAIL
	SUCCESS退款申请接收成功，结果通过退款查询接口查询
	FAIL 提交业务失败
	*/
	ErrCode             string `xml:"err_code"`
	ErrCodeDes          string `xml:"err_code_des"`
	AppID               string `xml:"appid"`
	MchID               string `xml:"mch_id"`
	NonceStr            string `xml:"nonce_str"`
	Sign                string `xml:"sign"`
	TransactionID       string `xml:"transaction_id"`
	OutTradeNo          string `xml:"out_trade_no"`
	OutRefundNo         string `xml:"out_refund_no"`
	RefundID            string `xml:"refund_id"`
	RefundFee           int64  `xml:"refund_fee"`
	SettlementRefundFee int64  `xml:"settlement_refund_fee"`
	TotalFee            int64  `xml:"total_fee"`
	SettlementTotalFee  int64  `xml:"settlement_total_fee"`
	FeeType             string `xml:"fee_type"`
	CashFee             int64  `xml:"cash_fee"`
	CashFeeType         string `xml:"cash_fee_type"`
	CashRefundFee       int64  `xml:"cash_refund_fee"`
	CouponRefundFee     int64  `xml:"coupon_refund_fee"`
	CouponRefundCount   int64  `xml:"coupon_refund_count"`
}

func (wx *Wechat) Refund(req *RefundReq, certKey, cert string) (*refundRes, error) {
	req.AppID = wx.AppID

	xmlParamsByte, xmlErr := xml.Marshal(req)
	if xmlErr != nil {
		return nil, xmlErr
	}

	res, httpErr := net.HttpPost(refundUrl, string(xmlParamsByte), nil, true, cert, certKey)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println(string(res))
	fmt.Println(fmt.Sprintf("%+v", string(res)))

	wx.ZeroValueProcess(res)

	var data refundRes
	if xmlErr := xml.Unmarshal(res, &data); xmlErr != nil {
		return nil, xmlErr
	}

	return &data, nil
}

// ==================== 退款通知 ====================
type RefundNotifyReq struct {
	ReturnCode string `xml:"return_code" validate:"required"`
	ReturnMsg  string `xml:"return_msg"`
	AppID      string `xml:"appid" validate:"required"`
	MchID      string `xml:"mch_id" validate:"required"`
	NonceStr   string `xml:"nonce_str" validate:"required"`
	ReqInfo    string `xml:"req_info" validate:"required"`
}

type RefundStatus string

const (
	RefundStatusSuccess     = RefundStatus("SUCCESS")
	RefundStatusChange      = RefundStatus("CHANGE")
	RefundStatusRefundClose = RefundStatus("REFUNDCLOSE")
)

type refundReqInfo struct {
	TransactionID       string       `xml:"transaction_id" validate:"required"`
	OutTradeNo          string       `xml:"out_trade_no"  validate:"required"`
	RefundID            string       `xml:"refund_id" validate:"required"`
	OutRefundNo         string       `xml:"out_refund_no"  validate:"required"`
	TotalFee            int64        `xml:"total_fee" validate:"required"`
	SettlementTotalFee  int64        `xml:"settlement_total_fee"` //当该订单有使用非充值券时，返回此字段。应结订单金额=订单金额-非充值代金券金额，应结订单金额<=订单金额。
	RefundFee           int64        `xml:"refund_fee" validate:"required"`
	SettlementRefundFee int64        `xml:"settlement_refund_fee" validate:"required"` //退款金额=申请退款金额-非充值代金券退款金额，退款金额<=申请退款金额
	RefundStatus        RefundStatus `xml:"refund_status" validate:"required"`
	SuccessTime         string       `xml:"success_time"` //资金退款至用户账号的时间，格式2017-12-15 09:46:01
	RefundRecvAccout    string       `xml:"refund_recv_accout" validate:"required"`
	/*退款入账账户
	取当前退款单的退款入账方
	1）退回银行卡：
	{银行名称}{卡类型}{卡尾号}
	2）退回支付用户零钱:
	支付用户零钱
	3）退还商户:
	商户基本账户
	商户结算银行账户
	4）退回支付用户零钱通:
	支付用户零钱通
	*/
	RefundAccount string `xml:"refund_account" validate:"required"`
	/*退款资金来源
	REFUND_SOURCE_RECHARGE_FUNDS 可用余额退款/基本账户
	REFUND_SOURCE_UNSETTLED_FUNDS 未结算资金退款
	*/
	RefundRequestSource string `xml:"refund_request_source" validate:"required"`
	/*退款发起来源
	API 接口
	VENDOR_PLATFORM 商户平台
	*/
}

type NotifyCode string

const (
	NotifySuccessReturnCode = NotifyCode("SUCCESS")
	NotifySuccessReturnMsg  = "OK"

	NotifyFailReturnCode = NotifyCode("FAIL")
)

type NotifyRes struct {
	XMLName    xml.Name   `xml:"xml"`
	ReturnCode NotifyCode `xml:"return_code"`
	ReturnMsg  string     `xml:"return_msg"`
}

func (wx *Wechat) DecodeRefundReqInfo(reqInfo, apiKey string) (*refundReqInfo, error) {
	cipher, base64Err := encode.Base64Decoded(reqInfo)
	if base64Err != nil {
		return nil, base64Err
	}
	secret, md5Err := encrypt.MD5(apiKey)
	if md5Err != nil {
		return nil, md5Err
	}

	aes := encrypt.NewAES(encrypt.ECB)
	res, aesErr := aes.Decrypt([]byte(cipher), secret)
	if aesErr != nil {
		return nil, aesErr
	}

	var data refundReqInfo
	if xmlErr := xml.Unmarshal(res, &data); xmlErr != nil {
		return nil, xmlErr
	}

	return &data, nil
}

// ==================== 获取access_token ====================
type accessTokenRes struct {
	ErrCode     float64 `json:"errcode"`
	ErrMsg      string  `json:"errmsg"`
	ExpiresIn   float64 `json:"expires_in"`   // 凭证有效时间，单位：秒。目前是7200秒之内的值。
	AccessToken string  `json:"access_token"` // 获取到的凭证
}

func (wx *Wechat) GetAccessToken() (*accessTokenRes, error) {
	queryParam := "?grant_type=client_credential&appid=" + wx.AppID + "&secret=" + wx.AppSecret
	res, httpErr := net.HttpGet(accessTokenUrl+queryParam, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	var data accessTokenRes

	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 获取二维码 ====================
type GetWxACodeUnLimitReq struct {
	Scene string `json:"scene"`
	/* 是
	最大32个可见字符，只支持数字，
	大小写英文以及部分特殊字符：!#$&'()*+,/:;=?@-._~，
	其它字符请自行编码为合法字符（因不支持%，中文无法使用 urlencode 处理，
	请使用其他编码方式）
	*/
	Page string `json:"page"`
	/* 否
	必须是已经发布的小程序存在的页面（否则报错），
	例如 pages/index/index, 根路径前不要填加 /,
	不能携带参数（参数请放在scene字段里），
	如果不填写这个字段，默认跳主页面
	*/
	Width     int64 `json:"width"`      // 否 二维码的宽度，单位 px，最小 280px，最大 1280px
	AutoColor bool  `json:"auto_color"` // 否 自动配置线条颜色，如果颜色依然是黑色，则说明不建议配置主色调，默认 false
	LineColor Color `json:"line_color"` // 否 auto_color 为 false 时生效，使用 rgb 设置颜色 例如 {"r":"xxx","g":"xxx","b":"xxx"} 十进制表示
	IsHyaline bool  `json:"is_hyaline"` // 否 是否需要透明底色，为 true 时，生成透明底色的小程序
}

type Color struct {
	R string `json:"r"`
	G string `json:"g"`
	B string `json:"b"`
}

type getWxACodeUnLimitRes struct {
	ErrCode float64 `json:"errcode"`
	ErrMsg  string  `json:"errmsg"`
	Qrcode  []byte  `json:"qrcode"`
}

func (wx *Wechat) GetWxACodeUnLimit(accessToken string, req *GetWxACodeUnLimitReq) (*getWxACodeUnLimitRes, error) {
	queryParam := "?access_token=" + accessToken
	res, httpErr := net.HttpPost(wxACodeUnLimitUrl+queryParam, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("============qrcode========")
	fmt.Println(string(res))
	var data getWxACodeUnLimitRes

	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
	//	ioutil.WriteFile("./output.jpg", []byte(qrcode), 0666)
	//
	//	return outputer.Json(Success, SuccessMsg, qrcode)
}

package ecommerce

import (
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/encode"
	"github.com/MangoMilk/go-kit/encrypt"
	"github.com/MangoMilk/go-kit/net"
)

// #################### 普通支付 ####################

// ==================== 小程序下单 ====================
type MiniProgramPayReq struct {
	SpAppID  string `json:"sp_appid"` // 是，服务商应用ID，[1,32]，服务商申请的公众号或移动应用appid。示例值：wx8888888888888888
	SpMchID  string `json:"sp_mchid"` // 是，服务商户号，[1,32]，由微信支付生成并下发。示例值：1230000109
	SubAppID string `json:"sub_appid"`
	/* 否，二级商户应用ID，[1,32]，
	二级商户申请的公众号或移动应用appid。若sub_openid有传的情况下，sub_appid必填，且sub_appid需与sub_openid对应。
	示例值：wxd678efh567hg6999
	*/
	SubMchID    string `json:"sub_mchid"`   // 是，二级商户号，[1,32]，二级商户的商户号，由微信支付生成并下发。示例值：1900000109
	Description string `json:"description"` // 是，商品详细描述，[1,127]，商品描述。示例值：Image形象店-深圳腾大-QQ公仔
	OutTradeNo  string `json:"out_trade_no"`
	/* 是，商户系统内部订单号，[6,32]
	只能是数字、大小写字母_-*且在同一个商户号下唯一，详见【商户订单号】。
	特殊规则：最小字符长度为6。
	示例值：1217752501201407033233368018
	*/
	TimeExpire string `json:"time_expire"`
	/* 否，交易结束时间，[1,64]
	 订单失效时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，
	T出现在字符串中，表示time元素的开头，
	HH:mm:ss表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日 13点29分35秒。
	示例值：2018-06-08T10:34:56+08:00
	*/
	Attach    string `json:"attach"` // 否，附加数据，[1,128]，在查询API和支付通知中原样返回，可作为自定义参数使用。示例值：自定义数据
	NotifyUrl string `json:"notify_url"`
	/* 是，通知地址，[1,256]
	通知URL必须为直接可访问的URL，不允许携带查询串。
	格式：URL
	示例值：https://www.weixin.qq.com/wxpay/pay.php
	*/
	GoodsTag   string     `json:"goods_tag"`   // 否，订单优惠标记，[1,32]，订单优惠标记。示例值：WXG
	SettleInfo settleInfo `json:"settle_info"` // 否，结算信息
	Amount     amount     `json:"amount"`      // 是，订单金额信息
	Payer      payer      `json:"payer"`       // 是，支付者信息
	Detail     detail     `json:"detail"`      // 否，优惠功能
	SceneInfo  sceneInfo  `json:"scene_info"`  // 否，支付场景描述
}

type settleInfo struct {
	ProfitSharing bool `json:"profit_sharing"`
	/* 否	是否指定分账，枚举值
	true：是
	false：否
	示例值：true
	*/
	SubsidyAmount int64 `json:"subsidy_amount"`
	/* 补差金额	否	SettleInfo.profit_sharing为true时，该金额才生效。
	注意：单笔订单最高补差金额为5000元
	示例值：10
	*/
}

type amount struct {
	Total    int    `json:"total"`    // 总金额		是	订单总金额，单位为分。示例值：100
	Currency string `json:"currency"` // 货币类型[1,16]	否	CNY：人民币，境内商户号仅支持人民币。示例值：CNY
	// 支付回调时有下面2个字段
	PayerTotal    int    `json:"payer_total"`    //用户支付金额		否	用户支付金额，单位为分。示例值：100
	PayerCurrency string `json:"payer_currency"` //用户支付币种	[1,16]	否	用户支付币种。	示例值：CNY
}

type payer struct {
	SpOpenid  string `json:"sp_openid"`  // 用户服务标识	[1,128]	二选一	用户在服务商appid下的唯一标识。示例值：oUpF8uMuAJO_M2pxb1Q9zNjWeS6o
	SubOpenid string `json:"sub_openid"` //用户子标识	[1,128]	用户在子商户appid下的唯一标识。若传sub_openid，那sub_appid必填。示例值：oUpF8uMuAJO_M2pxb1Q9zNjWeS6o

}

type detail struct {
	CostPrice int `json:"cost_price"`
	/* 订单原价	否
	1、商户侧一张小票订单可能被分多次支付，订单原价用于记录整张小票的交易金额。
	2、当订单原价与支付金额不相等，则不享受优惠。
	3、该字段主要用于防止同一张小票分多次支付，以享受多次优惠的情况，正常支付订单不必上传此参数。
	示例值：608800
	*/
	InvoiceID   string        `json:"invoice_id"`   // 商品小票ID[1,32]	否	商家小票ID。	示例值：微信123
	GoodsDetail []goodsDetail `json:"goods_detail"` //单品列表 否	单品列表信息。条目个数限制：【1，6000】

}

type goodsDetail struct {
	MerchantGoodsID  string `json:"merchant_goods_id"`  //商户侧商品编码	[1,32]	是	由半角的大小写字母、数字、中划线、下划线中的一种或几种组成。示例值：1246464644
	WechatPayGoodsID string `json:"wechatpay_goods_id"` //微信侧商品编码	 [1,32]	否	微信支付定义的统一商品编号（没有可不传）。示例值：1001
	GoodsName        string `json:"goods_name"`         //商品名称	[1,256]	否	商品的实际名称。示例值：iPhoneX 256G
	Quantity         int    `json:"quantity"`           //商品数量	是	用户购买的数量。示例值：1
	UnitPrice        int    `json:"unit_price"`         //商品单价	是	商品单价，单位为分。示例值：828800
}

type sceneInfo struct {
	PayerClientIP string    `json:"payer_client_ip"` //用户终端IP	[1,45]	是	用户的客户端IP，支持IPv4和IPv6两种格式的IP地址。示例值：14.23.150.211
	DeviceID      string    `json:"device_id"`       //商户端设备号	[1,32]	否	商户端设备号（门店号或收银设备ID）。	示例值：013467007045764
	StoreInfo     storeInfo `json:"store_info"`      //商户门店信息	否	商户门店信息
}

type storeInfo struct {
	ID       string `json:"id"`        //门店编号	[1,32]	是	商户侧门店编号。示例值：0001
	Name     string `json:"name"`      //门店名称	[1,256]	否	商户侧门店名称。示例值：腾讯大厦分店
	AreaCode string `json:"area_code"` //地区编码	[1,32]	否	地区编码，详细请见省市区编号对照表。示例值：440305
	Address  string `json:"address"`   //详细地址	[1,512]	否	详细的商户门店地址。示例值：广东省深圳市南山区科技中一道10000号
}

type miniProgramPayRes struct {
	//ReturnCode string `json:"return_code"`
	//ReturnMsg  string `json:"return_msg"`
	//AppID      string `json:"appid"`
	//MchID      string `json:"mch_id"`
	//DeviceInfo string `json:"device_info"`
	//NonceStr   string `json:"nonce_str"`
	//Sign       string `json:"sign"`
	//ResultCode string `json:"result_code"`
	//ErrCode    string `json:"err_code"`
	//ErrCodeDes string `json:"err_code_des"`
	//TradeType  TradeType `json:"trade_type"`
	PrepayID string `json:"prepay_id"`
	/* 是，预支付交易会话标识
	预支付交易会话标识。用于后续接口调用中使用，该值有效期为2小时
	示例值：wx201410272009395522657a690389285100
	*/
	//CodeUrl    string    `json:"code_url"`
}

func MiniProgramPay(req *MiniProgramPayReq) (*miniProgramPayRes, error) {
	api := "https://api.mch.weixin.qq.com/v3/pay/partner/transactions/jsapi"
	res, httpErr := net.HttpPost(api, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================MiniProgramPay")
	fmt.Println(string(res))

	var data miniProgramPayRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 支付/退款通知 ====================
type NotifyReq struct {
	ID         string `json:"id"` // 通知ID	[1,36]	是	通知的唯一ID。示例值：EV-2018022511223320873
	CreateTime string `json:"create_time"`
	/* 通知创建时间	[1,16]	是
	通知创建的时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，
	T出现在字符串中，表示time元素的开头，
	HH:mm:ss.表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示北京时间2015年05月20日13点29分35秒。
	示例值：2015-05-20T13:29:35+08:00
	*/
	EventType string `json:"event_type"`
	/* 通知类型	[1,32]	是	通知的类型，
	REFUND.SUCCESS：退款成功通知
	REFUND.ABNORMAL：退款异常通知
	REFUND.CLOSED：退款关闭通知
	TRANSACTION.SUCCESS：支付成功通知
	示例值：TRANSACTION.SUCCESS
	*/
	ResourceType string   `json:"resource_type"` // 通知数据类型	[1,32]	是	通知的资源数据类型，支付成功通知为encrypt-resource。	示例值：encrypt-resource
	Resource     resource `json:"resource"`      // 通知数据		是	通知资源数据，json格式，见示例
	Summary      string   `json:"summary"`       // 回调摘要	[1,64]	是	回调摘要。示例值：支付成功
}

type resource struct {
	Algorithm      string `json:"algorithm"`       // 加密算法类型	[1,32]	是	对开启结果数据进行加密的加密算法，目前只支持AEAD_AES_256_GCM。示例值：AEAD_AES_256_GCM
	Ciphertext     string `json:"ciphertext"`      // 数据密文	[1,1048576]	是	Base64编码后的开启/停用结果数据密文。示例值：sadsadsadsad
	AssociatedData string `json:"associated_data"` // 附加数据	[1,16]	否	附加数据。	示例值：fdasfwqewlkja484w
	Nonce          string `json:"nonce"`           // 随机串	[1,16]	是	加密使用的随机串。	示例值：fdasflkja484w
}

type orderDetail struct {
	SpAppID    string `json:"sp_appid"`  //服务商应用ID[1,32]	是	服务商申请的公众号或移动应用appid。示例值：wx8888888888888888
	SpMchID    string `json:"sp_mchid"`  //服务商户号	[1,32]	是	服务商户号，由微信支付生成并下发。	示例值：1230000109
	SubAppID   string `json:"sub_appid"` //二级商户应用ID	[1,32]	否	二级商户申请的公众号或移动应用appid。示例值：wxd678efh567hg6999
	SubMchID   string `json:"sub_mchid"` //二级商户号	[1,32]	是	二级商户的商户号，由微信支付生成并下发。	示例值：1900000109
	OutTradeNo string `json:"out_trade_no"`
	/*商户订单号	[1,32]	是	商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一。
	特殊规则：最小字符长度为6
	示例值：1217752501201407033233368018
	*/
	TransactionID string `json:"transaction_id"` //微信支付订单号	[1,32]	否	微信支付系统生成的订单号。	示例值：1217752501201407033233368018
	TradeType     string `json:"trade_type"`
	/*交易类型	[1,16]	否	交易类型，枚举值：
	JSAPI：公众号支付
	NATIVE：扫码支付
	APP：APP支付
	MICROPAY：付款码支付
	MWEB：H5支付
	FACEPAY：刷脸支付
	示例值：MICROPAY
	*/
	TradeState string `json:"trade_state"`
	/*交易状态	[1,32]	是	交易状态，枚举值：
	SUCCESS：支付成功
	REFUND：转入退款
	NOTPAY：未支付
	CLOSED：已关闭
	REVOKED：已撤销（付款码支付）
	USERPAYING：用户支付中（付款码支付）
	PAYERROR：支付失败(其他原因，如银行返回失败)
	示例值：SUCCESS
	*/
	TradeStateDesc string `json:"trade_state_desc"` //交易状态描述	[1,256]	是	交易状态描述。示例值：支付失败，请重新下单支付
	BankType       string `json:"bank_type"`        //付款银行	[1,16]	否	银行类型，采用字符串类型的银行标识。银行标识请参考《银行类型对照表》。示例值：CMC
	Attach         string `json:"attach"`           //附加数据	[1,128]	否	附加数据，在查询API和支付通知中原样返回，可作为自定义参数使用。示例值：自定义数据
	SuccessTime    string `json:"success_time"`
	/* 支付完成时间	[1,64]	否
	遵循rfc3339标准格式，格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，
	YYYY-MM-DD表示年月日，
	T出现在字符串中，表示time元素的开头，
	HH:mm:ss表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日 13点29分35秒。
	示例值：2018-06-08T10:34:56+08:00
	*/
	Payer           payer             `json:"payer"`            //支付者	否	支付者信息
	Amount          amount            `json:"amount"`           //订单金额	是	订单金额信息
	SceneInfo       sceneInfo         `json:"scene_info"`       //场景信息	否	支付场景信息描述
	PromotionDetail []promotionDetail `json:"promotion_detail"` //优惠功能	否	优惠功能，享受优惠时返回该字段。
}

type promotionDetail struct {
	CouponID string `json:"coupon_id"` //券ID	[1,32]	是	券ID。示例值：109519
	Name     string `json:"name"`      //优惠名称	[1,64]	否	优惠名称。示例值：单品惠-6
	Scope    string `json:"scope"`
	/*优惠范围	[1,32]	否
	GLOBAL：全场代金券
	SINGLE：单品优惠
	示例值：GLOBAL
	*/
	Type string `json:"type"`
	/*优惠类型	[1,32]	否	CASH：充值
	NOCASH：预充值
	示例值：CASH
	*/
	Amount              int                   `json:"amount"`               //优惠券面额	是	优惠券面额。示例值：100
	StockID             string                `json:"stock_id"`             //活动ID	[1,32]	否	活动ID。示例值：931386
	WechatPayContribute int                   `json:"wechatpay_contribute"` //微信出资		否	微信出资，单位为分。示例值：0
	MerchantContribute  int                   `json:"merchant_contribute"`  //商户出资		否	商户出资，单位为分。	示例值：0
	OtherContribute     int                   `json:"other_contribute"`     //其他出资	否	其他出资，单位为分。示例值：0
	Currency            string                `json:"currency"`             //优惠币种	[1,16]	否	CNY：人民币，境内商户号仅支持人民币。示例值：CNY
	GoodsDetail         []promotionGoodDetail `json:"good_detail"`          //单品列表		否	单品列表信息
}

type promotionGoodDetail struct {
	GoodsID        string `json:"goods_id"`        //商品编码	[1,32]	是	示例值：M1006
	DiscountAmount int    `json:"discount_amount"` //微信侧商品编码	 [1,32]	否	商品优惠金额。	示例值：0
	GoodsRemark    string `json:"goods_remark"`    //商品名称	[1,128]	否	商品备注信息。	示例值：商品备注信息
	Quantity       int    `json:"quantity"`        //商品数量	是	用户购买的数量。示例值：1
	UnitPrice      int    `json:"unit_price"`      //商品单价	是	商品单价，单位为分。示例值：828800
}

func DecodeNotifyCiphertext(notifyCiphertext, apiKey string) (*orderDetail, error) {
	cipher, base64Err := encode.Base64Decoded(notifyCiphertext)
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

	var data orderDetail
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

type NotifyCode string

const (
	NotifySuccessReturnCode = NotifyCode("SUCCESS")
	NotifySuccessReturnMsg  = "OK"

	NotifyFailReturnCode = NotifyCode("FAIL")
)

type NotifyRes struct {
	Code    NotifyCode `json:"code"`
	Message string     `json:"message"`
}

// ==================== 查询订单(微信支付订单号查询) ====================
func QueryOrderByTransactionID(spMchID, subMchID, transactionID string) (*orderDetail, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/pay/partner/transactions/id/%v?sp_mchid=%v&sub_mchid=%v", transactionID, spMchID, subMchID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryOrderByTransactionID")
	fmt.Println(string(res))

	var data orderDetail
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询订单(商户订单号查询) ====================
func QueryOrderByOutTradeNo(spMchID, subMchID, outTradeNo string) (*orderDetail, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/pay/partner/transactions/out-trade-no/%v?sp_mchid=%v&sub_mchid=%v", outTradeNo, spMchID, subMchID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryOrderByOutTradeNo")
	fmt.Println(string(res))

	var data orderDetail
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// #################### 合并支付 ####################

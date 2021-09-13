package ecommerce

import (
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/net"
)

type RefundReq struct {
	SpAppID    string `json:"sp_appid"`  //服务商应用ID[1,32]	是	服务商申请的公众号或移动应用appid。示例值：wx8888888888888888
	SubAppID   string `json:"sub_appid"` //二级商户应用ID	[1,32]	否	二级商户申请的公众号或移动应用appid。示例值：wxd678efh567hg6999
	SubMchID   string `json:"sub_mchid"` //二级商户号	[1,32]	是	二级商户的商户号，由微信支付生成并下发。	示例值：1900000109
	OutTradeNo string `json:"out_trade_no"`
	/*商户订单号	[1,32]	是	商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一。
	特殊规则：最小字符长度为6
	示例值：1217752501201407033233368018
	*/
	TransactionID string `json:"transaction_id"` //微信支付订单号	[1,32]	否	微信支付系统生成的订单号。	示例值：1217752501201407033233368018
	OutRefundNo   string //商户退款单号	[1,64]	是	body 商户系统内部的退款单号，商户系统内部唯一，只能是数字、大小写字母_-|*@，同一退款单号多次请求只退一笔。示例值：1217752501201407033233368018
	Reason        string `json:"reason"`
	/*退款原因	[1,80]	否	body 若商户传入，会在下发给用户的退款消息中体现退款原因。
	注意：若订单退款金额≤1元，且属于部分退款，则不会在退款消息中体现退款原因
	示例值：商品已售完
	*/
	Amount        refundAmount `json:"amount"`     //订单金额	是	订单金额信息
	NotifyUrl     string       `json:"notify_url"` //退款结果回调url	[1,256]	否	body 异步接收微信支付退款结果通知的回调地址，通知url必须为外网可访问的url，不能携带参数。 如果参数中传了notify_url，则商户平台上配置的回调地址将不会生效，优先回调当前传的地址。示例值：https://weixin.qq.com
	RefundAccount string       `json:"refund_account"`
	/*退款出资商户	[1, 32]	否	body电商平台垫资退款专用参数。
	需先确认已开通此功能后，才能使用。若需要开通，请联系微信支付客服。
	枚举值：
	REFUND_SOURCE_PARTNER_ADVANCE : 电商平台垫付，需要向微信支付申请开通
	REFUND_SOURCE_SUB_MERCHANT : 二级商户，默认值
	注意：
	1、电商平台垫资退款专用参数，需先确认已开通此功能后，才能使用。 若需要开通，请联系微信支付客服。
	2、若传入REFUND_SOURCE_PARTNER_ADVANCE，代表使用垫付退款功能。实际出款账户为开通该功能时商户指定的出款账户，实际以退款申请受理结果或查单结果为准。
	示例值：REFUND_SOURCE_SUB_MERCHANT
	*/
	FundsAccount string `json:"funds_account"`
	/*资金账户	[1, 32]	否
	body若订单处于待分账状态，可以传入此参数，指定退款资金来源账户。
	当该字段不存在时，默认使用订单交易资金所在账户出款，即待分账时使用不可用余额的资金进行退款，已分账或无分账时使用可用余额的资金进行退款。
	AVAILABLE：可用余额
	示例值：AVAILABLE
	*/
}

type refundAmount struct {
	Total          int    `json:"total"`           // 是，原支付交易的订单总金额，币种的最小单位，只能为整数。示例值：888
	Currency       string `json:"currency"`        // 是，退款币种 符合ISO 4217标准的三位字母代码，目前只支持人民币：CNY。
	Refund         int    `json:"refund"`          // 是，退款金额，币种的最小单位，只能为整数，不能超过原订单支付金额。示例值：888
	PayerRefund    int    `json:"payer_refund"`    //用户退款金额	是	退款给用户的金额，不包含所有优惠券金额。示例值：888
	DiscountRefund int    `json:"discount_refund"` //优惠退款金额		是	优惠券的退款金额，原支付单的优惠按比例退款。	示例值：888
	// 退款回调时有下面1个字段
	PayerTotal int `json:"payer_total"` //用户支付金额		否	用户支付金额，单位为分。示例值：100
}

type refundRes struct {
	RefundID        string                  `json:"refund_id"`        //微信退款单号	[1,32]	是	微信支付退款订单号。示例值：1217752501201407033233368018
	OutRefundNo     string                  `json:"out_refund_no"`    //商户退款单号	[1,64]	是	商户系统内部的退款单号，商户系统内部唯一，同一退款单号多次请求只退一笔。示例值：1217752501201407033233368018
	CreateTime      string                  `json:"create_time"`      //退款创建时间	[1,64]	是	退款受理时间，遵循rfc3339标准格式，格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，T出现在字符串中，表示time元素的开头，HH:mm:ss表示时分秒，TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日13点29分35秒。示例值：2018-06-08T10:34:56+08:00
	Amount          refundAmount            `json:"amount"`           //订单金额	是	订单金额信息
	PromotionDetail []refundPromotionDetail `json:"promotion_detail"` //优惠退款详情		否	优惠退款功能信息，discount_refund>0时，返回该字段。示例值：见示例

	RefundAccount string `json:"refund_account"`
	/*退款资金来源	[1, 32]	否	枚举值：
	REFUND_SOURCE_PARTNER_ADVANCE : 电商平台垫付
	REFUND_SOURCE_SUB_MERCHANT : 二级商户，默认值
	示例值：REFUND_SOURCE_SUB_MERCHANT
	*/
}

type refundPromotionDetail struct {
	PromotionID string `json:"promotion_id"` //券ID	[1,32]	是	券或者立减优惠id。示例值：109519
	Scope       string `json:"scope"`
	/*优惠范围	[1,32]	是	枚举值：
	GLOBAL：全场代金券
	SINGLE：单品优惠
	示例值：SINGLE
	*/
	Type string `json:"type"`
	/*优惠类型[1,32]		是	枚举值：
	COUPON：充值型代金券，商户需要预先充值营销经费
	DISCOUNT：免充值型优惠券，商户不需要预先充值营销经费
	示例值：DISCOUNT
	*/
	Amount       int `json:"amount"`        //优惠券面额	是	用户享受优惠的金额（优惠券面额=微信出资金额+商家出资金额+其他出资方金额 ）。示例值：5
	RefundAmount int `json:"refund_amount"` //优惠退款金额	是	代金券退款金额<=退款金额，退款金额-代金券或立减优惠退款金额为现金，说明详见《代金券或立减优惠》 。示例值：100
}

func Refund(req *RefundReq) (*refundRes, error) {
	api := "https://api.mch.weixin.qq.com/v3/ecommerce/refunds/apply"
	res, httpErr := net.HttpPost(api, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================Refund")
	fmt.Println(string(res))

	var data refundRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 退款密文 ====================
type RefundCiphertext struct {
	SpMchID    string `json:"sp_mchid"`  //服务商户号	[1,32]	是	服务商户号，由微信支付生成并下发。	示例值：1230000109
	SubMchID   string `json:"sub_mchid"` //二级商户号	[1,32]	是	二级商户的商户号，由微信支付生成并下发。	示例值：1900000109
	OutTradeNo string `json:"out_trade_no"`
	/*商户订单号	[1,32]	是	商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一。
	特殊规则：最小字符长度为6
	示例值：1217752501201407033233368018
	*/
	TransactionID string `json:"transaction_id"` //微信支付订单号	[1,32]	否	微信支付系统生成的订单号。	示例值：1217752501201407033233368018
	RefundID      string `json:"refund_id"`      //微信退款单号	[1,32]	是	微信支付退款订单号。示例值：1217752501201407033233368018
	OutRefundNo   string `json:"out_refund_no"`  //商户退款单号	[1,64]	是	商户系统内部的退款单号，商户系统内部唯一，同一退款单号多次请求只退一笔。示例值：1217752501201407033233368018
	RefundStatus  string `json:"refund_status"`  /*退款状态	[1,16]	是
	退款状态，枚举值：
	SUCCESS：退款成功
	CLOSE：退款关闭
	ABNORMAL：退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往【服务商平台—>交易中心】，手动处理此笔退款
	示例值：SUCCESS
	*/
	SuccessTime string `json:"success_time"`
	/*退款成功时间	[1,64]	否
	1、退款成功时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，
	YYYY-MM-DD表示年月日，
	T出现在字符串中，
	表示time元素的开头，
	HH:mm:ss表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日13点29分35秒。
	2、当退款状态为退款成功时返回此参数。
	示例值：2018-06-08T10:34:56+08:00
	*/
	UserReceivedAccount string `json:"user_received_account"`
	/*退款入账账户	[1,64]	是
	取当前退款单的退款入账方。
	退回银行卡：{银行名称}{卡类型}{卡尾号}
	退回支付用户零钱: 支付用户零钱
	退还商户: 商户基本账户、商户结算银行账户
	退回支付用户零钱通：支付用户零钱通
	示例值：招商银行信用卡0403
	*/
	Amount        refundAmount `json:"amount"` //订单金额	是	订单金额信息
	RefundAccount string       `json:"refund_account"`
	/*退款出资商户	[1, 32]	否	body电商平台垫资退款专用参数。
	需先确认已开通此功能后，才能使用。若需要开通，请联系微信支付客服。
	枚举值：
	REFUND_SOURCE_PARTNER_ADVANCE : 电商平台垫付，需要向微信支付申请开通
	REFUND_SOURCE_SUB_MERCHANT : 二级商户，默认值
	注意：
	1、电商平台垫资退款专用参数，需先确认已开通此功能后，才能使用。 若需要开通，请联系微信支付客服。
	2、若传入REFUND_SOURCE_PARTNER_ADVANCE，代表使用垫付退款功能。实际出款账户为开通该功能时商户指定的出款账户，实际以退款申请受理结果或查单结果为准。
	示例值：REFUND_SOURCE_SUB_MERCHANT
	*/
	FundsAccount string `json:"funds_account"`
	/*资金账户	[1, 32]	否
	body若订单处于待分账状态，可以传入此参数，指定退款资金来源账户。
	当该字段不存在时，默认使用订单交易资金所在账户出款，即待分账时使用不可用余额的资金进行退款，已分账或无分账时使用可用余额的资金进行退款。
	AVAILABLE：可用余额
	示例值：AVAILABLE
	*/
}

type refundDetail struct {
	RefundID      string `json:"refund_id"`      //微信退款单号	[1,32]	是	微信支付退款订单号。示例值：1217752501201407033233368018
	OutRefundNo   string `json:"out_refund_no"`  //商户退款单号	[1,64]	是	商户系统内部的退款单号，商户系统内部唯一，同一退款单号多次请求只退一笔。示例值：1217752501201407033233368018
	TransactionID string `json:"transaction_id"` //微信支付订单号	[1,32]	否	微信支付系统生成的订单号。	示例值：1217752501201407033233368018
	OutTradeNo    string `json:"out_trade_no"`
	/*商户订单号	[1,32]	是	商户系统内部订单号，只能是数字、大小写字母_-*且在同一个商户号下唯一。
	特殊规则：最小字符长度为6
	示例值：1217752501201407033233368018
	*/
	channel string
	/*退款渠道	[1,16]	是	ORIGINAL：原路退款
	BALANCE：退回到余额
	OTHER_BALANCE：原账户异常退到其他余额账户
	OTHER_BANKCARD：原银行卡异常退到其他银行卡
	示例值： ORIGINAL
	*/
	UserReceivedAccount string `json:"user_received_account"`
	/*退款入账账户	[1,64]	是
	取当前退款单的退款入账方。
	退回银行卡：{银行名称}{卡类型}{卡尾号}
	退回支付用户零钱: 支付用户零钱
	退还商户: 商户基本账户、商户结算银行账户
	退回支付用户零钱通：支付用户零钱通
	示例值：招商银行信用卡0403
	*/
	SuccessTime string `json:"success_time"`
	/*退款成功时间	[1,64]	否
	1、退款成功时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，
	YYYY-MM-DD表示年月日，
	T出现在字符串中，
	表示time元素的开头，
	HH:mm:ss表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日13点29分35秒。
	2、当退款状态status为SUCCESS（退款成功）时返回此参数。
	示例值：2018-06-08T10:34:56+08:00
	*/
	CreateTime string `json:"create_time"`
	/*退款成功时间	[1,64]	否
	1、退款成功时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，
	YYYY-MM-DD表示年月日，
	T出现在字符串中，
	表示time元素的开头，
	HH:mm:ss表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示，北京时间2015年5月20日13点29分35秒。
	2、当退款状态为退款成功时返回此参数。
	示例值：2018-06-08T10:34:56+08:00
	*/
	Status string `json:"status"` /*退款状态	[1,16]	是
	退款状态，枚举值：
	SUCCESS：退款成功
	CLOSE：退款关闭
	ABNORMAL：退款异常，退款到银行发现用户的卡作废或者冻结了，导致原路退款银行卡失败，可前往【服务商平台—>交易中心】，手动处理此笔退款
	示例值：SUCCESS
	*/
	Amount          refundAmount            `json:"amount"`           //订单金额	是	订单退款金额信息
	PromotionDetail []refundPromotionDetail `json:"promotion_detail"` //营销详情		否	优惠退款功能信息，discount_refund>0时，返回该字段。示例值：见示例
	RefundAccount   string                  `json:"refund_account"`
	/*退款出资商户	[1, 32]	否	body电商平台垫资退款专用参数。
	需先确认已开通此功能后，才能使用。若需要开通，请联系微信支付客服。
	枚举值：
	REFUND_SOURCE_PARTNER_ADVANCE : 电商平台垫付，需要向微信支付申请开通
	REFUND_SOURCE_SUB_MERCHANT : 二级商户，默认值
	注意：
	1、电商平台垫资退款专用参数，需先确认已开通此功能后，才能使用。 若需要开通，请联系微信支付客服。
	2、若传入REFUND_SOURCE_PARTNER_ADVANCE，代表使用垫付退款功能。实际出款账户为开通该功能时商户指定的出款账户，实际以退款申请受理结果或查单结果为准。
	示例值：REFUND_SOURCE_SUB_MERCHANT
	*/
	FundsAccount string `json:"funds_account"`
	/*资金账户	[1, 32]	否
	body若订单处于待分账状态，可以传入此参数，指定退款资金来源账户。
	当该字段不存在时，默认使用订单交易资金所在账户出款，即待分账时使用不可用余额的资金进行退款，已分账或无分账时使用可用余额的资金进行退款。
	AVAILABLE：可用余额
	示例值：AVAILABLE
	*/
}

// ==================== 查询退款(微信支付退款单号查询) ====================
func QueryRefundByRefundID(subMchID, refundID string) (*refundDetail, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/refunds/id/%v?sub_mchid=%v", refundID, subMchID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryRefundByRefundID")
	fmt.Println(string(res))

	var data refundDetail
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询退款(商户退款单号查询) ====================
func QueryRefundByOutRefundNo(subMchID, outRefundNo string) (*refundDetail, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/refunds/out-refund-no/%v?sub_mchid=%v", outRefundNo, subMchID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryRefundByOutRefundNo")
	fmt.Println(string(res))

	var data refundDetail
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

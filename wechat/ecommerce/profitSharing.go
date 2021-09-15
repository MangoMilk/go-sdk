package ecommerce

import (
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/net"
)

// ==================== 分账 ====================
type ProfitSharingReq struct {
	AppID         string     `json:"appid"`          //公众账号ID	[1,32]	是	body 电商平台的appid（公众号APPID或者小程序APPID）。示例值：wx8888888888888888
	SubMchID      string     `json:"sub_mchid"`      //二级商户号	[1,32]	是	body 分账出资的电商平台二级商户，填写微信支付分配的商户号。	示例值：1900000109
	TransactionID string     `json:"transaction_id"` //微信订单号	[1,32]	是	body 微信支付订单号。示例值： 4208450740201411110007820472
	OutOrderNo    string     `json:"out_order_no"`   //商户分账单号	[1,64]	是	body 商户系统内部的分账单号，在商户系统内部唯一（单次分账、多次分账、完结分账应使用不同的商户分账单号），同一分账单号多次请求等同一次。示例值：P20150806125346
	Receivers     []receiver `json:"receivers"`      //分账接收方列表		是	body 分账接收方列表，支持设置出资商户作为分账接收方，单次分账最多可有5个分账接收方
	Finish        bool       `json:"finish"`
	/*是否分账完成		是	body 是否完成分账
	1、如果为true，该笔订单剩余未分账的金额会解冻回电商平台二级商户；
	2、如果为false，该笔订单剩余未分账的金额不会解冻回电商平台二级商户，可以对该笔订单再次进行分账。
	示例值：true
	*/
}

type receiver struct {
	Type string `json:"type"`
	/*分账接收方类型	[1,32]	是	分账接收方类型，枚举值：
	MERCHANT_ID：商户
	PERSONAL_OPENID：个人
	示例值：MERCHANT_ID
	*/
	ReceiverAccount string `json:"receiver_account"`
	/*分账接收方账号	[1,64]	是	分账接收方账号：
	类型是MERCHANT_ID时，是商户号（mch_id或者sub_mch_id）
	类型是PERSONAL_OPENID时，是个人openid，openid获取方法
	示例值：1900000109
	*/
	Amount       int    `json:"amount"`      //分账金额	是	分账金额，单位为分，只能为整数，不能超过原订单支付金额及最大分账比例金额。示例值：190
	Description  string `json:"description"` //分账描述	[1,80]	是	分账的原因描述，分账账单中需要体现。示例值：分给商户1900000109
	ReceiverName string `json:"receiver_name"`
	/*分账个人姓名	[1, 10240]	条件选填	可选项，在接收方类型为个人的时可选填，若有值，会检查与 receiver_name 是否实名匹配，不匹配会拒绝分账请求
	1、分账接收方类型是PERSONAL_OPENID时，是个人姓名的密文（选传，传则校验） 此字段的加密方法详见：敏感信息加密说明
	2、使用微信支付平台证书中的公钥
	3、使用RSAES-OAEP算法进行加密
	4、将请求中HTTP头部的Wechatpay-Serial设置为证书序列号
	示例值：hu89ohu89ohu89o
	*/
}

type profitSharingRes struct {
	SubMchID      string `json:"sub_mchid"`      //二级商户号	[1,32]	是	分账出资的电商平台二级商户，填写微信支付分配的商户号。示例值：1900000109
	TransactionID string `json:"transaction_id"` //微信订单号	[1,32]	是	微信支付订单号。示例值： 4208450740201411110007820472
	OutOrderNo    string `json:"out_order_no"`   //商户分账单号	[1,64]	是	商户系统内部的分账单号，在商户系统内部唯一（单次分账、多次分账、完结分账应使用不同的商户分账单号），同一分账单号多次请求等同一次。示例值：P20150806125346
	OrderID       string `json:"order_id"`       //微信分账单号	[1,64]	是	微信分账单号，微信系统返回的唯一标识。示例值： 6754760740201411110007865434
	Status        string `json:"status"`
	/*分账单状态	[1,32]	是	分账单状态，枚举值：
	PROCESSING：处理中
	FINISHED：处理完成
	示例值：FINISHED
	*/
	Receivers []resReceiver `json:"receivers"` //分账接收方列表		是	分账接收方列表
}

type resReceiver struct {
	Amount      int    `json:"amount"`      //分账金额	是	分账金额，单位为分，只能为整数，不能超过原订单支付金额及最大分账比例金额。示例值：190
	Description string `json:"description"` //分账描述	[1,80]	是	分账的原因描述，分账账单中需要体现。示例值：分给商户1900000109
	FailReason  string `json:"fail_reason"`
	/*分账失败原因	[1,32]	否	分账失败原因，当分账结果result为RETURNED（已转回分账方）或CLOSED（已关闭）时，返回该字段
	枚举值：
	ACCOUNT_ABNORMAL : 分账接收账户异常
	NO_RELATION : 分账关系已解除
	RECEIVER_HIGH_RISK : 高风险接收方
	RECEIVER_REAL_NAME_NOT_VERIFIED : 接收方未实名
	示例值：NO_RELATION
	*/
	DetailID   string `json:"detail_id"` //分账明细单号	[1,64]	是	微信分账明细单号，每笔分账业务执行的明细单号，可与资金账单对账使用。示例值：36011111111111111111111
	FinishTime string `json:"finish_time"`
	/*完成时间	[1,64]	是	分账完成时间，遵循RFC3339标准格式，格式为
	YYYY-MM-DDTHH:mm:ss.sss+TIMEZONE，
	YYYY-MM-DD表示年月日，
	T出现在字符串中，
	表示time元素的开头，
	HH:mm:ss.sss表示时分秒毫秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35:120+08:00表示北京时间2015年05月20日13点29分35秒。
	示例值： 2015-05-20T13:29:35.120+08:00
	*/
	ReceiverAccount string `json:"receiver_account"`
	/*分账接收方账号	[1,64]	是	分账接收方账号：
	类型是MERCHANT_ID时，是商户号（mch_id或者sub_mch_id）
	类型是PERSONAL_OPENID时，是个人openid，openid获取方法
	示例值：1900000109
	*/
	ReceiverMchID string `json:"receiver_mchid"` //分账接收商户号	[1,32]	是	仅分账接收方类型为MERCHANT_ID时，填写微信支付分配的商户号。示例值：1900000110
	Result        string `json:"result"`
	/*分账结果	[1,32]	是	分账结果，枚举值：
	PENDING：待分账
	SUCCESS：分账成功
	CLOSED：分账失败已关闭
	示例值：SUCCESS
	*/
	Type string `json:"type"`
	/*分账接收方类型	[1,32]	是	分账接收方类型，枚举值：
	MERCHANT_ID：商户号（mch_id或者sub_mch_id）
	PERSONAL_OPENID：个人openid（由服务商的APPID转换得到）
	PERSONAL_SUB_OPENID：个人sub_openid（由品牌主的APPID转换得到）
	示例值：MERCHANT_ID
	*/
}

func ProfitSharing(req *ProfitSharingReq) (*profitSharingRes, error) {

	api := "https://api.mch.weixin.qq.com/v3/ecommerce/profitsharing/orders"
	res, httpErr := net.HttpPost(api, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================ProfitSharing")
	fmt.Println(string(res))

	var data profitSharingRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询分账结果 ====================
type queryProfitSharingRes struct {
	SubMchID      string `json:"sub_mchid"`      //二级商户号	[1,32]	是	分账出资的电商平台二级商户，填写微信支付分配的商户号。示例值：1900000109
	TransactionID string `json:"transaction_id"` //微信订单号	[1,32]	是	微信支付订单号。	示例值： 4208450740201411110007820472
	OutOrderNo    string `json:"out_order_no"`   //商户分账单号	[1,64]	是	商户系统内部的分账单号，在商户系统内部唯一（单次分账、多次分账、完结分账应使用不同的商户分账单号），同一分账单号多次请求等同一次。示例值：P20150806125346
	OrderID       string `json:"order_id"`       //微信分账单号	[1,64]	是	微信分账单号，微信系统返回的唯一标识。示例值： 008450740201411110007820472
	Status        string `json:"status"`
	/*分账单状态	[1,32]	是	分账单状态，枚举值：
	PROCESSING：处理中
	FINISHED：分账完成
	示例值：FINISHED
	*/
	Receivers         []resReceiver `json:"receivers"`          //分账接收方列表	否	分账接收方列表。当查询分账完结的执行结果时，不返回该字段
	FinishAmount      int           `json:"finish_amount"`      //分账完结金额		否	分账完结的分账金额，单位为分， 仅当查询分账完结的执行结果时，存在本字段。示例值：100
	FinishDescription string        `json:"finish_description"` //分账完结描述	[1,80]	否	分账完结的原因描述，仅当查询分账完结的执行结果时，存在本字段。示例值：分账完结
}

func QueryProfitSharing(subMchID, transactionID, outOrderNo string) (*queryProfitSharingRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/profitsharing/orders?sub_mchid=%v&transaction_id=%v&out_order_no=%v", subMchID, transactionID, outOrderNo)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryProfitSharing")
	fmt.Println(string(res))

	var data queryProfitSharingRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询订单剩余待分账金额 ====================
type queryProfitSharingOrderAmountsRes struct {
	TransactionID string `json:"transaction_id"` //微信订单号	[1,32]	是	微信支付订单号。示例值：4208450740201411110007820472
	UnSplitAmount int    `json:"unsplit_amount"` //订单剩余待分金额	是	订单剩余待分金额，整数，单位为分。示例值：1000
}

func QueryProfitSharingOrderAmounts(transactionID string) (*queryProfitSharingOrderAmountsRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/profitsharing/orders/%v/amounts", transactionID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryProfitSharingOrderAmounts")
	fmt.Println(string(res))

	var data queryProfitSharingOrderAmountsRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 完结分账 ====================
type FinishProfitSharingReq struct {
	SubMchID      string `json:"sub_mchid"`      //二级商户号	[1,32]	是	分账出资的电商平台二级商户，填写微信支付分配的商户号。示例值：1900000109
	TransactionID string `json:"transaction_id"` //微信订单号	[1,32]	是	微信支付订单号。	示例值： 4208450740201411110007820472
	OutOrderNo    string `json:"out_order_no"`   //商户分账单号	[1,64]	是	商户系统内部的分账单号，在商户系统内部唯一（单次分账、多次分账、完结分
	Description   string `json:"description"`    //分账描述	[1,80]	是	分账的原因描述，分账账单中需要体现。示例值：分账完结
}

type finishProfitSharingRes struct {
	SubMchID      string `json:"sub_mchid"`      //二级商户号	[1,32]	是	分账出资的电商平台二级商户，填写微信支付分配的商户号。示例值：1900000109
	TransactionID string `json:"transaction_id"` //微信订单号	[1,32]	是	微信支付订单号。	示例值： 4208450740201411110007820472
	OutOrderNo    string `json:"out_order_no"`   //商户分账单号	[1,64]	是	商户系统内部的分账单号，在商户系统内部唯一（单次分账、多次分账、完结分账应使用不同的商户分账单号），同一分账单号多次请求等同一次。示例值：P20150806125346
	OrderID       string `json:"order_id"`       //微信分账单号	[1,64]	是	微信分账单号，微信系统返回的唯一标识。示例值： 008450740201411110007820472
}

func FinishProfitSharing(req *FinishProfitSharingReq) (*finishProfitSharingRes, error) {

	api := "https://api.mch.weixin.qq.com/v3/ecommerce/profitsharing/finish-order"
	res, httpErr := net.HttpPost(api, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================FinishProfitSharing")
	fmt.Println(string(res))

	var data finishProfitSharingRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

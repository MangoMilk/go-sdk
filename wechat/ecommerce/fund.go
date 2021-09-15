package ecommerce

import (
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/net"
)

type AccountType string

const (
	AccountTypeBasic     = AccountType("BASIC")
	AccountTypeOperation = AccountType("OPERATION")
	AccountTypeFees      = AccountType("FEES")
)

// ==================== 查询二级商户账户实时余额 ====================
type queryBalanceRes struct {
	SubMchID    string      `json:"sub_mchid"` //二级商户号	[1,32]	是	电商平台二级商户号，由微信支付生成并下发。示例值： 1900000109
	AccountType AccountType `json:"account_type"`
	/*账户类型	[1,16]	否	枚举值：
	BASIC：基本账户
	OPERATION：运营账户
	FEES：手续费账户
	示例值： BASIC
	*/
	AvailableAmount int64 `json:"available_amount"` //可用余额	是	可用余额（单位：分），此余额可做提现操作。示例值： 100
	PendingAmount   int64 `json:"pending_amount"`   //不可用余额	否	不可用余额（单位：分）。	示例值： 100
}

func QueryBalance(subMchID string, accountType AccountType) (*queryBalanceRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/fund/balance/%v?account_type=%v", subMchID, accountType)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryBalance")
	fmt.Println(string(res))

	var data queryBalanceRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询二级商户账户日终余额 ====================
type queryEndDayBalanceRes struct {
	SubMchID        string `json:"sub_mchid"`        //二级商户号	[1,32]	是	电商平台二级商户号，由微信支付生成并下发。示例值： 1900000109
	AvailableAmount int64  `json:"available_amount"` //可用余额	是	可用余额（单位：分），此余额可做提现操作。示例值： 100
	PendingAmount   int64  `json:"pending_amount"`   //不可用余额	否	不可用余额（单位：分）。	示例值： 100
}

//date 指定查询商户日终余额的日期，可查询90天内的日终余额。示例值：2019-08-17
func QueryEndDayBalance(subMchID string, date string) (*queryEndDayBalanceRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/fund/enddaybalance/%v?date=%v", subMchID, date)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryEndDayBalance")
	fmt.Println(string(res))

	var data queryEndDayBalanceRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询电商平台账户实时余额 ====================
type queryMerchantBalanceRes struct {
	AvailableAmount int64 `json:"available_amount"` //可用余额	是	可用余额（单位：分），此余额可做提现操作。示例值： 100
	PendingAmount   int64 `json:"pending_amount"`   //不可用余额	否	不可用余额（单位：分）。	示例值： 100
}

func QueryMerchantBalance(accountType AccountType) (*queryMerchantBalanceRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/merchant/fund/balance/%v", accountType)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryMerchantBalance")
	fmt.Println(string(res))

	var data queryMerchantBalanceRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 查询电商平台账户日终余额 ====================
//date 指定查询商户日终余额的日期，可查询90天内的日终余额。示例值：2019-08-17
func QueryMerchantEndDayBalance(accountType AccountType, date string) (*queryMerchantBalanceRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/merchant/fund/dayendbalance/%v?date=%v", accountType, date)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryMerchantEndDayBalance")
	fmt.Println(string(res))

	var data queryMerchantBalanceRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 二级商户余额提现 ====================
type WithdrawReq struct {
	SubMchID     string      `json:"sub_mchid"`      //二级商户号	[1,32]	是	电商平台二级商户号，由微信支付生成并下发。示例值： 1900000109
	OutRequestNo string      `json:"out_request_no"` //商户提现单号	[1, 32]	是	body商户提现单号，由商户自定义生成，必须是字母数字。示例值：20190611222222222200000000012122
	Amount       int         `json:"amount"`         //提现金额	是	body单位：分，金额不能超过8亿元。示例值：1
	Remark       string      `json:"remark"`         //提现备注	[1, 56]	否	body商户对提现单的备注，商户自定义字段。示例值：交易提现
	BankMemo     string      `json:"bank_memo"`      //银行附言	[1, 32]	否	body展示在收款银行系统中的附言，数字、字母最长32个汉字（能否成功展示依赖银行系统支持）。示例值：微信支付提现
	AccountType  AccountType `json:"account_type"`
	/*账户类型	[1,16]	否	枚举值：
	BASIC：基本账户
	OPERATION：运营账户
	FEES：手续费账户
	示例值： BASIC
	*/
}

type withdrawRes struct {
	SubMchID     string `json:"sub_mchid"`      //二级商户号	[1,32]	是	电商平台二级商户号，由微信支付生成并下发。示例值： 1900000109
	WithdrawID   string `json:"withdraw_id"`    //微信支付提现单号	[1, 128]	是	电商平台提交二级商户提现申请后，由微信支付返回的申请单号，作为查询申请状态的唯一标识。示例值：12321937198237912739132791732912793127931279317929791239112123
	OutRequestNo string `json:"out_request_no"` //商户提现单号	[1, 32]	是	body商户提现单号，由商户自定义生成，必须是字母数字。示例值：20190611222222222200000000012122
}

func Withdraw(req *WithdrawReq) (*withdrawRes, error) {
	api := "https://api.mch.weixin.qq.com/v3/ecommerce/fund/withdraw"
	res, httpErr := net.HttpPost(api, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================Withdraw")
	fmt.Println(string(res))

	var data withdrawRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 二级商户查询提现状态(微信支付提现单号查询) ====================
type WithdrawStatus string

const (
	WithdrawStatusCreateSuccess = WithdrawStatus("CREATE_SUCCESS") //受理成功
	WithdrawStatusSuccess       = WithdrawStatus("SUCCESS")        //提现成功
	WithdrawStatusFail          = WithdrawStatus("FAIL")           //提现失败
	WithdrawStatusRefund        = WithdrawStatus("REFUND")         //提现退票
	WithdrawStatusClose         = WithdrawStatus("CLOSE")          //关单
	WithdrawStatusInit          = WithdrawStatus("INIT")           //业务单已创建
)

type queryWithdrawRes struct {
	SubMchID string         `json:"sub_mchid"` //二级商户号	[1,32]	是	电商平台二级商户号，由微信支付生成并下发。示例值： 1900000109
	SpMchID  string         `json:"sp_mchid"`  //电商平台商户号	[1, 32]	是	电商平台商户号。示例值：1800000123
	Status   WithdrawStatus `json:"status"`
	/*提现单状态	[1,16]	是	枚举值：
	CREATE_SUCCESS：受理成功
	SUCCESS：提现成功
	FAIL：提现失败
	REFUND：提现退票
	CLOSE：关单
	INIT：业务单已创建
	示例值：CREATE_SUCCESS
	*/
	WithdrawID   string `json:"withdraw_id"`    //微信支付提现单号	[1, 128]	是	电商平台提交二级商户提现申请后，由微信支付返回的申请单号，作为查询申请状态的唯一标识。示例值：12321937198237912739132791732912793127931279317929791239112123
	OutRequestNo string `json:"out_request_no"` //商户提现单号	[1, 32]	是	body商户提现单号，由商户自定义生成，必须是字母数字。示例值：20190611222222222200000000012122
	Amount       int    `json:"amount"`         //提现金额	是	单位：分。示例值：1
	CreateTime   string `json:"create_time"`
	/* 发起提现时间	[29,29]	是
	通知创建的时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，
	T出现在字符串中，表示time元素的开头，
	HH:mm:ss.表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示北京时间2015年05月20日13点29分35秒。
	示例值：2015-05-20T13:29:35+08:00
	*/
	UpdateTime string `json:"update_time"`
	/* 提现状态更新时间	[29,29]	是
	通知创建的时间，遵循rfc3339标准格式，
	格式为YYYY-MM-DDTHH:mm:ss+TIMEZONE，YYYY-MM-DD表示年月日，
	T出现在字符串中，表示time元素的开头，
	HH:mm:ss.表示时分秒，
	TIMEZONE表示时区（+08:00表示东八区时间，领先UTC 8小时，即北京时间）。
	例如：2015-05-20T13:29:35+08:00表示北京时间2015年05月20日13点29分35秒。
	示例值：2015-05-20T13:29:35+08:00
	*/
	Reason      string      `json:"reason"`    //失败原因	[1, 255]	是	提现失败原因，仅在提现失败、退票、关单时有值。示例值：卡号错误
	Remark      string      `json:"remark"`    //提现备注	[1, 56]	是	商户对提现单的备注，若发起提现时未传入相应值或输入不合法，则该值为空。示例值：交易提现
	BankMemo    string      `json:"bank_memo"` //银行附言	[1, 32]	是	展示在收款银行系统中的附言，由数字、字母、汉字组成（能否成功展示依赖银行系统支持）。若发起提现时未传入相应值或输入不合法，则该值为空。示例值：微信提现
	AccountType AccountType `json:"account_type"`
	/*出款账户类型	[1,16]	是	枚举值：
	BASIC：基本账户
	OPERATION：运营账户
	FEES：手续费账户
	示例值： BASIC
	*/
	AccountNumber string `json:"account_number"` //入账银行账号后四位	[1, 4]	是	服务商提现入账的银行账号，仅显示后四位。示例值：1178
	AccountBank   string `json:"account_bank"`   //入账银行	[1, 10]	是	服务商提现入账的开户银行。	示例值：招商银行
	BankName      string `json:"bank_name"`      //入账银行全称（含支行）	[1, 128]	否	服务商提现入账的开户银行全称（含支行）。示例值：中国工商银行股份有限公司深圳软件园支行
}

func QueryWithdrawByWithdrawID(withdrawID string, subMchID string) (*queryWithdrawRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/fund/withdraw/%v?sub_mchid=%v", withdrawID, subMchID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryWithdrawByWithdrawID")
	fmt.Println(string(res))

	var data queryWithdrawRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 二级商户查询提现状态(商户提现单号查询) ====================
func QueryWithdrawByOutRequestNo(outRequestNo string, subMchID string) (*queryWithdrawRes, error) {
	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/fund/withdraw/out-request-no/%v?sub_mchid=%v", outRequestNo, subMchID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================QueryWithdrawByOutRequestNo")
	fmt.Println(string(res))

	var data queryWithdrawRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

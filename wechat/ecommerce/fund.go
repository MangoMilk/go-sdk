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

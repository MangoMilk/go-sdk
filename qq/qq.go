package qq

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/MangoMilk/go-kit/net"
	goKitTime "github.com/MangoMilk/go-kit/time"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	PayB2CApi = "https://api.qpay.qq.com/cgi-bin/epay/qpay_epay_b2c.cgi"
	PayHBApi  = "https://api.qpay.qq.com/cgi-bin/hongbao/qpay_hb_mch_send.cgi"
)

type QQ struct {
	AppID string
	MchID string
	OPUserID string
	OPUserPassword string
	ApiKey string
	CertFile string
	KeyFile string
}

func NewQQ() *QQ {
	return &QQ{

	}
}

func (qq *QQ)QPayB2C(clientIp string, openId string, apiKey string, tradeNo string, nonceStr string, totalFree int) error {
	var data = make(map[string]interface{})
	data["input_charset"] = "UTF-8"
	data["appid"] = qq.AppID
	data["openid"] = openId
	data["mch_id"] = qq.MchID
	data["nonce_str"] = nonceStr
	data["out_trade_no"] = tradeNo
	data["total_fee"] = totalFree
	data["op_user_id"] = qq.OPUserID
	data["op_user_passwd"] = qq.OPUserPassword
	data["spbill_create_ip"] = clientIp

	// 生成签名（sign）和请求参数（xml）
	var dataKeys []string
	for k, _ := range data {
		dataKeys = append(dataKeys, k)
	}
	sort.Strings(dataKeys)

	// joint data
	var xmlParams string = "<xml>"
	var signStr string = ""
	for _, v := range dataKeys {
		val := data[v].(string)

		xmlParams += "<" + v + ">" + val + "</" + v + ">"
		signStr += (v + "=" + val + "&")
	}

	h := md5.New()
	signByte := []byte(strings.Trim(signStr, "&") + "&key=" + apiKey)
	h.Write(signByte)
	var sign string = strings.ToUpper(hex.EncodeToString(h.Sum(nil)))

	xmlParams += "<sign>" + sign + "</sign>"
	xmlParams += "</xml>"
	res,httpErr:=net.HttpPost(PayB2CApi, xmlParams, nil, true,qq.CertFile,qq.KeyFile)
	if httpErr != nil {
		return httpErr
	}
	var response = make(map[string]interface{})

	if err := xml.Unmarshal(res, &response);err != nil {
		return err
	}

	return nil

	//if int(response["retcode"].(float64)) != 0 || response["return_code"].(string) != "SUCCESS" || response["result_code"].(string) != "SUCCESS" {
	//}
}

type QPayHBRes struct {
	RetCode string `retcode`
	ListID  string `listid`
}

func (qq *QQ)QPayHB(openId string, iconID string, bannerID string, apiKey string, mchBillno string, sender string, nonceStr string, totalFree int, actName string, wishing string) (*QPayHBRes,error) {
	var totalFreeStr = strconv.Itoa(totalFree)
	var data = make(map[string]string)
	data["charset"] = "1"
	data["nonce_str"] = nonceStr
	// mchBillno, _ := strconv.ParseInt(genMchBillNo(Conf.QQ_MCH_ID), 10, 64)
	data["mch_billno"] = mchBillno //genMchBillNo(mchID)
	data["mch_id"] = qq.MchID
	data["mch_name"] = sender //红包发送者名称
	data["qqappid"] = qq.AppID
	data["re_openid"] = openId
	data["total_amount"] = totalFreeStr //发放总金额(单位：分)
	data["total_num"] = "1"             //红包发送总人数(目前限制为1)
	data["wishing"] = wishing           //红包祝福语
	data["act_name"] = actName          //活动名称
	data["icon_id"] = iconID            //qq coin_id
	data["banner_id"] = bannerID        //qq banner_id
	data["min_value"] = "1"             //单个红包的最小金额(单位：分)
	data["max_value"] = totalFreeStr    //单个红包的最大金额(单位：分)

	// 生成签名（sign）和请求参数（xml）
	var dataKeys []string
	for k, _ := range data {
		dataKeys = append(dataKeys, k)
	}
	sort.Strings(dataKeys)

	// joint data
	var signStr string = ""
	var val string
	for _, v := range dataKeys {
		val = data[v]
		signStr += (v + "=" + val + "&")
	}
	// md5 && upper
	h := md5.New()
	signByte := []byte(strings.Trim(signStr, "&") + "&key=" + apiKey)
	h.Write(signByte)
	data["sign"] = strings.ToUpper(hex.EncodeToString(h.Sum(nil)))


	res,httpErr:=net.HttpPost(PayHBApi, data, nil, true,qq.CertFile,qq.KeyFile)
	if httpErr != nil {
		return nil,httpErr
	}

	var response QPayHBRes

	if jsonErr:=json.Unmarshal(res,&response);jsonErr!=nil{
		return nil,jsonErr
	}

	//if response.RetCode != "0" {
	//}

	return &response,nil
}

func (qq *QQ)genMchBillNo(mchId string) string {
	// Long() mch_id+yyyymmdd+10位一天内不能重复的的数字，如：144124610120161010 1234567890
	return mchId + time.Now().Format(goKitTime.FormatSeamlessDate) + fmt.Sprintf("%010d", goKitTime.GetCurDayPassMS())
}

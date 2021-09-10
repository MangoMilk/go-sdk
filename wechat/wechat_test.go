package wechat

import (
	"encoding/xml"
	"fmt"
	"github.com/MangoMilk/go-kit/encrypt"
	"github.com/MangoMilk/go-sdk/aliyun"
	"strconv"
	"testing"
	"time"
)

var (
	wx *Wechat

	appID           = ""
	appSecret       = ""
	mchID           = ""
	apiKey          = ""
	notifyUrl       = ""
	refundNotifyUrl = ""
	certKey         = ""
	cert            = ""
)

func setup() {
	wx = NewWechat(appID, appSecret)
}

func teardown() {

}

func TestMain(m *testing.M) {
	setup()
	m.Run()
	teardown()
}

func TestUnifiedOrder(t *testing.T) {
	nonceStr, _ := encrypt.MD5(strconv.FormatInt(time.Now().Unix(), 10))
	req := UnifiedOrderReq{
		MchID:          mchID,
		NonceStr:       nonceStr,
		Body:           "xx",
		OutTradeNo:     "12312",
		TotalFee:       "1",
		SpbillCreateIP: "127.0.0.1",
		//TimeStart : time.Now().Format(util.SecondSeamlessDateFormat),
		//TimeExpire : time.Now().Add(time.Minute * 16).Format(util.SecondSeamlessDateFormat),
		NotifyUrl: notifyUrl,
		TradeType: "JSAPI",
		OpenID:    "ad23sd12",
	}

	apiKey := ""
	req.Sign = wx.GenSign(req, apiKey)
	res, err := wx.UnifiedOrder(&req)
	if err != nil {
		t.Log(err)
	}

	t.Log(res)
}

func TestJsCode2Session(t *testing.T) {
	code := "asdsdf"
	res, err := wx.JsCode2Session(code)
	if err != nil {
		t.Log(err)
	}

	t.Log(res)
}

func TestRefund(t *testing.T) {
	nonceStr, _ := encrypt.MD5(strconv.FormatInt(time.Now().Unix(), 10))
	req := RefundReq{
		MchID:       mchID,
		NonceStr:    nonceStr,
		OutTradeNo:  "123",
		OutRefundNo: strconv.Itoa(int(time.Now().Unix())),
		TotalFee:    1,
		RefundFee:   1,
		NotifyUrl:   refundNotifyUrl,
		//RefundDesc:"主动退款"
	}

	req.Sign = wx.GenSign(req, apiKey)
	res, err := wx.Refund(&req, certKey, cert)
	if err != nil {
		t.Log(err)
	}

	t.Log(res)
	fmt.Println(fmt.Sprintf("%+v", res))
	fmt.Println(res.Sign)
	fmt.Println(wx.GenSign(*res, apiKey))
}

func TestGenSign(t *testing.T) {
	refundXml := []byte(`<xml>
<return_code><![CDATA[SUCCESS]]></return_code>
<return_msg><![CDATA[OK]]></return_msg>
<appid><![CDATA[]]></appid>
<mch_id><![CDATA[]]></mch_id>
<nonce_str><![CDATA[xldNakMnWgQtaWfV]]></nonce_str>
<sign><![CDATA[]]></sign>
<result_code><![CDATA[SUCCESS]]></result_code>
<transaction_id><![CDATA[]]></transaction_id>
<out_trade_no><![CDATA[]]></out_trade_no>
<out_refund_no><![CDATA[]]></out_refund_no>
<refund_id><![CDATA[]]></refund_id>
<refund_channel><![CDATA[]]></refund_channel>
<refund_fee>1</refund_fee>
<coupon_refund_fee>0</coupon_refund_fee>
<total_fee>1</total_fee>
<cash_fee>1</cash_fee>
<coupon_refund_count>0</coupon_refund_count>
<cash_refund_fee>1</cash_refund_fee>
</xml>`)

	wx.ZeroValueProcess(refundXml)
	var data refundRes
	if xmlUnmarshalErr := xml.Unmarshal(refundXml, &data); xmlUnmarshalErr != nil {
		t.Error(xmlUnmarshalErr)
	}

	fmt.Println(wx.GenSign(data, apiKey))
}

func TestGetWxACodeUnLimit(t *testing.T) {

	//fmt.Println(url.QueryEscape("store_id=1#from=store_code"))
	//return

	accessTokenInfo, getTokenErr := wx.GetAccessToken()
	if getTokenErr != nil {
		t.Error(getTokenErr)
	}

	t.Log(accessTokenInfo)
	accessToken := accessTokenInfo.AccessToken

	req := GetWxACodeUnLimitReq{
		Scene: "store_id=1#from=store_code", // 不支持 & 符号，文档坑人
		Page:  "pages/home/home",
	}
	qrcodeInfo, getQrcodeErr := wx.GetWxACodeUnLimit(accessToken, &req)
	if getQrcodeErr != nil {
		t.Error(getQrcodeErr)
	}

	//t.Log(qrcodeInfo)

	// remote save
	conf := aliyun.OssConfig{
		AccessKeyID:     "",
		AccessKeySecret: "",
		Endpoint:        "",
		Bucket:          "",
	}
	oss, _ := aliyun.NewOssStore(&conf)
	oss.UploadBytes(qrcodeInfo.Buffer, "test/qrcode.jpeg")

	// local save
	//f, _ := os.Create("./qrcode-10228.jpeg")
	//f.Write(qrcodeInfo.Buffer)
	//f.Close()

}

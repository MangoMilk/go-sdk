package wechat

import (
	"github.com/MangoMilk/go-kit/encrypt"
	"strconv"
	"testing"
	"time"
)

var (
	wx *Wechat
)

func setup() {
	appID := ""
	appSecret := ""
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
		MchID:          "",
		NonceStr:       nonceStr,
		Body:           "xx",
		OutTradeNo:     "12312",
		TotalFee:       "1",
		SpbillCreateIP: "127.0.0.1",
		//TimeStart : time.Now().Format(util.SecondSeamlessDateFormat),
		//TimeExpire : time.Now().Add(time.Minute * 16).Format(util.SecondSeamlessDateFormat),
		NotifyUrl: "https://127.0.0.1/xx/notify",
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
		MchID:       "",
		NonceStr:    nonceStr,
		OutTradeNo:  "12312",
		OutRefundNo: "12312",
		TotalFee:    1,
		RefundFee:   1,
		NotifyUrl:   "https://127.0.0.1/xx/notify",
		//RefundDesc:"主动退款"
	}

	apiKey := ""
	req.Sign = wx.GenSign(req, apiKey)
	certKey := ""
	cert := ""
	res, err := wx.Refund(&req, certKey, cert)
	if err != nil {
		t.Log(err)
	}

	t.Log(res)
}

package mi

import (
	"encoding/json"
	"sort"
	"strings"
	"github.com/MangoMilk/go-kit/net"
	"github.com/MangoMilk/go-kit/encrypt"
)

type CheckSessionRes struct {
	ErrCode float64 `json:"errcode"`
	ErrMsg  string  `json:"errMsg"`
	Adult   float64 `json:"adult"`
}

const (
	authUrl = "https://mis.migc.xiaomi.com/api/biz/service/loginvalidate"
)

func CheckSessionID(appID string, sessionID string, uid string, appSecret string) (*CheckSessionRes,error) {
	var postData = make(map[string]string)
	postData["appId"] = appID
	postData["session"] = sessionID
	postData["uid"] = uid
	signature,signErr := genXMSignature(postData, appSecret)
	if signErr !=nil{
		return nil,signErr
	}

	postData["signature"] = signature

	httpRes,err := net.HttpPost(authUrl, postData, nil, false,"","")
	if err != nil {
		return nil,err
	}
	var res CheckSessionRes

	if jsonErr:=json.Unmarshal(httpRes,&res);jsonErr != nil {
		return nil,jsonErr
	}

	return &res,nil
}

func genXMSignature(data map[string]string, appSecret string) (string,error) {
	// 生成签名（sign）
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

	// hmac-sha1 && key = app_secret
	return encrypt.HmacSHA1(strings.Trim(signStr, "&"), appSecret)
}

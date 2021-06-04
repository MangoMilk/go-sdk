package wechat

import (
	"encoding/json"
	"github.com/MangoMilk/go-kit/net"
)

type WeChatAuthRes struct {
	ErrCode    float64 `json:"errcode"`
	ErrMsg     string  `json:"errmsg"`
	SessionKey string  `json:"session_key"`
	OpenID     string  `json:"openid"`
	UnionID    string  `json:"unionid"`
}

const (
	AuthApi = "https://api.weixin.qq.com/sns/jscode2session"
)

func Auth(appID string, appSecret string, code string) (*WeChatAuthRes,error) {
	queryParam := "?appid=" + appID + "&secret=" + appSecret + "&js_code=" + code + "&grant_type=authorization_code"
	res,httpErr:=net.HttpGet(AuthApi+queryParam, nil)
	if httpErr != nil {
		return nil,httpErr
	}

	var data WeChatAuthRes

	if jsonErr := json.Unmarshal(res, &data);jsonErr != nil {
		return nil,jsonErr
	}

	return &data,nil
}

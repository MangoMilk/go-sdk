package qcloud

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"
)

var COS_DOMAIN string
var COS_BUCKET string
var COS_REGION string

/**
 * InitCOS
 */
func InitCOS(bucket string, region string) {
	COS_BUCKET = bucket
	COS_REGION = region
	COS_DOMAIN = COS_BUCKET + ".cos." + COS_REGION + ".myqcloud.com"
}

/**
 * GetObject
 */
//func GetCosObject(secretKey string, secretID string, uri string, data map[string]interface{}, header map[string]string) interface{} {
//	var api string = "http://" + COS_DOMAIN + uri
//
//	sign := CosSign(secretKey, secretID, uri, "get", data, header)
//
//	if header == nil {
//		header = make(map[string]string)
//	}
//	header["Authorization"] = sign
//
//	return util.JsonDecoded(modules.HttpGet(api, header))
//}

/**
 * CosSign
 */
func CosSign(secretKey string, secretID string, uri string, method string, data map[string]interface{}, header map[string]string) string {
	// 1.KeyTime
	var StartTimestamp int64 = time.Now().Unix()
	var EndTimestamp int64 = StartTimestamp + 300
	var KeyTime string = strconv.Itoa(int(StartTimestamp)) + ";" + strconv.Itoa(int(EndTimestamp))

	// 2.SignKey
	// alg = HS1
	var secret string = secretKey
	secretByte := []byte(secret)
	h := hmac.New(sha1.New, secretByte)
	h.Write([]byte(KeyTime))
	var SignKey string = hex.EncodeToString(h.Sum(nil))

	// 3.UrlParamList HttpParameters
	// sort data keys
	var dataKeys []string
	for k, _ := range data {
		dataKeys = append(dataKeys, k)
	}
	sort.Strings(dataKeys)

	// joint data
	var UrlParamList string = ""
	var HttpParameters string = ""
	for _, v := range dataKeys {
		val := data[v].(string)
		v = strings.ToLower(v)
		UrlParamList += (v + ";")
		v = url.QueryEscape(val)
		HttpParameters += (v + "=" + val + "&")
	}
	UrlParamList = strings.TrimRight(UrlParamList, ";")
	HttpParameters = strings.TrimRight(HttpParameters, "&")

	// 4.HeaderList HttpHeaders
	// sort header keys
	var headerKeys []string
	for k, _ := range header {
		headerKeys = append(headerKeys, k)
	}
	sort.Strings(headerKeys)

	// joint header
	var HeaderList string = ""
	var HttpHeaders string = ""
	for _, v := range headerKeys {
		val := header[v]
		v = strings.ToLower(v)
		HeaderList += (v + ";")
		v = url.QueryEscape(val)
		HttpHeaders += (v + "=" + val + "&")
	}
	HeaderList = strings.TrimRight(HeaderList, ";")
	HttpHeaders = strings.TrimRight(HttpHeaders, "&")

	// 5.HttpString
	var HttpMethod string = strings.ToLower(method)
	var UriPathname string = uri

	var HttpString string = HttpMethod + "\n" + UriPathname + "\n" + HttpParameters + "\n" + HttpHeaders + "\n"

	// 6.StringToSign
	// alg SHA1
	s := sha1.New()
	s.Write([]byte(HttpString))
	var SHA1HttpString string = hex.EncodeToString(s.Sum(nil))
	var StringToSign string = "sha1\n" + KeyTime + "\n" + SHA1HttpString + "\n"

	// 7.Signature
	// alg HS1
	var hSecretByte = []byte(SignKey)
	hh := hmac.New(sha1.New, hSecretByte)
	hh.Write([]byte(StringToSign))
	var Signature string = hex.EncodeToString(hh.Sum(nil)) //hash_hmac('sha1', $StringToSign, $SignKey);

	// 8.生成
	var SecretId string = secretID
	var sign string = "q-sign-algorithm=sha1"
	sign += "&q-ak=" + SecretId
	sign += "&q-sign-time=" + KeyTime
	sign += "&q-key-time=" + KeyTime
	sign += "&q-header-list=" + HeaderList
	sign += "&q-url-param-list=" + UrlParamList
	sign += "&q-signature=" + Signature

	return sign
}

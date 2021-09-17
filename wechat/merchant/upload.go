package merchant

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/encrypt"
	"github.com/MangoMilk/go-kit/net"
	"strings"
	"time"
)

// ==================== 上传图片 ====================
type UploadReq struct {
	File bytes.Buffer `json:"file"`
	/*图片文件		是	body
	将媒体图片进行二进制转换，得到的媒体图片二进制内容，在请求body中上传此二进制内容。
	媒体图片只支持JPG、BMP、PNG格式，文件大小不能超过2M。
	示例值：pic1

	视频文件 是 body
	将媒体视频进行二进制转换，得到的媒体视频二进制内容，在请求body中上传此二进制内容。
	媒体视频只支持avi、wmv、mpeg、mp4、mov、mkv、flv、f4v、m4v、rmvb格式，文件大小不能超过5M。
	示例值：pic1
	*/
	Meta metaData `json:"meta"` //媒体文件元信息	是	body 媒体文件元信息，使用json表示，包含两个对象：filename、sha256。
}

type metaData struct {
	Filename string `json:"filename"`
	/*文件名称	[1,128]	是
	商户上传的媒体图片的名称，商户自定义，必须以JPG、BMP、PNG为后缀。
	示例值：filea.jpg

	视频名称	[1,128]	是
	商户上传的媒体视频的名称，商户自定义，必须以avi、wmv、mpeg、mp4、mov、mkv、flv、f4v、m4v、rmvb为后缀。
	示例值：file_test.mp4
	*/
	Sha256 string `json:"sha256"` //文件摘要	[1,64]	是	图片/视频文件的文件摘要，即对图片/视频文件的二进制内容进行sha256计算得到的值。示例值：hjkahkjsjkfsjk78687dhjahdajhk
}

type uploadRes struct {
	MediaID string `json:"media_id"`
	/*媒体文件标识 Id [1,512]	是	微信返回的媒体文件标识Id。
	示例值：6uqyGjGrCf2GtyXP8bxrbuH9-aAoTjH-rKeSl3Lf4_So6kdkQu4w8BYVP3bzLtvR38lxt4PjtCDXsQpzqge_hQEovHzOhsLleGFQVRF-U_0
	*/
}

func UploadImage(req *UploadReq, mchID, serialNo, sign string) (*uploadRes, error) {

	api := "https://api.mch.weixin.qq.com/v3/merchant/media/upload"

	t := time.Now()
	randStr, _ := encrypt.MD5(fmt.Sprintf("%v", t.UnixNano()))
	boundary := randStr
	nonceStr, _ := encrypt.MD5(api + randStr)

	header := map[string]string{
		"Authorization": fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%v\",nonce_str=\"%v\",timestamp=\"%v\",serial_no=\"%v\",signature=\"%v\"", mchID, nonceStr, t.Unix(), serialNo, sign),
		"Content-Type":  fmt.Sprintf("multipart/form-data;boundary=%v", boundary),
		"Accept":        "application/json",
	}

	fileNameArr := strings.Split(req.Meta.Filename, ".")
	fileType := fileNameArr[len(fileNameArr)-1]

	body := fmt.Sprintf(`
--%v  
Content-Disposition: form-data; name="meta";
Content-Type: application/json

{ "filename": "%v", "sha256": "%v" }
--boundary
Content-Disposition: form-data; name="file"; filename="%v";
Content-Type: image/%v

%v
--%v--
`, boundary, req.Meta.Filename, req.Meta.Sha256, req.Meta.Filename, fileType, req.File, boundary)

	res, httpErr := net.HttpPost(api, body, header, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================UploadImage")
	fmt.Println(string(res))

	var data uploadRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== 上传视频 ====================
func UploadVideo(req *UploadReq, mchID, serialNo, sign string) (*uploadRes, error) {
	api := "https://api.mch.weixin.qq.com/v3/merchant/media/video_upload"

	t := time.Now()
	randStr, _ := encrypt.MD5(fmt.Sprintf("%v", t.UnixNano()))
	boundary := randStr
	nonceStr, _ := encrypt.MD5(api + randStr)

	header := map[string]string{
		"Authorization": fmt.Sprintf("WECHATPAY2-SHA256-RSA2048 mchid=\"%v\",nonce_str=\"%v\",timestamp=\"%v\",serial_no=\"%v\",signature=\"%v\"", mchID, nonceStr, t.Unix(), serialNo, sign),
		"Content-Type":  fmt.Sprintf("multipart/form-data;boundary=%v", boundary),
		"Accept":        "application/json",
	}

	fileNameArr := strings.Split(req.Meta.Filename, ".")
	fileType := fileNameArr[len(fileNameArr)-1]

	body := fmt.Sprintf(`
--%v  
Content-Disposition: form-data; name="meta";
Content-Type: application/json

{ "filename": "%v", "sha256": "%v" }
--boundary
Content-Disposition: form-data; name="file"; filename="%v";
Content-Type: video/%v

%v
--%v--
`, boundary, req.Meta.Filename, req.Meta.Sha256, req.Meta.Filename, fileType, req.File, boundary)

	res, httpErr := net.HttpPost(api, body, header, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("====================UploadVideo")
	fmt.Println(string(res))

	var data uploadRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

package ecommerce

import (
	"encoding/json"
	"fmt"
	"github.com/MangoMilk/go-kit/net"
)

// ==================== 进件 ====================
type OrganizationType string

const (
	OrganizationTypeMicroStore   = OrganizationType("2401")
	OrganizationTypePersonSeller = OrganizationType("2500")
	OrganizationTypeIICH         = OrganizationType("4") // Individual industrial and commercial households
	OrganizationTypeCompany      = OrganizationType("2")
	OrganizationTypeGovernment   = OrganizationType("3")
	OrganizationTypeOther        = OrganizationType("1708")
)

type ApplyReq struct {
	OutRequestNo string `json:"out_request_no"`
	/* 必填，业务申请编号，长度 1~124
	1、服务商自定义的商户唯一编号。
	2、每个编号对应一个申请单，每个申请单审核通过后会生成一个微信支付商户号。
	3、若申请单被驳回，可填写相同的“业务申请编号”，即可覆盖修改原申请单信息 。
	示例值：APPLYMENT_00000000001
	*/
	OrganizationType OrganizationType `json:"organization_type"`
	/* 必填，主体类型，长度 1~4
	非小微的主体类型需与营业执照/登记证书上一致，可参考选择主体指引，枚举值如下。
	2401：小微商户，指无营业执照的个人商家。
	2500：个人卖家，指无营业执照，已持续从事电子商务经营活动满6个月，且期间经营收入累计超过20万元的个人商家。（若选择该主体，请在“补充说明”填写相关描述）
	4：个体工商户，营业执照上的主体类型一般为个体户、个体工商户、个体经营。
	2：企业，营业执照上的主体类型一般为有限公司、有限责任公司。
	3：党政、机关及事业单位，包括国内各级、各类政府机构、事业单位等（如：公安、党团、司法、交通、旅游、工商税务、市政、医疗、教育、学校等机构）。
	1708：其他组织，不属于企业、政府/事业单位的组织机构（如社会团体、民办非企业、基金会），要求机构已办理组织机构代码证。
	示例值：2401
	*/
	BusinessLicenseInfo businessLicenseInfo `json:"business_license_info"`
	/* 条件选填，营业执照/登记证书信息
	1、主体为“小微/个人卖家”时，不填。
	2、主体为“个体工商户/企业”时，请上传营业执照。
	3、主体为“党政、机关及事业单位/其他组织”时，请上传登记证书。
	*/
	OrganizationCertInfo organizationCertInfo `json:"organization_cert_info"`
	/* 条件选填，组织机构代码证信息
	主体为企业/党政、机关及事业单位/其他组织，且证件号码不是18位时必填。
	注：
	若营业执照未三证合一 ，该参数必传;
	若营业执照三证合一 ，该参数可不传。
	*/
	IdDocType string `json:"id_doc_type"`
	/* 否，经营者/法人证件类型，长度1~64
	1、主体为“小微/个人卖家”，可选择：身份证。
	2、主体为“个体户/企业/党政、机关及事业单位/其他组织”，可选择：以下任一证件类型。
	3、若没有填写，系统默认选择：身份证。
	枚举值:
	IDENTIFICATION_TYPE_MAINLAND_IDCARD：中国大陆居民-身份证
	IDENTIFICATION_TYPE_OVERSEA_PASSPORT：其他国家或地区居民-护照
	IDENTIFICATION_TYPE_HONGKONG：中国香港居民–来往内地通行证
	IDENTIFICATION_TYPE_MACAO：中国澳门居民–来往内地通行证
	IDENTIFICATION_TYPE_TAIWAN：中国台湾居民–来往大陆通行证
	示例值：IDENTIFICATION_TYPE_MACAO
	*/
	IdCardInfo idCardInfo `json:"id_card_info"`
	/* 条件选填，经营者/法人身份证信息
	请填写经营者/法人的身份证信息
	证件类型为“身份证”时填写。
	*/
	IdDocInfo idDocInfo `json:"id_doc_info"`
	/* 条件选填，经营者/法人其他类型证件信息
	证件类型为“来往内地通行证、来往大陆通行证、护照”时填写。
	*/
	NeedAccountInfo bool `json:"need_account_info"`
	/* 是  是否填写结算银行账户
	1、可根据实际情况，填写“true”或“false”。
	   1）若为“true”，则需填写结算银行账户。
	   2）若为“false”，则无需填写结算银行账户。
	2、若入驻时未填写结算银行账户，则需入驻后调用修改结算账户API补充信息，才能发起提现。
	3、当超级管理员类型为负责人时，该字段只能传true，即结算银行账户必填
	示例值：true
	*/
	AccountInfo accountInfo `json:"account_info"`
	/* 条件选填    结算银行账户
	   若"是否填写结算账户信息"填写为“true”, 则必填，填写为“false”不填 。
	*/
	ContactInfo contactInfo `json:"contact_info"`
	/* 是 超级管理员信息
	   请填写店铺的超级管理员信息。
	   超级管理员需在开户后进行签约，并可接收日常重要管理信息和进行资金操作，请确定其为商户法定代表人或负责人。
	*/
	SalesSceneInfo    salesSceneInfo `json:"sales_scene_info"` // 是  店铺信息 请填写店铺信息
	MerchantShortname string         `json:"merchant_shortname"`
	/* 是   商户简称 [1,64]
	UTF-8格式，中文占3个字节，即最多21个汉字长度。将在支付完成页向买家展示，需与商家的实际售卖商品相符 。
	示例值：腾讯
	*/
	Qualifications string `json:"qualifications"`
	/* 否 特殊资质  [1,1024]
	1、根据商户经营业务要求提供相关资质，详情查看《行业对应特殊资质》。
	2、请提供为“申请商家主体”所属的特殊资质，可授权使用总公司/分公司的特殊资 质；
	3、最多可上传5张照片，请填写通过图片上传接口预先上传图片生成好的MediaID 。
	示例值：[\"jTpGmxUX3FBWVQ5NJInE4d2I6_H7I4\"]
	*/
	BusinessAdditionPics string `json:"business_addition_pics"`
	/*  否   补充材料  [1,1024]
	根据实际审核情况，额外要求提供。最多可上传5张照片，请填写通过图片上传接口预先上传图片生成好的MediaID 。
	示例值：[\"jTpGmg05InE4d2I6_H7I4\"]
	*/
	BusinessAdditionDesc string `json:"business_addition_desc"`
	/* 否   补充说明   [1,256]
	1、可填写512字以内 。
	2、若主体为“个人卖家”，该字段必传，则需填写描述“ 该商户已持续从事电子商务经营活动满6个月，且期间经营收入累计超过20万元。”
	示例值：特殊情况，说明原因
	*/
}

type businessLicenseInfo struct {
	BusinessLicenseCopy string `json:"business_license_copy"` /* 必填，证件扫描件，长度 1~256
		1、主体为“个体工商户/企业”时，请上传营业执照的证件图片。
	2、主体为“党政、机关及事业单位/其他组织”时，请上传登记证书的证件图片。
	3、可上传1张图片，请填写通过图片上传接口预先上传图片生成好的MediaID 。
	4、图片要求：
	（1）请上传证件的彩色扫描件或彩色数码拍摄件，黑白复印件需加盖公章（公章信息需完整） 。
	（2）不得添加无关水印（非微信支付商户申请用途的其他水印）。
	（3）需提供证件的正面拍摄件，完整、照面信息清晰可见。信息不清晰、扭曲、压缩变形、反光、不完整均不接受。
	（4）不接受二次剪裁、翻拍、PS的证件照片。
	示例值： 47ZC6GC-vnrbEny_Ie_An5-tCpqxucuxi-vByf3Gjm7KE53JXvGy9tqZm2XAUf-4KGprrKhpVBDIUv0OF4wFNIO4kqg05InE4d2I6_H7I4
	*/
	BusinessLicenseNumber string `json:"business_license_number"`
	/*是，证件注册号，长度 15~ 18
	1、主体为“个体工商户/企业”时，请填写营业执照上的注册号/统一社会信用代码，须为15位数字或 18位数字|大写字母。
	2、主体为“党政、机关及事业单位/其他组织”时，请填写登记证书的证书编号。
	示例值：123456789012345678
	商户名称	merchant_name	string[1,128]	是	1、请填写营业执照/登记证书的商家名称，2~110个字符，支持括号 。
	2、个体工商户/党政、机关及事业单位，不能以“公司”结尾。
	3、个体工商户，若营业执照上商户名称为空或为“无”，请填写"个体户+经营者姓名"，如“个体户张三” 。
	示例值：腾讯科技有限公司
	*/
	LegalPerson    string `json:"legal_person"`    //是，经营者/法定代表人姓名，长度 1~128	请填写证件的经营者/法定代表人姓名。示例值：张三
	CompanyAddress string `json:"company_address"` //条件选填，注册地址，长度 1~128	主体为“党政、机关及事业单位/其他组织”时必填，请填写登记证书的注册地址。	示例值：深圳南山区科苑路
	BusinessTime   string `json:"business_time"`
	/* 条件选填，营业期限，长度 1~256
	1、主体为“党政、机关及事业单位/其他组织”时必填，请填写证件有效期。
	2、若证件有效期为长期，请填写：长期。
	3、结束时间需大于开始时间。
	4、有效期必须大于60天，即结束时间距当前时间需超过60天。
	示例值：[\"2014-01-01\",\"长期\"]
	*/
}

type organizationCertInfo struct {
	OrganizationCopy string `json:"organization_copy"`
	/*是，组织机构代码证照片，长度 1~256
	可上传1张图片，请填写通过图片上传接口预先上传图片生成好的MediaID。
	示例值：vByf3Gjm7KE53JXv\prrKhpVBDIUv0OF4wFNIO4kqg05InE4d2I6_H7I4
	*/
	OrganizationNumber string `json:"organization_number"`
	/* 是，组织机构代码，长度 1~256
	1、请填写组织机构代码证上的组织机构代码。
	2、可填写9或10位 数字|字母|连字符。
	示例值：12345679-A
	*/
	OrganizationTime string `json:"organization_time"`
	/* 是，组织机构代码有效期限	，1~256
	1、请填写组织机构代码证的有效期限，注意参照示例中的格式。
	2、若证件有效期为长期，请填写：长期。
	3、结束时间需大于开始时间。
	4、有效期必须大于60天，即结束时间距当前时间需超过60天。
	示例值：[\"2014-01-01\",\"长期\"]
	*/
}

type idCardInfo struct {
	IdCardCopy string `json:"id_card_copy"`
	/* 身份证人像面照片	[1,256]	是
	1、请上传经营者/法定代表人的身份证人像面照片。
	2、可上传1张图片，请填写通过图片上传接口预先上传图片生成好的MediaID。
	示例值：xpnFuAxhBTEO_PvWkfSCJ3zVIn001D8daLC-ehEuo0BJqRTvDujqhThn4ReFxikqJ5YW6zFQ
	*/
	IdCardNational string `json:"id_card_national"`
	/* 身份证国徽面照片，[1,256]	是
	1、请上传经营者/法定代表人的身份证国徽面照片。
	2、可上传1张图片，请填写通过图片上传接口预先上传图片生成好的MediaID 。
	示例值：vByf3Gjm7KE53JXvGy9tqZm2XAUf-4KGprrKhpVBDIUv0OF4wFNIO4kqg05InE4d2I6_H7I4
	*/
	IdCardName string `json:"id_card_name"`
	/* 身份证姓名，[1,256]	是
	1、请填写经营者/法定代表人对应身份证的姓名，2~30个中文字符、英文字符、符号。
	2、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：pVd1HJ6v/69bDnuC4EL5Kz4jBHLiCa8MRtelw/wDa4SzfeespQO/0kjiwfqdfg==
	*/
	IdCardNumber string `json:"id_card_number"`
	/* 身份证号码	[15,18]	是
	1、请填写经营者/法定代表人对应身份证的号码。
	2、15位数字或17位数字+1位数字|X ，该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：zV+BEmytMNQCqQ8juwEc4P4TG5xzchG/5IL9DBd+Z0zZXkw==4
	*/
	IdCardValidTime string `json:"id_card_valid_time"`
	/* 身份证有效期限	，[1,128]	是
	1、请填写身份证有效期的结束时间，注意参照示例中的格式。
	2、若证件有效期为长期，请填写：长期。
	3、证件有效期需大于60天。
	示例值：2026-06-06
	*/
}

type idDocInfo struct {
	IdDocName string `json:"id_doc_name"`
	/* 证件姓名	[1,128]	是
	1、请填写经营者/法人姓名。
	2、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：jTpGmxUX3FBWVQ5NJTZvlKX_gdU4LC-ehEuo0BJqRTvDujqhThn4ReFxikqJ5YW6zFQ
	*/
	IdDocNumber string `json:"id_doc_number"`
	/* 证件号码	[1,128]	是
	7~11位 数字|字母|连字符 。
	该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：jTpGmxUX3FBWVQ5NJTZvlKX_go0BJqRTvDujqhThn4ReFxikqJ5YW6zFQ
	*/
	IdDocCopy string `json:"id_doc_copy"`
	/* 证件照片	[1,256]	是
	1、可上传1张图片，请填写通过图片上传接口预先上传图片生成好的MediaID。
	2、2M内的彩色图片，格式可为bmp、png、jpeg、jpg或gif 。
	示例值：xi-vByf3Gjm7KE53JXvGy9tqZm2XAUf-4KGprrKhpVBDIUv0OF4wFNIO4kqg05InE4d2I6_H7I4
	*/
	DocPeriodEnd string `json:"doc_period_end"`
	/* 证件结束日期	 [1,128]	是
	1、请按照示例值填写。
	2、若证件有效期为长期，请填写：长期。
	3、证件有效期需大于60天 。
	示例值：2020-01-02
	*/
}

type accountInfo struct {
	BankAccountType string `json:"bank_account_type"`
	/* 账户类型	[1,2]	是
	1、若主体为企业/党政、机关及事业单位/其他组织，可填写：74-对公账户。
	2、主体为“小微/个人卖家”，可选择：75-对私账户。
	3、若主体为个体工商户，可填写：74-对公账户、75-对私账户。
	示例值：75
	*/
	AccountBank string `json:"account_bank"`
	/* 开户银行	 [1,10]	是
	详细参见开户银行对照表。
	注：
	17家直连银行，请根据开户银行对照表直接填写银行名 ;
	非17家直连银行，该参数请填写为“其他银行”。
	示例值：工商银行
	*/
	AccountName string `json:"account_name"`
	/* 开户名称	[1,128]	是
	1、选择经营者个人银行卡时，开户名称必须与身份证姓名一致。
	2、选择对公账户时，开户名称必须与营业执照上的“商户名称”一致。
	3、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：AOZdYGISxo4yw96uY1Pk7Rq79Jtt7+I8juwEc4P4TG5xzchG/5IL9DBd+Z0zZXkw==
	*/
	BankAddressCode string `json:"bank_address_code"`
	/* 开户银行省市编码	[1,12]	是
	至少精确到市，详细参见省市区编号对照表。
	注：
	仅当省市区编号对照表中无对应的省市区编号时，可向上取该银行对应市级编号或省级编号。
	示例值：110000
	*/
	BankBranchID string `json:"bank_branch_id"`
	/* 开户银行联行号	 [1,64]	条件选填
	1、17家直连银行无需填写，如为其他银行，开户银行全称（含支行）和开户银行联行号二选一。
	2、详细参见开户银行全称（含支行）对照表。
	示例值：402713354941
	*/
	BankName string `json:"bank_name"`
	/* 开户银行全称 （含支行）[1,128]	条件选填
	1、17家直连银行无需填写，如为其他银行，开户银行全称（含支行）和开户银行联行号二选一。
	2、需填写银行全称，如"深圳农村商业银行XXX支行" 。
	3、详细参见开户银行全称（含支行）对照表。
	示例值：施秉县农村信用合作联社城关信用社
	*/
	AccountNumber string `json:"account_number"`
	/* 银行账号	[1,128]	是
	1、数字，长度遵循系统支持的对公/对私卡号长度要求表。
	2、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值： d+xT+MQCvrLHUVDWv/8MR/dB7TkXLVfSrUxMPZy6jWWYzpRrEEaYQE8ZRGYoeorwC+w==
	*/
}

type contactInfo struct {
	ContactType string `json:"contact_type"`
	/* 超级管理员类型	[1,2]	是
	1、主体为“小微/个人卖家 ”，可选择：65-经营者/法人。
	2、主体为“个体工商户/企业/党政、机关及事业单位/其他组织”，可选择：65-经营者/法人、66- 负责人。 （负责人：经商户授权办理微信支付业务的人员，授权范围包括但不限于签约，入驻过程需完成账户验证）。
	示例值：65
	*/
	ContactName string `json:"contact_name"`
	/* 超级管理员姓名 [1,256]	是
	1、若管理员类型为“法人”，则该姓名需与法人身份证姓名一致。
	2、若管理员类型为“负责人”，则可填写实际负责人的姓名。
	3、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	（后续该管理员需使用实名微信号完成签约）
	示例值： pVd1HJ6zyvPedzGaV+X3IdGdbDnuC4Eelw/wDa4SzfeespQO/0kjiwfqdfg==
	*/
	ContactIdCardNumber string `json:"contact_id_card_number"`
	/* 超级管理员身份证件号码 [1,256]	是
	1、若管理员类型为法人，则该身份证号码需与法人身份证号码一致。若管理员类型为负责人，则可填写实际负责人的身份证号码。
	2、可传身份证、来往内地通行证、来往大陆通行证、护照等证件号码。
	3、超级管理员签约时，校验微信号绑定的银行卡实名信息，是否与该证件号码一致。
	4、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：pVd1HJ6zmty7/mYNxLMpRSvMRtelw/wDa4SzfeespQO/0kjiwfqdfg==
	*/
	MobilePhone string `json:"mobile_phone"`
	/* 超级管理员手机	 [1,256]	是
	1、请填写管理员的手机号，11位数字， 用于接收微信支付的重要管理信息及日常操作验证码 。
	2、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：pVd1HJ6zyvPedzGaV+X3qtmrq9bb9tPROvwia4ibL+F6mfjbzQIzfb3HHLEjZ4YiNWWNeespQO/0kjiwfqdfg==
	*/
	ContactEmail string `json:"contact_email"`
	/* 超级管理员邮箱[1,256]	条件选填
	1、主体类型为“小微商户/个人卖家”可选填，其他主体需必填。
	2、用于接收微信支付的开户邮件及日常业务通知。
	3、需要带@，遵循邮箱格式校验 。
	4、该字段需进行加密处理，加密方法详见敏感信息加密说明。(提醒：必须在HTTP头中上送Wechatpay-Serial)
	示例值：pVd1HJ6zyvPedzGaV+X3qtmrq9bb9tPROvwia4ibL+FWWNUlw/wDa4SzfeespQO/0kjiwfqdfg==
	*/
}

type salesSceneInfo struct {
	StoreName string `json:"store_name"` //店铺名称	 [1,256]	是	请填写店铺全称。	示例值：爱烧烤

	StoreUrl string `json:"store_url"`
	/* 店铺链接	[1,1024]	二选一
	1、店铺二维码or店铺链接二选一必填。
	2、请填写店铺主页链接，需符合网站规范。
	示例值：http://www.qq.com
	*/
	StoreQrCode string `json:"store_qr_code"`
	/* 店铺二维码[1,256]
	1、店铺二维码 or 店铺链接二选一必填。
	2、若为电商小程序，可上传店铺页面的小程序二维码。
	3、请填写通过图片上传接口预先上传图片生成好的MediaID，仅能上传1张图片 。
	示例值：jTpGmxUX3FBWVQ5NJTZvlKX_gdU4cRz7z5NxpnFuAxhBTEO1D8daLC-ehEuo0BJqRTvDujqhThn4ReFxikqJ5YW6zFQ
	*/
	MiniProgramSubAppID string `json:"mini_program_sub_appid"`
	/* 小程序AppID [1,256]	否
	1、商户自定义字段，可填写已认证的小程序AppID，认证主体需与二级商户主体一致；
	2、完成入驻后， 系统发起二级商户号与该AppID的绑定（即配置为sub_appid，可在发起支付时传入）
	示例值：wxd678efh567hg6787
	*/
}

type applyRes struct {
	ApplymentID  uint64 `json:"applyment_id"` // 微信支付申请单号		是	微信支付分配的申请单号 。示例值：2000002124775691
	OutRequestNo string `json:"out_request_no"`
	/*业务申请编号	[1,124]	是
	服务商自定义的商户唯一编号。每个编号对应一个申请单，每个申请单审核通过后会生成一个微信支付商户号。
	示例值：APPLYMENT_00000000001
	*/
}

func Apply(req *ApplyReq) (*applyRes, error) {
	api := "https://api.mch.weixin.qq.com/v3/ecommerce/applyments/"
	res, httpErr := net.HttpPost(api, *req, nil, false, "", "")
	if httpErr != nil {
		return nil, httpErr
	}

	fmt.Println("=================Apply===============Apply")
	fmt.Println(string(res))

	var data applyRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		//data.Buffer = res
	}

	return &data, nil
}

// ==================== 查询进件状态 ====================
type ApplymentState string

const (
	ApplymentStateChecking          = ApplymentState("CHECKING")
	ApplymentStateAccountNeedVerify = ApplymentState("ACCOUNT_NEED_VERIFY")
	ApplymentStateAuditing          = ApplymentState("AUDITING")
	ApplymentStateRejected          = ApplymentState("REJECTED")
	ApplymentStateNeedSign          = ApplymentState("NEED_SIGN")
	ApplymentStateFinish            = ApplymentState("FINISH")
	ApplymentStateFrozen            = ApplymentState("FROZEN")
)

type SignState string

const (
	SignStateUnsigned    = SignState("UNSIGNED")
	SignStateSigned      = SignState("SIGNED")
	SignStateNotSignable = SignState("NOT_SIGNABLE")
)

type getApplyStatusRes struct {
	ApplymentState ApplymentState `json:"applyment_state"`
	/* 申请状态	[1,32]	是	枚举值：
	CHECKING：资料校验中
	ACCOUNT_NEED_VERIFY：待账户验证
	AUDITING：审核中
	REJECTED：已驳回
	NEED_SIGN：待签约
	FINISH：完成
	FROZEN：已冻结
	示例值：FINISH
	申请状态描述	applyment_state_desc	string[1,32]	是	申请状态描述
	示例值：“审核中”
	*/
	SignState SignState `json:"sign_state"`
	/* 签约状态	[1,16]	否
	1、UNSIGNED：未签约。该状态下，电商平台可查询获取签约链接，引导二级商户的超级管理员完成签约；
	2、SIGNED ：已签约。指二级商户的超级管理员已完成签约。注意：若申请单被驳回，商户修改了商户主体名称、法人名称、超级管理员信息、主体类型等信息，则需重新签约。
	3、NOT_SIGNABLE：不可签约。该状态下，暂不支持超级管理员签约。一般为申请单处于已驳回、已冻结、机器校验中状态，无法签约。
	示例值：SIGNED
	*/
	SignUrl string `json:"sign_url"`
	/* 签约链接 [1,256]	否
	1、当申请状态为NEED_SIGN 或 签约状态为UNSIGNED时返回，该链接为永久有效；
	2、申请单中的超级管理者，需用已实名认证的微信扫码打开，完成签约。
	示例值：https://pay.weixin.qq.com/public/apply4ec_sign/s?applymentId=2000002126198476&sign=b207b673049a32c858f3aabd7d27c7ec
	*/
	SubMchid           string            `json:"sub_mchid"`          // 电商平台二级商户号[1,32]	否	当申请状态为NEED_SIGN或FINISH时才返回。	示例值：1542488631
	AccountValidation  accountValidation `json:"account_validation"` // 汇款账户验证信息		否	当申请状态为ACCOUNT_NEED_VERIFY 时有返回，可根据指引汇款，完成账户验证。
	AuditDetail        []auditDetail     `json:"audit_detail"`       // 驳回原因详情	否	各项资料的审核情况。当申请状态为REJECTED或 FROZEN时才返回。
	LegalValidationUrl string            `json:"legal_validation_url"`
	/* 法人验证链接	[1,256]	否
	1、当申请状态为
	ACCOUNT_NEED_VERIFY，且通过系统校验的申请单，将返回链接。
	2、建议将链接转为二维码展示，让商户法人用微信扫码打开，完成账户验证。
	注：商户申请单进入审核状态后，微信侧会校验法人证件号码是否跟营业执照匹配，或匹配，返回该字段 ; 若不匹配，不支持法人扫码验证，不返回该字段。
	示例值： https://pay.weixin.qq.com/public/apply4ec_sign/s?applymentId=2000002126198476&sign=b207b673049a32c858f3aabd7d27c7ec
	*/
	OutRequestNo string `json:"out_request_no"` // 业务申请编号[1,124]	是	提交接口填写的业务申请编号。	示例值：APPLYMENT_00000000001
	ApplymentID  uint64 `json:"applyment_id"`   // 微信支付申请单号	是	微信支付分配的申请单号。	示例值：2000002124775691
}

type accountValidation struct {
	AccountName string `json:"account_name"`
	/* 付款户名	[1,128]	是
	需商户使用该户名的账户进行汇款。
	该字段需进行解密处理，解密方法详见敏感信息加解密说明。
	示例值： rDdICA3ZYXshYqeOSslSjSMf+MhhC4oaujiISFzq3AE+as7mAEDJly+DgRuVs74msmKUH8pl+3oA==
	*/
	AccountNo string `json:"account_no"`
	/* 付款卡号 [1,128]	否
	结算账户为对私时会返回，商户需使用该付款卡号进行汇款。
	该字段需进行解密处理，解密方法详见敏感信息加解密说明。
	示例值：9nZYDEvBT4rDdICA3ZYXshYqeOSslSjSauAE+as7mAEDJly+DgRuVs74msmKUH8pl+3oA==
	*/
	PayAmount                int    `json:"pay_amount"`                 // 汇款金额	是	需要汇款的金额(单位：分)。	示例值：124
	DestinationAccountNumber string `json:"destination_account_number"` // 收款卡号	[1,128]	是	收款账户的卡号。	示例值：7222223333322332
	DestinationAccountName   string `json:"destination_account_name"`   // 收款户名	[1,128]	是	收款账户名。	示例值：财付通支付科技有限公司
	DestinationAccountBank   string `json:"destination_account_bank"`   // 开户银行	[1,128]	是	收款账户的开户银行名称。	示例值：招商银行威盛大厦支行
	City                     string `json:"city"`                       // 省市信息	[1,128]	是	收款账户的省市。	示例值：深圳
	Remark                   string `json:"remark"`                     // 备注信息	[1,128]	是	商户汇款时，需要填写的备注信息。	示例值：入驻账户验证
	Deadline                 string `json:"deadline"`                   // 汇款截止时间	[1,20]	是	请在此时间前完成汇款。	示例值：2018-12-10 17:09:01
}

type auditDetail struct {
	ParamName    string `json:"param_name"`    // 参数名称[1,32]	是	提交申请单的资料项名称。	示例值：id_card_copy
	RejectReason string `json:"reject_reason"` // 	驳回原因	[1,32]	是	提交资料项被驳回原因。示例值：身份证背面识别失败，请上传更清晰的身份证图片
}

func GetApplyStatusByApplymentID(applymentID uint64) (*getApplyStatusRes, error) {

	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/applyments/%v", applymentID)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}

	var data getApplyStatusRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

func GetApplyStatusByOutRequestNo(outRequestNo string) (*getApplyStatusRes, error) {

	api := fmt.Sprintf("https://api.mch.weixin.qq.com/v3/ecommerce/applyments/out-request-no/%v", outRequestNo)
	res, httpErr := net.HttpGet(api, nil)
	if httpErr != nil {
		return nil, httpErr
	}
	var data getApplyStatusRes
	if jsonErr := json.Unmarshal(res, &data); jsonErr != nil {
		return nil, jsonErr
	}

	return &data, nil
}

// ==================== TODO 获取证书 ====================
func GetCertificates() {
	//api := "https://api.mch.weixin.qq.com/v3/certificates"
}

// ==================== TODO 修改结算账号 ====================

// ==================== TODO 查询结算账号 ====================

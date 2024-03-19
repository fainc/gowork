package wepay

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
	"github.com/gogf/gf/frame/g"
)

type v2Service struct{}

var v2Var = v2Service{}

func V2() *v2Service {
	return &v2Var
}

var v2Client *wechat.Client // 单例对象

var lock sync.Mutex

type InitV2ClientParams struct {
	AppId         string
	MchId         string
	ApiKey        string
	Debug         bool
	CertPath      string
	ReadConfigKey string // 是否读取文件配置，指定前缀
}

func (receiver *v2Service) NewClient(params *InitV2ClientParams) (*wechat.Client, error) {

	var AppId string
	var MchId string
	var ApiKey string
	var Debug bool
	var CertPath string

	if params.ReadConfigKey != "" {
		AppId = g.Cfg().GetString(params.ReadConfigKey + ".AppId")
		if AppId == "" {
			return nil, errors.New("读取指定微信支付配置无效")
		}
		MchId = g.Cfg().GetString(params.ReadConfigKey + ".MchId")
		ApiKey = g.Cfg().GetString(params.ReadConfigKey + ".ApiKey")
		Debug = g.Cfg().GetBool(params.ReadConfigKey + ".Debug")
		CertPath = g.Cfg().GetString(params.ReadConfigKey + ".CertPath")
	} else {
		AppId = params.AppId
		MchId = params.MchId
		ApiKey = params.ApiKey
		Debug = params.Debug
		CertPath = params.CertPath
	}

	client := wechat.NewClient(AppId, MchId, ApiKey, true)
	if Debug {
		client.DebugSwitch = gopay.DebugOn
	}
	if CertPath != "" {
		err := client.AddCertPkcs12FilePath(CertPath) // 使用微信pkcs12证书
		if err != nil {
			return client, errors.New("初始化微信pkcs12证书失败")
		}
	}
	return client, nil
}

// ClientInstance 懒汉单例
func (receiver *v2Service) ClientInstance(params *InitV2ClientParams) (*wechat.Client, error) {
	if v2Client == nil {
		lock.Lock()
		defer lock.Unlock()
		if v2Client == nil {
			client, err := receiver.InitInstance(params)
			if err != nil {
				return nil, err
			}
			v2Client = client
		}
	}
	return v2Client, nil
}

// GetInstance 饿汉单例，使用前需要先调用InitInstance初始化
func (receiver *v2Service) GetInstance() *wechat.Client {
	return v2Client
}

// InitInstance 初始化单例
func (receiver *v2Service) InitInstance(params *InitV2ClientParams) (*wechat.Client, error) {
	client, err := receiver.NewClient(params)
	if err != nil {
		return nil, err
	}
	v2Client = client
	return client, nil
}

type UnifiedOrderV2Params struct {
	AppType    string // * 应用类型 JSAPI/WEAPP/APP，H5和公众号使用JSAPI应用类型
	Body       string // * 商品描述
	OutTradeNo string // * 订单号
	TotalFee   string // * 付款金额单位为【分】
	NotifyUrl  string // * 通知url
	TradeType  string // * JSAPI--JSAPI支付（或小程序支付）、NATIVE--Native支付、APP--app支付，MWEB--H5支付
	AppId      string // *
	ApiKey     string // *
	OpenId     string // trade_type=JSAPI时（即JSAPI支付），此参数必传
	Attach     string // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
}

// UnifiedOrder 微信统一下单
func (receiver *v2Service) UnifiedOrder(client *wechat.Client, params *UnifiedOrderV2Params) (string, error) {
	if client == nil || client.AppId == "" {
		return "", errors.New("微信支付Client未初始化")
	}
	// 初始化 BodyMap
	bm := make(gopay.BodyMap)
	bm.Set("nonce_str", util.RandomString(32)).
		Set("body", params.Body).
		Set("out_trade_no", params.OutTradeNo).
		Set("total_fee", params.TotalFee).
		Set("spbill_create_ip", "127.0.0.1").
		Set("notify_url", params.NotifyUrl).
		Set("trade_type", params.TradeType)
	if params.Attach != "" {
		bm.Set("attach", params.Attach)
	}
	if params.TradeType == "JSAPI" && params.OpenId != "" {
		bm.Set("openid", params.OpenId)
	}
	wxRsp, err := client.UnifiedOrder(context.TODO(), bm)
	if err != nil {
		return "", errors.New("微信统一下单失败")
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 获取小程序支付需要的paySign
	pac := "prepay_id=" + wxRsp.PrepayId
	var paySign = ""
	if params.AppType == "WEAPP" {
		paySign = wechat.GetMiniPaySign(params.AppId, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, params.ApiKey)
	}
	if params.AppType == "JSAPI" {
		paySign = wechat.GetJsapiPaySign(params.AppId, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, params.ApiKey)
	}
	if params.AppType == "APP" {
		paySign = wechat.GetAppPaySign(params.AppId, "", wxRsp.NonceStr, wxRsp.PrepayId, wechat.SignType_MD5, timeStamp, params.ApiKey)
	}
	return paySign, nil
}

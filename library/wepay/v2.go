package wepay

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/pkg/util"
	"github.com/go-pay/gopay/wechat"
)

type wePayV2Service struct{}

var wePayV2Var = wePayV2Service{}

func WePayV2() *wePayV2Service {
	return &wePayV2Var
}

type InitV2ClientParams struct {
	AppId    string
	MchId    string
	ApiKey   string
	Debug    bool
	CertPath string
}

// InitClient 初始化微信配置
func (receiver *wePayV2Service) InitClient(params *InitV2ClientParams) (*wechat.Client, error) {
	client := wechat.NewClient(params.AppId, params.MchId, params.ApiKey, true)

	// 打开调试
	if params.Debug {
		client.DebugSwitch = gopay.DebugOn
	}

	// 使用微信pkcs12证书
	if params.CertPath != "" {
		err := client.AddCertPkcs12FilePath(params.CertPath)
		if err != nil {
			return client, errors.New("初始化微信pkcs12证书失败")
		}
	}
	return client, nil
}

type UnifiedOrderV2Params struct {
	AppType      string // * 应用类型 JSAPI/WEAPP/APP，H5和公众号使用JSAPI应用类型
	Body         string // * 商品描述
	OutTradeNo   string // * 订单号
	TotalFee     string // * 付款金额单位为【分】
	NotifyUrl    string // * 通知url
	TradeType    string // * JSAPI--JSAPI支付（或小程序支付）、NATIVE--Native支付、APP--app支付，MWEB--H5支付
	OpenId       string // trade_type=JSAPI时（即JSAPI支付），此参数必传
	Attach       string // 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据
	ClientParams *InitV2ClientParams
}

// UnifiedOrder 微信统一下单
func (receiver *wePayV2Service) UnifiedOrder(params *UnifiedOrderV2Params) (error, string) {
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
	client, err := receiver.InitClient(params.ClientParams)
	// 请求支付下单，成功后得到结果
	wxRsp, err := client.UnifiedOrder(context.TODO(), bm)
	if err != nil {
		return errors.New("微信统一下单失败"), ""
	}
	timeStamp := strconv.FormatInt(time.Now().Unix(), 10)

	// 获取小程序支付需要的paySign
	pac := "prepay_id=" + wxRsp.PrepayId
	var paySign = ""
	if params.AppType == "WEAPP" {
		paySign = wechat.GetMiniPaySign(params.ClientParams.AppId, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, params.ClientParams.ApiKey)
	}
	if params.AppType == "JSAPI" {
		paySign = wechat.GetJsapiPaySign(params.ClientParams.AppId, wxRsp.NonceStr, pac, wechat.SignType_MD5, timeStamp, params.ClientParams.ApiKey)
	}
	if params.AppType == "APP" {
		paySign = wechat.GetAppPaySign(params.ClientParams.AppId, "", wxRsp.NonceStr, wxRsp.PrepayId, wechat.SignType_MD5, timeStamp, params.ClientParams.ApiKey)
	}
	return nil, paySign
}

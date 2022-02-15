package sms

import (
	"errors"

	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysms "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	"github.com/alibabacloud-go/tea/tea"
)

var Service = smsService{}

type smsService struct{}

func (s *smsService) createClient(accessKeyId *string, accessKeySecret *string) (_result *dysms.Client, _err error) {
	config := &openapi.Config{
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
	}
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result, _err = dysms.NewClient(config)
	return _result, _err
}

type SendSmsParams struct {
	PhoneNumbers    string `dc:"手机号码，多号码用,分隔"`
	SignName        string `dc:"签名"`
	TemplateCode    string `dc:"短信模板编码"`
	AccessKeyId     string `dc:"阿里云 AccessKeyId"`
	AccessKeySecret string `dc:"阿里云 AccessKeySecret"`
	TemplateParam   string `dc:"(可选)短信模板变量参数，接受JSON字符串"`
}

// SendSms 接口是短信发送接口，支持在一次请求中向多个不同的手机号码发送同样内容的短信。
func (s *smsService) SendSms(params *SendSmsParams) (err error) {
	client, dyErr := s.createClient(tea.String(params.AccessKeyId), tea.String(params.AccessKeySecret))
	if dyErr != nil {
		return errors.New("短信对接错误，阿里云短信账户初始化异常")
	}
	sendSmsRequest := &dysms.SendSmsRequest{
		PhoneNumbers:  tea.String(params.PhoneNumbers),
		SignName:      tea.String(params.SignName),
		TemplateCode:  tea.String(params.TemplateCode),
		TemplateParam: tea.String(params.TemplateParam),
	}
	result, dyErr := client.SendSms(sendSmsRequest)
	if dyErr != nil {
		return errors.New("短信对接错误，请检查阿里云短信密钥是否可用、参数是否完整")
	}
	if tea.StringValue(result.Body.Code) != "OK" {
		return errors.New("短信发送失败，" + tea.StringValue(result.Body.Message))
	}
	return nil
}

// SendBatchSms 接口是短信批量发送接口，支持在一次请求中分别向多个不同的手机号码发送不同签名的短信。
// 手机号码等参数均为JSON格式，字段个数相同，一一对应，短信服务根据字段在JSON中的顺序判断发往指定手机号码的签名。
// 如果您需要往多个手机号码中发送同样签名的短信，请使用SendSms接口实现。
func (s *smsService) SendBatchSms(params *SendSmsParams) (err error) {
	client, dyErr := s.createClient(tea.String(params.AccessKeyId), tea.String(params.AccessKeySecret))
	if dyErr != nil {
		return errors.New("短信对接错误，阿里云短信账户初始化异常")
	}
	sendSmsRequest := &dysms.SendBatchSmsRequest{
		PhoneNumberJson:   tea.String(params.PhoneNumbers),
		SignNameJson:      tea.String(params.SignName),
		TemplateCode:      tea.String(params.TemplateCode),
		TemplateParamJson: tea.String(params.TemplateParam),
	}
	result, dyErr := client.SendBatchSms(sendSmsRequest)
	if dyErr != nil {
		return errors.New("短信对接错误，请检查阿里云短信密钥是否可用、参数是否完整")
	}
	if tea.StringValue(result.Body.Code) != "OK" {
		return errors.New("短信发送失败，" + tea.StringValue(result.Body.Message))
	}
	return nil
}

package email

import (
	"crypto/tls"
	"errors"

	"gopkg.in/gomail.v2"
)

type ConfigParams struct {
	Host               string
	Port               int
	User               string
	Password           string
	InsecureSkipVerify bool
}

type SendEMailParams struct {
	From         string        // 发送人完整邮箱账户
	FromName     string        // 发送人账户别名（收件方显示的发件人名称，如不设置默认显示邮箱号开头英文）
	To           []string      // 收件人邮箱号（支持多邮箱）
	Subject      string        // 邮件主题
	Body         string        // 邮件正文（html格式）
	AttachPath   string        // （可选）附件本地文件地址
	CcAddress    []string      // （可选）抄送 完整邮箱账户（支持多邮箱）
	ConfigParams *ConfigParams //  基本配置
}

func SendEMail(params *SendEMailParams) (err error) {
	if len(params.To) == 0 {
		return errors.New("邮件收件人列表不得为空")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(params.From, params.FromName))
	m.SetHeader("To", params.To...)
	if len(params.To) != 0 {
		m.SetHeader("Cc", params.CcAddress...)
	}
	m.SetHeader("Subject", params.Subject)
	m.SetBody("text/html", params.Body)
	if params.AttachPath != "" {
		m.Attach(params.AttachPath)
	}
	d := gomail.NewDialer(params.ConfigParams.Host, params.ConfigParams.Port, params.ConfigParams.User, params.ConfigParams.Password)
	if params.ConfigParams.InsecureSkipVerify {
		d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	}
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

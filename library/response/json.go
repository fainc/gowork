package response

import (
	"github.com/gogf/gf/net/ghttp"
)

var Json = jsonResponse{}

type jsonResponse struct{}

type Res struct {
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	RequestId interface{} `json:"requestId"`
}

func (c *jsonResponse) Success(r *ghttp.Request, data interface{}) {
	c.output(r, 200, 200, "success", data)
}

func (c *jsonResponse) Error(r *ghttp.Request, status, code int, message string) {
	c.output(r, status, code, message, nil)
}

func (c *jsonResponse) StandError(r *ghttp.Request, message string) {
	c.Error(r, 400, 400, message)
}

func (c *jsonResponse) ServerError(r *ghttp.Request, message string) {
	c.Error(r, 500, 500, message)
}

func (c *jsonResponse) NotFound(r *ghttp.Request) {
	c.Error(r, 404, 404, "not found")
}

func (c *jsonResponse) Authorization(r *ghttp.Request, message string) {
	c.Error(r, 401, 401, message)
}

//output 设置状态码、清空buffer后输出json并退出当前业务流程(基于gf框架，非gf框架请勿使用)
//notice:based on gf
func (c *jsonResponse) output(r *ghttp.Request, status int, code int, message string, data interface{}) {
	r.Response.WriteStatus(status)
	r.Response.ClearBuffer()
	_ = r.Response.WriteJsonExit(Res{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestId: r.GetParam("requestId"),
	})
}

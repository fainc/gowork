//notice:based on gf
package response

import (
	"strings"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/guid"
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

func (c *jsonResponse) StandardError(r *ghttp.Request, message string) {
	c.Error(r, 400, 400, message)
}

func (c *jsonResponse) ServerError(r *ghttp.Request) {
	c.Error(r, 500, 500, "服务器发生错误")
}

func (c *jsonResponse) NotFound(r *ghttp.Request) {
	c.Error(r, 404, 404, "request uri not found")
}

func (c *jsonResponse) Authorization(r *ghttp.Request, message string) {
	c.Error(r, 401, 401, message)
}

func (c *jsonResponse) output(r *ghttp.Request, status int, code int, message string, data interface{}) {
	r.Response.WriteStatus(status)
	r.Response.ClearBuffer()
	_ = r.Response.WriteJson(Res{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestId: r.GetCtxVar("requestId", strings.ToUpper(guid.S())),
	})
	r.Response.Header().Set("Content-Type", "application/json;charset=utf-8")
	r.Exit()
}

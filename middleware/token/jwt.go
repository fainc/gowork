package token

import (
	"github.com/gogf/gf/net/ghttp"
	"gowork/library/jwt"
	"gowork/library/response"
)

var JwtToken = jwtToken{}

type jwtToken struct{}

func (*jwtToken) StandardAuth(r *ghttp.Request) {
	uuid, err := jwt.Helper.Parse(r, "user")
	if err != nil {
		response.Json.Authorization(r, err.Error())
	}
	r.SetCtxVar("UUID", uuid)
	r.Middleware.Next()
}

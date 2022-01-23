package token

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"gowork/library/jwt"
)

func JwtAuth(r *ghttp.Request) {
	scopes := g.SliceStr{
		"user",
	}
	whiteTables := g.SliceStr{
		"/hello/test",
	}
	jwt.Helper.StandardAuth(r, scopes, whiteTables)
	r.Middleware.Next()
}

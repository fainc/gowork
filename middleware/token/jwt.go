package token

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"gowork/library/jwt"
	"gowork/library/response"
)

var Jwt = jwtToken{}

type jwtToken struct{}

// standardAuth 强验证 验证失败后非白名单URI强制返回401并中断请求,白名单内的URI获取UUID可能为0或有效UUID（类似弱验证）
func (*jwtToken) standardAuth(r *ghttp.Request, scope string) {
	whiteTable := garray.NewStrArrayFrom(g.SliceStr{
		"/api/account/login",
		"/api/account/valid",
	})
	uuid, scopeKey, err := jwt.Helper.Parse(r, scope)
	if err != nil {
		if !whiteTable.ContainsI(r.RequestURI) {
			response.Json.Authorization(r, err.Error())
		} else {
			r.SetCtxVar("UUID", 0)
			r.SetCtxVar("SCOPE", "UNKNOWN")
		}
	} else {
		r.SetCtxVar("UUID", uuid)
		r.SetCtxVar("SCOPE", scopeKey)
	}
}

// weakAuth 弱验证，验证失败UUID赋值为0代表未登录，不强制返回401中断请求，弱验证忽略白名单
func (*jwtToken) weakAuth(r *ghttp.Request, scope string) {
	uuid, scopeKey, err := jwt.Helper.Parse(r, scope)
	if err != nil {
		r.SetCtxVar("UUID", 0)
		r.SetCtxVar("SCOPE", "UNKNOWN")
	} else {
		r.SetCtxVar("UUID", uuid)
		r.SetCtxVar("SCOPE", scopeKey)
	}
}

type jwtParse struct {
	UUID  int64
	SCOPE string
}

func (*jwtToken) Parse(r *ghttp.Request) *jwtParse {
	UUID := r.GetCtxVar("UUID", 0)
	SCOPE := r.GetCtxVar("SCOPE", "UNKNOWN")
	return &jwtParse{
		UUID:  gconv.Int64(UUID),
		SCOPE: gconv.String(SCOPE),
	}
}

//在此实现自己的验证逻辑

//AdminAuth 管理员jwt验证，scope为admin，验证失败强制返回401
func (*jwtToken) AdminAuth(r *ghttp.Request) {
	Jwt.standardAuth(r, "admin")
	r.Middleware.Next()
}

//UserAuth 用户jwt验证，scope为user，验证失败强制返回401
func (*jwtToken) UserAuth(r *ghttp.Request) {
	Jwt.standardAuth(r, "user")
	r.Middleware.Next()
}

//UserOptionalAuth 用户可选验证，UUID可为0，不强制返回401
func (*jwtToken) UserOptionalAuth(r *ghttp.Request) {
	Jwt.weakAuth(r, "user")
	r.Middleware.Next()
}

//AdminOptionalAuth 管理员可选验证，UUID可为0，不强制返回401
func (*jwtToken) AdminOptionalAuth(r *ghttp.Request) {
	Jwt.weakAuth(r, "admin")
	r.Middleware.Next()
}

package cors

import "github.com/gogf/gf/net/ghttp"

func Standard(r *ghttp.Request) {
	corsOptions := r.Response.DefaultCORSOptions()
	corsOptions.AllowDomain = []string{"*"}
	r.Response.CORS(corsOptions)
	r.Middleware.Next()
}

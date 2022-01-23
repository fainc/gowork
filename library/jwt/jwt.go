package jwt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/golang-jwt/jwt/v4"
	"strings"
	"time"
)

var Helper = jwtHelper{}

type jwtHelper struct{}

func (*jwtHelper) Parse(r *ghttp.Request, scopesSlice g.SliceStr) (string, string, error) {
	secret := g.Cfg().GetString("jwt.secret") //配置修改会自动刷新
	if secret == "" {
		return "", "", errors.New("jwt secret invalid")
	}
	tokenString := r.GetHeader("Authorization")
	if tokenString == "" {
		return "", "", errors.New("authorization invalid")
	}
	tokenMap := strings.Split(tokenString, "Bearer ")
	if len(tokenMap) != 2 {
		return "", "", errors.New("bearer invalid")
	}
	tokenString = tokenMap[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return "", "", errors.New(err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", "", errors.New(err.Error())
	}
	uuid := claims["uuid"]
	if uuid == nil {
		return "", "", errors.New("signature user key invalid")
	}
	encodeString, err := hex.DecodeString(gconv.String(uuid))
	if err != nil {
		return "", "", errors.New("signature user key decode hex fail")
	}
	uuid, err = gaes.Decrypt(encodeString, []byte(secret))
	if err != nil {
		return "", "", errors.New("signature user key decrypt fail")
	}
	scope := claims["scope"]
	scopes := garray.NewStrArrayFrom(scopesSlice)
	if scope == nil || !scopes.ContainsI(gconv.String(scope)) {
		return "", "", errors.New("scope invalid")
	}
	return gconv.String(uuid), gconv.String(scope), nil
}

func (*jwtHelper) Generate(uuid string, scope string, duration time.Duration) (string, error) {
	secret := g.Cfg().GetString("jwt.secret")
	if secret == "" {
		return "", errors.New("jwt secret invalid")
	}
	type MyCustomClaims struct {
		Uuid  string `json:"uuid"`
		Scope string `json:"scope"`
		jwt.RegisteredClaims
	}
	uuidEncode, _ := gaes.Encrypt([]byte(uuid), []byte(secret))
	claims := MyCustomClaims{
		hex.EncodeToString(uuidEncode),
		scope,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jwtHelper",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString, nil
}

type user struct {
	UUID  int64
	SCOPE string
}

func (*jwtHelper) GetUser(r *ghttp.Request) *user {
	UUID := r.GetCtxVar("UUID", 0)
	SCOPE := r.GetCtxVar("SCOPE", "UNKNOWN")
	return &user{
		UUID:  gconv.Int64(UUID),
		SCOPE: gconv.String(SCOPE),
	}
}

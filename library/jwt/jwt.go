package jwt

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/crypto/gaes"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/golang-jwt/jwt/v4"
)

var Helper = jwtHelper{}

type jwtHelper struct{}

type ParseParams struct {
	Token  string     // * jwt字符串
	Scopes g.SliceStr // * jwt scope可用范围
	Secret string     // * jwt密钥
}

// Parse jwt解析、解密
func (*jwtHelper) Parse(params ParseParams) (int64, string, error) {
	// secret := g.Cfg().GetString("jwt.secret") //配置修改会自动刷新
	if params.Secret == "" {
		return 0, "", errors.New("jwt secret invalid")
	}
	// tokenString := r.GetHeader("Authorization")
	if params.Token == "" {
		return 0, "", errors.New("authorization invalid")
	}
	tokenMap := strings.Split(params.Token, "Bearer ")
	if len(tokenMap) != 2 {
		return 0, "", errors.New("bearer invalid")
	}
	tokenString := tokenMap[1]
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(params.Secret), nil
	})
	if err != nil {
		return 0, "", errors.New(err.Error())
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New(err.Error())
	}
	uuid := claims["uuid"]
	if uuid == nil {
		return 0, "", errors.New("signature user key invalid")
	}
	encodeString, err := hex.DecodeString(gconv.String(uuid))
	if err != nil {
		return 0, "", errors.New("signature user key decode hex fail")
	}
	uuid, err = gaes.Decrypt(encodeString, []byte(params.Secret))
	if err != nil {
		return 0, "", errors.New("signature user key decrypt fail")
	}
	scope := claims["scope"]
	scopes := garray.NewStrArrayFrom(params.Scopes)
	if scope == nil || !scopes.ContainsI(gconv.String(scope)) {
		return 0, "", errors.New("scope invalid")
	}
	return gconv.Int64(uuid), gconv.String(scope), nil
}

type GenerateParams struct {
	Uuid     int64         // * 非0用户ID
	Scope    string        // * 授权scope标志
	Duration time.Duration // * 授权时长
	Secret   string        // * jwt及加密密钥
}

// Generate 生成jwt
func (*jwtHelper) Generate(params GenerateParams) (string, error) {
	// secret := g.Cfg().GetString("jwt.secret")
	if params.Uuid == 0 || params.Scope == "" || params.Duration == 0 || params.Secret == "" {
		return "", errors.New("generate jwt params invalid")
	}

	type MyCustomClaims struct {
		Uuid  string `json:"uuid"`
		Scope string `json:"scope"`
		jwt.RegisteredClaims
	}
	uuidEncode, _ := gaes.Encrypt([]byte(gconv.String(params.Uuid)), []byte(params.Secret))
	claims := MyCustomClaims{
		hex.EncodeToString(uuidEncode),
		params.Scope,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(params.Duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "jwtHelper",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(params.Secret))
	return tokenString, nil
}

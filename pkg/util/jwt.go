package util

import (
	"Gin-blog-example/pkg/setting"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtSecret = []byte(setting.AppSetting.JwtSecret)

type Claims struct {
	Username string `json:"username"`
	Password string `json:"password"`
	jwt.StandardClaims
}

//

// 生成 token
func GenerateToken(username, password string) (string, error) {
	nowTime := time.Now()
	//生成过期时间
	expireTime := nowTime.Add(3 * time.Hour)
	claims := Claims{
		username,
		password,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间
			Issuer:    "gin-blog",        //签发人
		},
	}
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString(jwtSecret)
}

// 解析 token
func ParseToken(token string) (*Claims, error) {
	//对token进行解析
	tokenClaims, e := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (i interface{}, e error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		//如果解析成功，并且验证通过，说明jwt符合，进行还原数据返回
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, e
}

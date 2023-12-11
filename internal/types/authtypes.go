package types

import "github.com/golang-jwt/jwt/v4"

type TokenClaims struct {
	UserId   int64
	Phonenum string
	UserType string
	Role     []string
	Token    string
	Ext      map[string]any
	jwt.RegisteredClaims
}

type AuthRequest struct {
	Phonenum  string `valid:"Required" form:"phonenum" json:"phonenum"`
	Password  string `valid:"Required" form:"password" json:"password"`
	CaptchaId string `valid:"Required" form:"captcha_id" json:"captcha_id"`
	Captcha   string `valid:"Required" json:"captcha" form:"captcha"`
}

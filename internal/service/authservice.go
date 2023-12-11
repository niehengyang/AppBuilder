package service

import (
	"appBuilder/internal/service/resp"
	"appBuilder/internal/svc"
	"appBuilder/internal/types"
	"appBuilder/model"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logc"
	"net/http"
	"time"
)

// GetJwtToken
// @secretKey: JWT 加解密密钥
// @seconds: 过期时间，单位秒
// @token: 数据载体
func GetJwtToken(svcCtx *svc.ServiceContext, token *types.TokenClaims) (string, error) {
	token.RegisteredClaims = jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * time.Duration(svcCtx.Config.JwtAuth.AccessExpire))),
		Issuer:    "nhy",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, token)
	mySigningKey := []byte(svcCtx.Config.JwtAuth.AccessSecret)
	signedString, err := claims.SignedString(mySigningKey)
	return signedString, err
}

func Secret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}
}

func ParseToken(tokenss string, svcCtx *svc.ServiceContext) (*types.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenss, &types.TokenClaims{}, Secret(svcCtx.Config.JwtAuth.AccessSecret))
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, errors.New("token不合法")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, errors.New("token已失效")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				return nil, errors.New("token未激活")
			} else {
				return nil, errors.New("无法处理此token")
			}
		}
	}
	if claims, ok := token.Claims.(*types.TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("无法处理此token")
}

// GetLoginUser 获取当前登录用户---普通用户
func GetLoginUser(r *http.Request, svcCtx *svc.ServiceContext) (*model.User, error) {
	user := model.User{}
	tokenString := r.Header.Get("Authorization")
	claims, err := ParseToken(tokenString, svcCtx)
	if err != nil {
		logc.Error(context.Background(), "[GetLoginUser] error,", err.Error())
		return nil, resp.NewError401("")
	}
	queryRes := svcCtx.DB.First(&user, "id = ?", claims.UserId)
	if queryRes.Error != nil {
		return &user, resp.NewError404()
	}

	return &user, nil
}

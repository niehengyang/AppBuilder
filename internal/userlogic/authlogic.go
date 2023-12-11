package userlogic

import (
	"appBuilder/internal/service"
	"appBuilder/internal/service/captcha"
	"appBuilder/internal/service/resp"
	"appBuilder/internal/service/rsa"
	"appBuilder/internal/service/utils"
	"appBuilder/internal/types"
	"appBuilder/model"
	"context"
	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logc"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"

	"appBuilder/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type AuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	req    *http.Request
}

func NewAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext, req *http.Request) *AuthLogic {
	return &AuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		req:    req,
	}
}

func (l *AuthLogic) getJwtToken(secretKey string, iat, seconds, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secretKey))
}

func (l *AuthLogic) Login(req *types.AuthRequest) (*types.Response, error) {
	// 校验图形验证码 Test关闭
	if ok := captcha.CaptchaVerify(req.CaptchaId, req.Captcha); !ok {
		return nil, resp.NewError(http.StatusUnauthorized, "验证码错误")
	}

	// 账号解密
	phonenum, errLo := rsa.RsaDecrypt(req.Phonenum)
	if errLo != nil {
		logc.Error(context.Background(), "[login]账号解密失败", errLo)
		return nil, resp.NewError(resp.Err_LoginFailed, "登录失败，请联系管理员")
	}

	// 密码解密
	pwdText, errCr := rsa.RsaDecrypt(req.Password)
	if errCr != nil {
		logc.Error(context.Background(), "[login]密码解密失败", errLo)
		return nil, resp.NewError(resp.Err_LoginFailed, "登录失败，请联系管理员")
	}

	user := new(model.User)
	result := l.svcCtx.DB.Where("phonenum = ?", phonenum).First(&user)
	if result.Error != nil {
		return nil, resp.NewError401("验证码错误")
	}

	if user.Status != model.AccountStatus_Avaliable {
		return nil, resp.NewError401("账号已禁用")
	}

	// 密码对比
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwdText))
	if err != nil {
		return nil, resp.NewError401("账号或密码错误")
	}

	claims := types.TokenClaims{
		UserId:   int64(user.ID),
		Phonenum: phonenum,
		UserType: "web",
		Role:     []string{"user"},
		Ext:      map[string]any{"ex1": "val1", "ext2": "val2"},
	}

	jwtToken, err1 := service.GetJwtToken(l.svcCtx, &claims)
	if err1 != nil {
		logc.Error(context.Background(), "[login]token生成失败", errLo)
		return nil, resp.NewError(resp.Err_LoginFailed, "登录失败，请联系管理员")
	}

	user.LastLogin = utils.TimestampFormat(time.Now().Unix(), utils.TimeStrTemplate1)
	user.Token = jwtToken
	l.svcCtx.DB.Save(&user)

	return &types.Response{
		Message: "success",
		Data: map[string]string{
			"token": jwtToken,
		},
	}, nil
}

func (l *AuthLogic) CurrentAccount() (*types.Response, error) {
	logc.Debug(context.Background(), "获取当前用户")

	account, _ := service.GetLoginUser(l.req, l.svcCtx)

	return &types.Response{
		Message: "Success",
		Data:    account,
	}, nil
}

// AuthLogout
// 退出登录
func (l *AuthLogic) AuthLogout() (*resp.Response, error) {
	loginAccount, _ := service.GetLoginUser(l.req, l.svcCtx)

	loginAccount.Token = ""
	res := l.svcCtx.DB.Save(&loginAccount)
	if res.Error != nil {
		return nil, resp.NewError500()
	}

	return resp.Success(nil), nil
}

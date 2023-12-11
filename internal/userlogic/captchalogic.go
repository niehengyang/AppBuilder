package userlogic

import (
	"appBuilder/internal/service/captcha"
	resp2 "appBuilder/internal/service/resp"
	"appBuilder/internal/svc"
	"appBuilder/internal/types"
	"context"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type CaptchaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	req    *http.Request
}

func NewCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext, req *http.Request) *CaptchaLogic {
	return &CaptchaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
		req:    req,
	}
}

func (l *CaptchaLogic) GetCaptcha() (resp *types.Response, err error) {
	id, image, _, err := captcha.CaptchaGenarate()
	if err != nil {
		logc.Error(context.Background(), "获取图形验证码失败,", err)
		return nil, resp2.NewError500()
	}

	return &types.Response{
		Message: "success",
		Data: map[string]string{
			"captcha_id":    id,
			"captcha_image": image,
		},
	}, nil
}

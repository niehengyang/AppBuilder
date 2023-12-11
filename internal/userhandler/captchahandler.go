package userhandler

import (
	"appBuilder/internal/svc"
	"appBuilder/internal/userlogic"
	"github.com/zeromicro/go-zero/rest/httpx"
	"net/http"
)

func GetCaptchaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := userlogic.NewCaptchaLogic(r.Context(), svcCtx, r)
		resp, err := l.GetCaptcha()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}

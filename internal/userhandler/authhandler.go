package userhandler

import (
	"appBuilder/internal/service/resp"
	"net/http"

	"appBuilder/internal/svc"
	"appBuilder/internal/types"
	"appBuilder/internal/userlogic"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func AuthLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.AuthRequest

		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.Error(w, resp.NewError400(err.Error()))
			return
		}
		l := userlogic.NewAuthLogic(r.Context(), svcCtx, r)
		resP, err := l.Login(&req)
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resP)
		}
	}
}

func CurrentAccountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := userlogic.NewAuthLogic(r.Context(), svcCtx, r)
		resP, err := l.CurrentAccount()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resP)
		}
	}
}

// AuthLogoutHandler
// 退出登录
func AuthLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := userlogic.NewAuthLogic(r.Context(), svcCtx, r)
		res, err := l.AuthLogout()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, res)
		}
	}
}

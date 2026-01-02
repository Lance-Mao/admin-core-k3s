package token

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/logic/token"
	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/svc"
	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/types"
)

// swagger:route post /token/update token UpdateToken
//
// Update token information | 更新 Token
//
// Update token information | 更新 Token
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: TokenInfo
//
// Responses:
//  200: BaseMsgResp

func UpdateTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TokenInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := token.NewUpdateTokenLogic(r.Context(), svcCtx)
		resp, err := l.UpdateToken(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

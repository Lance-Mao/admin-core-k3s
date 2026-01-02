package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/logic/role"
	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/svc"
	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/types"
)

// swagger:route post /role/create role CreateRole
//
// Create role information | 创建角色
//
// Create role information | 创建角色
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: RoleInfo
//
// Responses:
//  200: BaseMsgResp

func CreateRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewCreateRoleLogic(r.Context(), svcCtx)
		resp, err := l.CreateRole(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

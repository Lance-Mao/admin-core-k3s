package smslog

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/logic/smslog"
	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/svc"
	"github.com/Lance-Mao/admin-core-k3s/app/api/internal/types"
)

// swagger:route post /sms_log/update smslog UpdateSmsLog
//
// Update sms log information | 更新短信日志
//
// Update sms log information | 更新短信日志
//
// Parameters:
//  + name: body
//    require: true
//    in: body
//    type: SmsLogInfo
//
// Responses:
//  200: BaseMsgResp

func UpdateSmsLogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SmsLogInfo
		if err := httpx.Parse(r, &req, true); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := smslog.NewUpdateSmsLogLogic(r.Context(), svcCtx)
		resp, err := l.UpdateSmsLog(&req)
		if err != nil {
			err = svcCtx.Trans.TransError(r.Context(), err)
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

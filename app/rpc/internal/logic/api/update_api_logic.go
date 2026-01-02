package api

import (
	"context"

	"github.com/Lance-Mao/admin-core-k3s/app/rpc/internal/svc"
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/internal/utils/dberrorhandler"
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/suyuan32/simple-admin-common/i18n"
)

type UpdateApiLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateApiLogic {
	return &UpdateApiLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateApiLogic) UpdateApi(in *core.ApiInfo) (*core.BaseResp, error) {
	err := l.svcCtx.DB.API.UpdateOneID(*in.Id).
		SetNotNilPath(in.Path).
		SetNotNilDescription(in.Description).
		SetNotNilAPIGroup(in.ApiGroup).
		SetNotNilMethod(in.Method).
		SetNotNilIsRequired(in.IsRequired).
		SetNotNilServiceName(in.ServiceName).
		Exec(l.ctx)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	return &core.BaseResp{Msg: i18n.UpdateSuccess}, nil
}

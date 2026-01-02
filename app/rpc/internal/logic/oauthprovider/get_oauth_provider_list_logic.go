package oauthprovider

import (
	"context"

	"github.com/suyuan32/simple-admin-common/utils/pointy"

	"github.com/Lance-Mao/admin-core-k3s/app/rpc/ent/oauthprovider"
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/ent/predicate"

	"github.com/Lance-Mao/admin-core-k3s/app/rpc/internal/svc"
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/internal/utils/dberrorhandler"
	"github.com/Lance-Mao/admin-core-k3s/app/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOauthProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOauthProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOauthProviderListLogic {
	return &GetOauthProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOauthProviderListLogic) GetOauthProviderList(in *core.OauthProviderListReq) (*core.OauthProviderListResp, error) {
	var predicates []predicate.OauthProvider
	if in.Name != nil {
		predicates = append(predicates, oauthprovider.NameContains(*in.Name))
	}
	result, err := l.svcCtx.DB.OauthProvider.Query().Where(predicates...).Page(l.ctx, in.Page, in.PageSize)
	if err != nil {
		return nil, dberrorhandler.DefaultEntError(l.Logger, err, in)
	}

	resp := &core.OauthProviderListResp{}
	resp.Total = result.PageDetails.Total

	for _, v := range result.List {
		resp.Data = append(resp.Data, &core.OauthProviderInfo{
			Id:           &v.ID,
			CreatedAt:    pointy.GetPointer(v.CreatedAt.UnixMilli()),
			UpdatedAt:    pointy.GetPointer(v.UpdatedAt.UnixMilli()),
			Name:         &v.Name,
			ClientId:     &v.ClientID,
			ClientSecret: &v.ClientSecret,
			RedirectUrl:  &v.RedirectURL,
			Scopes:       &v.Scopes,
			AuthUrl:      &v.AuthURL,
			TokenUrl:     &v.TokenURL,
			AuthStyle:    &v.AuthStyle,
			InfoUrl:      &v.InfoURL,
		})
	}

	return resp, nil
}

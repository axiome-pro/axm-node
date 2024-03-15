package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axiome-pro/axm-node/x/vote/types"
)

type QueryServer struct {
	Keeper
}

var _ types.QueryServer = QueryServer{}

func (qs QueryServer) History(ctx context.Context, req *types.HistoryRequest) (*types.HistoryResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = qs.Keeper
	)
	limit := req.Limit
	page := req.Page

	if limit <= 0 || limit > 100 {
		limit = 100
	}

	if page < 1 {
		page = 1
	}

	data := k.GetHistory(sdkCtx, limit, page)
	return &types.HistoryResponse{
		History: data,
	}, nil
}

func (qs QueryServer) Government(ctx context.Context, _ *types.GovernmentRequest) (*types.GovernmentResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = qs.Keeper
	)
	data := k.GetGovernment(sdkCtx)
	return &types.GovernmentResponse{
		Members: data.Members,
	}, nil
}

func (qs QueryServer) Current(ctx context.Context, _ *types.CurrentRequest) (*types.CurrentResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = qs.Keeper
	)
	var (
		proposal = k.GetCurrentProposal(sdkCtx)
		gov,
		agreed,
		disagreed types.Government
	)
	if proposal == nil {
		proposal = new(types.Proposal)
	} else {
		gov = k.GetGovernment(sdkCtx)
		agreed = k.GetAgreed(sdkCtx)
		disagreed = k.GetDisagreed(sdkCtx)
	}
	return &types.CurrentResponse{
		Proposal:   *proposal,
		Government: gov.Strings(),
		Agreed:     agreed.Strings(),
		Disagreed:  disagreed.Strings(),
	}, nil
}

func (qs QueryServer) Params(ctx context.Context, _ *types.ParamsRequest) (*types.ParamsResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = qs.Keeper
	)
	data := k.GetParams(sdkCtx)
	return &types.ParamsResponse{
		Params: data,
	}, nil
}

func (qs QueryServer) Poll(ctx context.Context, _ *types.PollRequest) (*types.PollResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = qs.Keeper
	)
	poll, ok := k.GetCurrentPoll(sdkCtx)
	if !ok {
		return nil, status.Error(codes.NotFound, "There is no active poll at the moment")
	}
	yes, no := k.GetPollStatus(sdkCtx)

	return &types.PollResponse{
		Poll: poll,
		Yes:  yes,
		No:   no,
	}, nil
}

func (qs QueryServer) PollHistory(ctx context.Context, req *types.PollHistoryRequest) (*types.PollHistoryResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = qs.Keeper
	)
	data := k.GetPollHistory(sdkCtx, req.Limit, req.Page)
	return &types.PollHistoryResponse{History: data}, nil
}

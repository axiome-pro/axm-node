package keeper

import (
	"context"

	"github.com/axiome-pro/axm-node/x/referral/types"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ types.QueryServer = Querier{}

type Querier struct {
	Keeper
}

func (qs Querier) Exists(ctx context.Context, request *types.ExistsRequest) (*types.ExistsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	return &types.ExistsResponse{
		Exists: qs.exists(sdkCtx, request.AccAddress),
	}, nil
}

func (qs Querier) Children(ctx context.Context, request *types.ChildrenRequest) (*types.ChildrenResponse, error) {
	// TODO: validate acc address for correct bech32
	if request.AccAddress == "" {
		return &types.ChildrenResponse{}, nil
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	children, err := qs.Keeper.GetChildren(sdkCtx, request.AccAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot obtain account children data")
	}

	return &types.ChildrenResponse{
		Children: children,
	}, nil
}

func NewQuerier(keeper Keeper) Querier {
	return Querier{Keeper: keeper}
}

func (qs Querier) Get(ctx context.Context, request *types.GetRequest) (*types.GetResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	info, err := qs.Keeper.Get(sdkCtx, request.AccAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot obtain account data")
	}

	return &types.GetResponse{Info: info}, nil
}

func (qs Querier) Coins(ctx context.Context, request *types.CoinsRequest) (*types.CoinsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	d, err := qs.Keeper.GetDelegatedInNetwork(sdkCtx, request.AccAddress)
	if err != nil {
		return nil, errors.Wrap(err, "cannot obtain delegated coins")
	}
	return &types.CoinsResponse{
		Delegated: d,
	}, nil
}

func (qs Querier) CheckStatus(ctx context.Context, request *types.CheckStatusRequest) (*types.CheckStatusResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	result, err := qs.Keeper.AreStatusRequirementsFulfilled(sdkCtx, request.AccAddress, request.Status)
	if err != nil {
		return nil, errors.Wrap(err, "cannot obtain data")
	}
	return &types.CheckStatusResponse{Result: result}, nil
}

func (qs Querier) Params(ctx context.Context, _ *types.ParamsRequest) (*types.ParamsResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	params := qs.Keeper.GetParams(sdkCtx)

	return &types.ParamsResponse{Params: params}, nil
}

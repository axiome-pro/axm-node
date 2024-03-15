package keeper

import (
	"github.com/axiome-pro/axm-node/x/vote/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// GetParams returns the total set of referral parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	params, _ = k.Params.Get(ctx)
	return params
}

// SetParams sets the referral parameters to the param space.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.Logger(ctx).Debug("SetParams", "params", params)
	_ = k.Params.Set(ctx, params)
}

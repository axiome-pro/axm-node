package keeper

import (
	"github.com/axiome-pro/axm-node/x/referral/types"
	"context"
	sdkerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	stakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
)

// Wrapper struct
type Hooks struct {
	k Keeper
}

var _ stakingtypes.RefStakingHooks = Hooks{}

// Create new distribution hooks
func (k Keeper) Hooks() Hooks {
	return Hooks{k}
}

// Staked coins amount modified
func (h Hooks) DelegationCoinsModified(ctx context.Context, delAddr, valAddr string, oldCoins, newCoins math.Int) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := h.k.Logger(sdkCtx)

	logger.Debug("Shared modified hook", "del", delAddr, "val", valAddr, "old", oldCoins, "new", newCoins)

	return h.k.OnBalanceChanged(sdkCtx, delAddr, newCoins.Sub(oldCoins))
}

func (h Hooks) CheckDelegationAvailable(ctx context.Context, delAddr, valAddr string) error {
	if !h.k.exists(sdk.UnwrapSDKContext(ctx), delAddr) {
		return sdkerrors.Wrap(types.ErrNotFound, delAddr)
	}

	return nil
}

func (h Hooks) SpendCoinsForRef(ctx context.Context, addr string, totalAmount math.Int) (math.Int, error) {
	return h.k.PayUpFees(sdk.UnwrapSDKContext(ctx), addr, totalAmount)
}

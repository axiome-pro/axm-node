package types

import (
	"context"
	"cosmossdk.io/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"cosmossdk.io/math"
)

// combine multiple staking hooks, all hook functions are run in array sequence
var _ RefStakingHooks = &MultiRefStakingHooks{}

type MultiRefStakingHooks []RefStakingHooks

func NewMultiRefStakingHooks(hooks ...RefStakingHooks) MultiRefStakingHooks {
	return hooks
}

func (h MultiRefStakingHooks) DelegationCoinsModified(ctx context.Context, delAddr, valAddr string, oldCoins, newCoins math.Int) error {
	for i := range h {
		if err := h[i].DelegationCoinsModified(ctx, delAddr, valAddr, oldCoins, newCoins); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiRefStakingHooks) CheckDelegationAvailable(ctx context.Context, delAddr, valAddr string) error {
	for i := range h {
		if err := h[i].CheckDelegationAvailable(ctx, delAddr, valAddr); err != nil {
			return err
		}
	}

	return nil
}

func (h MultiRefStakingHooks) SpendCoinsForRef(ctx context.Context, addr string, totalAmount math.Int) (math.Int, error) {
	for i := range h {
		remain, err := h[i].SpendCoinsForRef(ctx, addr, totalAmount)

		if remain.LT(math.ZeroInt()) {
			return math.ZeroInt(), errors.Wrap(sdkerrors.ErrInsufficientFunds, "unable to execute all spend coins hooks")
		}

		totalAmount = remain

		if err != nil {
			return math.ZeroInt(), err
		}
	}

	return totalAmount, nil
}

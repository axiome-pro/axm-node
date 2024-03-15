package keeper

import (
	"context"

	"cosmossdk.io/math"
)

// GetCommunityTax returns the current distribution community tax.
func (k Keeper) GetCommunityTax(ctx context.Context) (math.LegacyDec, error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	return params.CommunityTax, nil
}

// GetWithdrawAddrEnabled returns the current distribution withdraw address
// enabled parameter.
func (k Keeper) GetWithdrawAddrEnabled(ctx context.Context) (enabled bool, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return false, err
	}

	return params.WithdrawAddrEnabled, nil
}

// GetBurnRate returns the current burn rate on fee distribution
func (k Keeper) GetBurnRate(ctx context.Context) (dec math.LegacyDec, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	return params.BurnRate, nil
}

// GetValidatorCommissionRate returns the current validator commission rate
func (k Keeper) GetValidatorCommissionRate(ctx context.Context) (dec math.LegacyDec, err error) {
	params, err := k.Params.Get(ctx)
	if err != nil {
		return math.LegacyZeroDec(), err
	}

	return params.ValidatorCommissionRate, nil
}

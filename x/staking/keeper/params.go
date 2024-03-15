package keeper

import (
	"context"
	"time"

	"cosmossdk.io/math"

	"github.com/axiome-pro/axm-node/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UnbondingTime - The time duration for unbonding
func (k Keeper) UnbondingTime(ctx context.Context) (time.Duration, error) {
	params, err := k.GetParams(ctx)
	return params.UnbondingTime, err
}

// RedelegationTime - The time duration for unbonding
func (k Keeper) RedelegationTime(ctx context.Context) (time.Duration, error) {
	params, err := k.GetParams(ctx)
	return params.RedelegationTime, err
}

// MaxValidators - Maximum number of validators
func (k Keeper) MaxValidators(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	return params.MaxValidators, err
}

// MaxEntries - Maximum number of simultaneous unbonding
// delegations or redelegations (per pair/trio)
func (k Keeper) MaxEntries(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	return params.MaxEntries, err
}

// HistoricalEntries = number of historical info entries
// to persist in store
func (k Keeper) HistoricalEntries(ctx context.Context) (uint32, error) {
	params, err := k.GetParams(ctx)
	return params.HistoricalEntries, err
}

// BondDenom - Bondable coin denomination
func (k Keeper) BondDenom(ctx context.Context) (string, error) {
	params, err := k.GetParams(ctx)
	return params.BondDenom, err
}

// MinSelfStake - min self stake of validator
func (k Keeper) MinSelfDelegation(ctx context.Context) (math.Int, error) {
	params, err := k.GetParams(ctx)
	return params.MinSelfDelegation, err
}

// PowerReduction - is the amount of staking tokens required for 1 unit of consensus-engine power.
// Currently, this returns a global variable that the app developer can tweak.
// TODO: we might turn this into an on-chain param:
// https://github.com/cosmos/cosmos-sdk/issues/8365
func (k Keeper) PowerReduction(ctx context.Context) math.Int {
	return sdk.DefaultPowerReduction
}

func (k Keeper) EmissionTable(ctx context.Context) ([]*types.EmissionRange, error) {
	params, err := k.GetParams(ctx)
	return params.EmissionTable, err
}

// ValidatorEmissionRate - validator to delegator emission rate
func (k Keeper) ValidatorEmissionRate(ctx context.Context) (math.LegacyDec, error) {
	params, err := k.GetParams(ctx)
	return params.ValidatorEmissionRate, err
}

func (k Keeper) MaximumMonthlyPoints(ctx context.Context) (uint64, error) {
	params, err := k.GetParams(ctx)
	return params.MaximumMonthlyPoints, err
}

// SetParams sets the x/staking module parameters.
// CONTRACT: This method performs no validation of the parameters.
func (k Keeper) SetParams(ctx context.Context, params types.Params) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&params)
	if err != nil {
		return err
	}
	return store.Set(types.ParamsKey, bz)
}

// GetParams gets the x/staking module parameters.
func (k Keeper) GetParams(ctx context.Context) (params types.Params, err error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ParamsKey)
	if err != nil {
		return params, err
	}

	if bz == nil {
		return params, nil
	}

	err = k.cdc.Unmarshal(bz, &params)
	return params, err
}

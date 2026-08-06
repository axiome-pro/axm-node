package app

import (
	"context"
	"time"

	axmstakingkeeper "github.com/axiome-pro/axm-node/x/staking/keeper"
	axmstakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
	sdkstakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// ibcStakingKeeper adapts the forked Axiome staking keeper to the interface the
// IBC client keeper expects (ibc-go/modules/core/02-client/types.StakingKeeper).
//
// Axiome forked x/staking into its own proto package, so GetHistoricalInfo
// returns axiome.staking.v1beta1.HistoricalInfo while IBC is typed against
// cosmos.staking.v1beta1.HistoricalInfo. The two messages are structurally the
// same apart from the validator commission, which Axiome does not have: rewards
// are distributed through the referral module instead. Everything this adapter
// touches is a plain field copy.
//
// For the record, ibc-go v8 only reads HistoricalInfo.Header (timestamp, app
// hash and next validators hash) in GetSelfConsensusState, and never looks at
// the validator set. The valset is still converted below so the adapter stays
// correct if that ever changes.
type ibcStakingKeeper struct {
	keeper *axmstakingkeeper.Keeper
}

func newIBCStakingKeeper(k *axmstakingkeeper.Keeper) ibcStakingKeeper {
	return ibcStakingKeeper{keeper: k}
}

// GetHistoricalInfo returns the historical info at the given height, translated
// into the cosmos-sdk representation.
func (a ibcStakingKeeper) GetHistoricalInfo(ctx context.Context, height int64) (sdkstakingtypes.HistoricalInfo, error) {
	histInfo, err := a.keeper.GetHistoricalInfo(ctx, height)
	if err != nil {
		return sdkstakingtypes.HistoricalInfo{}, err
	}

	return sdkstakingtypes.HistoricalInfo{
		Header: histInfo.Header,
		Valset: convertValidators(histInfo.Valset),
	}, nil
}

// UnbondingTime returns the staking unbonding period. IBC compares it against
// the unbonding period declared by a counterparty light client.
func (a ibcStakingKeeper) UnbondingTime(ctx context.Context) (time.Duration, error) {
	return a.keeper.UnbondingTime(ctx)
}

func convertValidators(vals []axmstakingtypes.Validator) []sdkstakingtypes.Validator {
	if len(vals) == 0 {
		return nil
	}

	out := make([]sdkstakingtypes.Validator, 0, len(vals))
	for _, v := range vals {
		out = append(out, sdkstakingtypes.Validator{
			OperatorAddress: v.OperatorAddress,
			ConsensusPubkey: v.ConsensusPubkey,
			Jailed:          v.Jailed,
			// The BondStatus enums are declared with identical numeric values in
			// both packages (0 unspecified, 1 unbonded, 2 unbonding, 3 bonded).
			Status:          sdkstakingtypes.BondStatus(v.Status),
			Tokens:          v.Tokens,
			DelegatorShares: v.DelegatorShares,
			Description: sdkstakingtypes.Description{
				Moniker:         v.Description.Moniker,
				Identity:        v.Description.Identity,
				Website:         v.Description.Website,
				SecurityContact: v.Description.SecurityContact,
				Details:         v.Description.Details,
			},
			UnbondingHeight: v.UnbondingHeight,
			UnbondingTime:   v.UnbondingTime,
			// Axiome validators have no commission. A zero commission is the
			// truthful representation rather than a placeholder.
			Commission:              sdkstakingtypes.Commission{},
			MinSelfDelegation:       v.MinSelfDelegation,
			UnbondingOnHoldRefCount: v.UnbondingOnHoldRefCount,
			UnbondingIds:            v.UnbondingIds,
		})
	}

	return out
}

package types

import (
	distrtypes "github.com/axiome-pro/axm-node/x/distribution/types"
	"github.com/axiome-pro/axm-node/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"context"
	//nolint:staticcheck
)

type StakingKeeper interface {
	// BondDenom - Bondable coin denomination
	BondDenom(ctx context.Context) (string, error)
	// GetValidator get a single validator
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator types.Validator, err error)
	// GetBondedValidatorsByPower get the current group of bonded validators sorted by power-rank
	GetBondedValidatorsByPower(ctx context.Context) ([]types.Validator, error)
	// GetAllDelegatorDelegations return all delegations for a delegator
	GetAllDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress) ([]types.Delegation, error)
	// GetDelegation return a specific delegation
	GetDelegation(ctx context.Context,
		delAddr sdk.AccAddress, valAddr sdk.ValAddress) (types.Delegation, error)
	// HasReceivingRedelegation check if validator is receiving a redelegation
	HasReceivingRedelegation(ctx context.Context,
		delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) (bool, error)
}

type DistributionKeeper interface {
	DelegatorWithdrawAddress(c context.Context, req *distrtypes.QueryDelegatorWithdrawAddressRequest) (*distrtypes.QueryDelegatorWithdrawAddressResponse, error)
	DelegationRewards(c context.Context, req *distrtypes.QueryDelegationRewardsRequest) (*distrtypes.QueryDelegationRewardsResponse, error)
	DelegationTotalRewards(c context.Context, req *distrtypes.QueryDelegationTotalRewardsRequest) (*distrtypes.QueryDelegationTotalRewardsResponse, error)
	DelegatorValidators(c context.Context, req *distrtypes.QueryDelegatorValidatorsRequest) (*distrtypes.QueryDelegatorValidatorsResponse, error)
}

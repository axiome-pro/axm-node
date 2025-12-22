package types

import (
	"context"

	axdisttypes "github.com/axiome-pro/axm-node/x/distribution/types"
	"github.com/axiome-pro/axm-node/x/staking/types"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	disttypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

type WrappedStakingKeeper struct {
	Keeper StakingKeeper
}

func (w WrappedStakingKeeper) BondDenom(ctx context.Context) (string, error) {
	return w.Keeper.BondDenom(ctx)
}

func (w WrappedStakingKeeper) GetValidator(ctx context.Context, addr sdk.ValAddress) (stakingtypes.Validator, error) {
	val, err := w.Keeper.GetValidator(ctx, addr)
	if err != nil {
		return stakingtypes.Validator{}, err
	}
	return convertValidator(val), nil
}

func (w WrappedStakingKeeper) GetBondedValidatorsByPower(ctx context.Context) ([]stakingtypes.Validator, error) {
	validators, err := w.Keeper.GetBondedValidatorsByPower(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]stakingtypes.Validator, len(validators))
	for i, val := range validators {
		result[i] = convertValidator(val)
	}
	return result, nil
}

func (w WrappedStakingKeeper) GetAllDelegatorDelegations(ctx context.Context, delegator sdk.AccAddress) ([]stakingtypes.Delegation, error) {
	delegations, err := w.Keeper.GetAllDelegatorDelegations(ctx, delegator)
	if err != nil {
		return nil, err
	}

	result := make([]stakingtypes.Delegation, len(delegations))
	for i, del := range delegations {
		result[i] = stakingtypes.Delegation{
			DelegatorAddress: del.DelegatorAddress,
			ValidatorAddress: del.ValidatorAddress,
			Shares:           del.Shares,
		}
	}
	return result, nil
}

func (w WrappedStakingKeeper) GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (stakingtypes.Delegation, error) {
	delegation, err := w.Keeper.GetDelegation(ctx, delAddr, valAddr)
	if err != nil {
		return stakingtypes.Delegation{}, err
	}

	return stakingtypes.Delegation{
		DelegatorAddress: delegation.DelegatorAddress,
		ValidatorAddress: delegation.ValidatorAddress,
		Shares:           delegation.Shares,
	}, nil
}

func (w WrappedStakingKeeper) HasReceivingRedelegation(ctx context.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) (bool, error) {
	return w.Keeper.HasReceivingRedelegation(ctx, delAddr, valDstAddr)
}

func convertValidator(val types.Validator) stakingtypes.Validator {
	return stakingtypes.Validator{
		OperatorAddress:         val.OperatorAddress,
		ConsensusPubkey:         val.ConsensusPubkey,
		Jailed:                  val.Jailed,
		Status:                  stakingtypes.BondStatus(val.Status),
		Tokens:                  val.Tokens,
		DelegatorShares:         val.DelegatorShares,
		Description:             convertDescription(val.Description),
		UnbondingHeight:         val.UnbondingHeight,
		UnbondingTime:           val.UnbondingTime,
		Commission:              stakingtypes.NewCommission(math.LegacyZeroDec(), math.LegacyZeroDec(), math.LegacyZeroDec()),
		MinSelfDelegation:       val.MinSelfDelegation,
		UnbondingOnHoldRefCount: val.UnbondingOnHoldRefCount,
		UnbondingIds:            val.UnbondingIds,
	}
}

func convertDescription(desc types.Description) stakingtypes.Description {
	return stakingtypes.Description{
		Moniker:         desc.Moniker,
		Identity:        desc.Identity,
		Website:         desc.Website,
		SecurityContact: desc.SecurityContact,
		Details:         desc.Details,
	}
}

type WrappedDistributionKeeper struct {
	Keeper DistributionKeeper
}

func (w WrappedDistributionKeeper) DelegatorWithdrawAddress(c context.Context, req *disttypes.QueryDelegatorWithdrawAddressRequest) (*disttypes.QueryDelegatorWithdrawAddressResponse, error) {
	if req == nil {
		return nil, nil
	}
	axReq := &axdisttypes.QueryDelegatorWithdrawAddressRequest{DelegatorAddress: req.DelegatorAddress}
	axRes, err := w.Keeper.DelegatorWithdrawAddress(c, axReq)
	if err != nil {
		return nil, err
	}
	return &disttypes.QueryDelegatorWithdrawAddressResponse{WithdrawAddress: axRes.WithdrawAddress}, nil
}

func (w WrappedDistributionKeeper) DelegationRewards(c context.Context, req *disttypes.QueryDelegationRewardsRequest) (*disttypes.QueryDelegationRewardsResponse, error) {
	if req == nil {
		return nil, nil
	}
	axReq := &axdisttypes.QueryDelegationRewardsRequest{
		DelegatorAddress: req.DelegatorAddress,
		ValidatorAddress: req.ValidatorAddress,
	}
	axRes, err := w.Keeper.DelegationRewards(c, axReq)
	if err != nil {
		return nil, err
	}
	// Combine rewards and emitted to preserve total value in SDK response
	totalRewards := sdk.DecCoins{}
	if axRes.Rewards != nil {
		totalRewards = totalRewards.Add(axRes.Rewards...)
	}
	if axRes.Emitted != nil {
		totalRewards = totalRewards.Add(axRes.Emitted...)
	}
	return &disttypes.QueryDelegationRewardsResponse{Rewards: totalRewards}, nil
}

func (w WrappedDistributionKeeper) DelegationTotalRewards(c context.Context, req *disttypes.QueryDelegationTotalRewardsRequest) (*disttypes.QueryDelegationTotalRewardsResponse, error) {
	if req == nil {
		return nil, nil
	}
	axReq := &axdisttypes.QueryDelegationTotalRewardsRequest{DelegatorAddress: req.DelegatorAddress}
	axRes, err := w.Keeper.DelegationTotalRewards(c, axReq)
	if err != nil {
		return nil, err
	}
	// Map each reward, summing reward and emission for SDK type
	mapped := make([]disttypes.DelegationDelegatorReward, 0, len(axRes.Rewards))
	for _, r := range axRes.Rewards {
		combined := sdk.DecCoins{}
		if r.Reward != nil {
			combined = combined.Add(r.Reward...)
		}
		if r.Emission != nil {
			combined = combined.Add(r.Emission...)
		}
		mapped = append(mapped, disttypes.DelegationDelegatorReward{
			ValidatorAddress: r.ValidatorAddress,
			Reward:           combined,
		})
	}
	// Total from axRes already reflects both reward and emission
	return &disttypes.QueryDelegationTotalRewardsResponse{Rewards: mapped, Total: axRes.Total}, nil
}

func (w WrappedDistributionKeeper) DelegatorValidators(c context.Context, req *disttypes.QueryDelegatorValidatorsRequest) (*disttypes.QueryDelegatorValidatorsResponse, error) {
	if req == nil {
		return nil, nil
	}
	axReq := &axdisttypes.QueryDelegatorValidatorsRequest{DelegatorAddress: req.DelegatorAddress}
	axRes, err := w.Keeper.DelegatorValidators(c, axReq)
	if err != nil {
		return nil, err
	}
	return &disttypes.QueryDelegatorValidatorsResponse{Validators: axRes.Validators}, nil
}

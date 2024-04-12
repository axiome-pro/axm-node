package app

import (
	"github.com/axiome-pro/axm-node/x/referral"
	"github.com/axiome-pro/axm-node/x/referral/keeper"
	stakingkeeper "github.com/axiome-pro/axm-node/x/staking/keeper"
	stakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
	"context"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"

	referraltypes "github.com/axiome-pro/axm-node/x/referral/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const UpgradeNameV102 = "v1.0.2"
const UpgradeNameV103 = "v1.0.3"

func (app *AxmApp) RegisterUpgradeHandlers() {
	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeNameV102,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			err := upgradeToV102(ctx, app.ReferralKeeper)
			if err != nil {
				return nil, err
			}
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)

	app.UpgradeKeeper.SetUpgradeHandler(
		UpgradeNameV103,
		func(ctx context.Context, _ upgradetypes.Plan, fromVM module.VersionMap) (module.VersionMap, error) {
			err := upgradeToV103(ctx, *app.StakingKeeper, app.AccountKeeper)
			if err != nil {
				return nil, err
			}
			return app.ModuleManager.RunMigrations(ctx, app.Configurator(), fromVM)
		},
	)
}

func upgradeToV102(ctx context.Context, k referral.Keeper) error {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger := k.Logger(sdkCtx)
	logger.Info("Starting account teams aggregation refresh ...")
	k.Iterate(sdkCtx, func(acc string, info *referraltypes.Info) (changed, checkForStatusUpdate bool) {
		oldACtiveCount := info.ActiveCount
		info.ActiveCount = &referraltypes.ActiveAggregations{}
		// Save 3x0 and 3x3 status aggregations, because it's not error on it
		info.ActiveCount.FirstLine = oldACtiveCount.FirstLine
		info.ActiveCount.FirstLineBy3 = oldACtiveCount.FirstLineBy3

		// Find all active children and recalculates it's team sizes
		children, err := k.GetChildren(sdkCtx, acc)
		if err != nil {
			panic(err)
		}
		for _, rAddr := range children {
			rInfo, err := k.Get(sdkCtx, rAddr)
			if err != nil {
				logger.Error("Account %s not found", rAddr)
				panic(err)
			}

			// If account is active, it impacts team aggregations
			if rInfo.Active {
				keeper.ChangeTeamActive(info.ActiveCount, rInfo.GetActiveRefsCountFromLevelToLevel(1, 13), 1)
			}
		}

		if !oldACtiveCount.Eqals(*info.ActiveCount) {
			logger.Info("Status aggregations updated", "acc", acc, "old", oldACtiveCount, "new", info.ActiveCount)
			return true, true
		}

		return false, false
	})
	logger.Info("... done")
	return nil
}

// warn: this upgrade complete all pending redelegations
func upgradeToV103(ctx context.Context, k stakingkeeper.Keeper, accountKeeper authkeeper.AccountKeeper) error {
	logger := k.Logger(ctx)
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	logger.Info("Starting redelegation queue cleaning ...")
	unbondingTime, err := k.UnbondingTime(ctx)
	if err != nil {
		panic(err)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations, err := k.DequeueAllMatureRedelegationQueue(ctx, sdkCtx.BlockTime().Add(unbondingTime))
	if err != nil {
		panic(err)
	}

	for _, dvvTriplet := range matureRedelegations {
		addressCodec := k.ValidatorAddressCodec()
		valSrcAddr, err := addressCodec.StringToBytes(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := addressCodec.StringToBytes(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress, err := accountKeeper.AddressCodec().StringToBytes(dvvTriplet.DelegatorAddress)
		if err != nil {
			panic(err)
		}

		red, err := k.GetRedelegation(ctx, delegatorAddress, valSrcAddr, valDstAddr)
		if err != nil {
			continue
		}

		// Fix completion time for all entries
		for i := 0; i < len(red.Entries); i++ {
			red.Entries[i].CompletionTime = sdkCtx.BlockTime()
		}

		// Save redelegation entries
		err = k.SetRedelegation(ctx, red)

		if err != nil {
			panic(err)
		}

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)

		logger.Info("... processing ", "del", dvvTriplet.DelegatorAddress, "err", err)

		if err != nil {
			continue
		}

		sdkCtx.EventManager().EmitEvent(
			sdk.NewEvent(
				stakingtypes.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(stakingtypes.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(stakingtypes.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(stakingtypes.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	logger.Info("... done")
	return nil
}

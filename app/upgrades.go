package app

import (
	"github.com/axiome-pro/axm-node/x/referral"
	"github.com/axiome-pro/axm-node/x/referral/keeper"
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	upgradetypes "cosmossdk.io/x/upgrade/types"

	referraltypes "github.com/axiome-pro/axm-node/x/referral/types"
	"github.com/cosmos/cosmos-sdk/types/module"
)

const UpgradeNameV102 = "v1.0.2"

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

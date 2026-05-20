package keeper

import (
	"time"

	"github.com/axiome-pro/axm-node/util"
	"github.com/axiome-pro/axm-node/x/referral/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpgradeDeactivateBelowThreshold iterates over all referral accounts,
// deactivates those with SelfDelegated < 1000 axm (new threshold),
// and schedules status downgrades for ancestors whose status requirements
// are no longer met.
func (k Keeper) UpgradeDeactivateBelowThreshold(ctx sdk.Context) error {
	logger := k.Logger(ctx)
	logger.Info("Starting activation threshold upgrade to 1000 axm ...")

	// Phase 1: Deactivate accounts below new threshold
	var deactivated int64
	k.Iterate(ctx, func(acc string, info *types.Info) (changed, checkForStatusUpdate bool) {
		if !info.Active {
			return false, false
		}
		// New threshold: 1000 axm = 1_000_000_000 uaxm
		if info.SelfDelegated.Int64() < 1_000_000_000 {
			bu := newBunchUpdater(k, ctx)
			err := k.SetActive(ctx, acc, false, false, bu)
			if err != nil {
				logger.Error("SetActive(false) failed", "acc", acc, "error", err)
				return false, false
			}
			if err := bu.commit(); err != nil {
				logger.Error("commit failed after deactivation", "acc", acc, "error", err)
				return false, false
			}
			deactivated++
			logger.Info("Deactivated account", "acc", acc, "delegated", info.SelfDelegated)
		}
		return false, false
	})

	logger.Info("Deactivation phase done", "deactivated", deactivated)

	// Phase 2: Check all accounts for status requirements and schedule downgrades
	var downgradeCounter int64
	k.Iterate(ctx, func(acc string, info *types.Info) (changed, checkForStatusUpdate bool) {
		if info.Status == types.STATUS_UNSPECIFIED || info.Status == types.STATUS_NEW {
			return false, false
		}

		checkResult, err := checkStatusRequirements(info.Status, *info)
		if err != nil {
			logger.Error("checkStatusRequirements failed", "acc", acc, "error", err)
			return false, false
		}

		if !checkResult.Overall {
			// Requirements not met — check if downgrade is already scheduled
			if info.StatusDowngradeAt != nil {
				// Already scheduled — do nothing
				return false, false
			}

			// Schedule downgrade: upgrade time + 1 hour + N seconds
			downgradeAt := ctx.BlockTime().Add(1*time.Hour + time.Duration(downgradeCounter)*time.Second)
			downgradeCounter++

			if err := k.ScheduleStatusDowngrade(ctx, acc, downgradeAt); err != nil {
				logger.Error("ScheduleStatusDowngrade failed", "acc", acc, "error", err)
				return false, false
			}

			info.StatusDowngradeAt = &downgradeAt

			util.EmitEvent(ctx,
				&types.EventStatusWillBeDowngraded{
					Address: acc,
					Time:    downgradeAt,
				},
			)

			logger.Info("Scheduled status downgrade", "acc", acc, "at", downgradeAt)
			return true, false
		}

		return false, false
	})

	logger.Info("... upgrade done", "deactivated", deactivated, "downgrades_scheduled", downgradeCounter)
	return nil
}

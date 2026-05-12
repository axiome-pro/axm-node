package keeper

import (
	"time"

	"github.com/axiome-pro/axm-node/util"
	"github.com/axiome-pro/axm-node/x/referral/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// UpgradeRecalculateStatuses iterates over all referral accounts, checks status
// requirements with the current (updated) rules, and schedules downgrades with
// 1-second intervals or promotes accounts that now qualify for a higher status.
func (k Keeper) UpgradeRecalculateStatuses(ctx sdk.Context) error {
	logger := k.Logger(ctx)
	logger.Info("Starting status recalculation for v2.2.0 upgrade ...")

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
			// Requirements not met — schedule downgrade

			// Remove existing downgrade schedule if present
			if info.StatusDowngradeAt != nil {
				if err := k.RemoveStatusDowngradeSchedule(ctx, acc, *info.StatusDowngradeAt); err != nil {
					logger.Error("RemoveStatusDowngradeSchedule failed", "acc", acc, "error", err)
					return false, false
				}
			}

			downgradeAt := ctx.BlockTime().Add(7*24*time.Hour + time.Duration(downgradeCounter)*time.Second)
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

		// Requirements met — cancel pending downgrade if any, and check for promotion
		if info.StatusDowngradeAt != nil {
			if err := k.RemoveStatusDowngradeSchedule(ctx, acc, *info.StatusDowngradeAt); err != nil {
				logger.Error("RemoveStatusDowngradeSchedule failed", "acc", acc, "error", err)
				return false, false
			}
			info.StatusDowngradeAt = nil

			util.EmitEvent(ctx,
				&types.EventStatusDowngradeCanceled{
					Address: acc,
				},
			)
		}

		// Check for possible promotion
		nextStatus := info.Status
		for {
			if nextStatus == types.MaximumStatus {
				break
			}
			nextStatus++

			cr, err := checkStatusRequirements(nextStatus, *info)
			if err != nil {
				logger.Error("checkStatusRequirements for promotion failed", "acc", acc, "error", err)
				nextStatus--
				break
			}
			if !cr.Overall {
				nextStatus--
				break
			}
		}

		if nextStatus > info.Status {
			util.EmitEvent(ctx,
				&types.EventStatusUpdated{
					Address: acc,
					Before:  info.Status,
					After:   nextStatus,
				},
			)
			k.setStatus(ctx, info, nextStatus, acc)
			logger.Info("Status promoted", "acc", acc, "from", info.Status, "to", nextStatus)
			return true, false
		}

		return false, false
	})

	logger.Info("... status recalculation done", "downgrades_scheduled", downgradeCounter)
	return nil
}

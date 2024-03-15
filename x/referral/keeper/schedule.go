package keeper

import (
	"github.com/axiome-pro/axm-node/x/referral/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

func (k Keeper) ScheduleStatusDowngrade(ctx sdk.Context, acc string, downgradeAt time.Time) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Set(types.GetDowngradeQueueKey(acc, downgradeAt), []byte{0x01})
}

func (k Keeper) RemoveStatusDowngradeSchedule(ctx sdk.Context, acc string, downgradeAt time.Time) error {
	store := k.storeService.OpenKVStore(ctx)
	return store.Delete(types.GetDowngradeQueueKey(acc, downgradeAt))
}

func (k Keeper) PerfomStatusDowngradeSchedule(ctx sdk.Context) error {
	store := k.storeService.OpenKVStore(ctx)

	itr, err := store.Iterator(types.GetDowngradeQueueIteratorStartKey(), types.GetDowngradeQueueIteratorEndKey(ctx.BlockTime()))
	if err != nil {
		return nil
	}

	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		acc := types.ExtractAccFromDowngradeQueueKey(itr.Key())
		k.Logger(ctx).Info("Downgrade status", "acc", acc)
		err = k.performDowngrade(ctx, acc)
		if err != nil {
			return err
		}
		err = store.Delete(itr.Key())
		if err != nil {
			return err
		}
	}

	return nil
}

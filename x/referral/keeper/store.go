package keeper

import (
	"github.com/axiome-pro/axm-node/x/referral/types"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

func (k Keeper) Iterate(ctx sdk.Context, callback func(acc string, r *types.Info) (changed, checkForStatusUpdate bool)) {
	bu := newBunchUpdater(k, ctx)
	store := k.storeService.OpenKVStore(ctx)
	it, err := store.Iterator(types.InfoPrefix, storetypes.PrefixEndBytes(types.InfoPrefix))

	// TODO: we seriously need to panic here?
	if err != nil {
		panic(err)
	}

	defer func() {
		if it != nil {
			it.Close()
		}
	}()
	for ; it.Valid(); it.Next() {
		var acc = string(it.Key()[len(types.InfoPrefix):])
		var item types.Info
		if err := k.cdc.Unmarshal(it.Value(), &item); err != nil {
			panic(errors.Wrapf(err, `cannot unmarshal info for "%s"`, acc))
		}
		if changed, checkForStatusUpdate := callback(acc, &item); changed || checkForStatusUpdate {
			var f func(r *types.Info) error
			if changed {
				f = func(r *types.Info) error {
					*r = item
					return nil
				}
			} else {
				f = func(_ *types.Info) error {
					return nil
				}
			}
			err := bu.update(acc, checkForStatusUpdate, f)
			if err != nil {
				panic(err)
			}
		}
	}
	it.Close()
	it = nil
	err = bu.commit()
	if err != nil {
		panic(err)
	}
}

// Get returns all the data for an account (status, parent, children)
func (k Keeper) Get(ctx sdk.Context, acc string) (types.Info, error) {
	store := k.storeService.OpenKVStore(ctx)
	var item types.Info

	bz, err := store.Get(types.GetInfoAddrKey(acc))

	if err != nil {
		err = errors.Wrapf(
			err,
			"no data for %s", acc,
		)
	}

	err = errors.Wrapf(
		k.cdc.Unmarshal(bz, &item),
		"no data for %s", acc,
	)
	return item, err
}

func (k Keeper) set(ctx sdk.Context, acc string, value types.Info) error {
	store := k.storeService.OpenKVStore(ctx)
	keyBytes := types.GetInfoAddrKey(acc)
	valueBytes, err := k.cdc.Marshal(&value)
	if err != nil {
		return err
	}
	return store.Set(keyBytes, valueBytes)
}

func (k Keeper) exists(ctx sdk.Context, acc string) bool {
	store := k.storeService.OpenKVStore(ctx)
	keyBytes := types.GetInfoAddrKey(acc)

	present, err := store.Has(keyBytes)
	if err != nil {
		return false
	}
	return present
}

func (k Keeper) update(ctx sdk.Context, acc string, callback func(value types.Info) types.Info) error {
	store := k.storeService.OpenKVStore(ctx)

	keyBytes := types.GetInfoAddrKey(acc)
	var value types.Info
	bz, err := store.Get(keyBytes)
	if err != nil {
		return err
	}

	err = k.cdc.Unmarshal(bz, &value)
	if err != nil {
		return err
	}
	value = callback(value)
	valueBytes, err := k.cdc.Marshal(&value)
	if err != nil {
		return err
	}
	store.Set(keyBytes, valueBytes)
	return nil
}

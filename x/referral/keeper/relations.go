package keeper

import (
	"github.com/axiome-pro/axm-node/x/referral/types"
	sdkioerrors "cosmossdk.io/errors"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/pkg/errors"
)

// GetTopLevelAccounts returns all accounts without parents and is supposed to be used during genesis export
func (k Keeper) GetTopLevelAccounts(ctx sdk.Context) (topLevel []string, err error) {
	store := k.storeService.OpenKVStore(ctx)
	itr, err := store.Iterator(types.InfoPrefix, storetypes.PrefixEndBytes(types.InfoPrefix))

	if err != nil {
		return nil, err
	}

	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		v := itr.Value()
		var record types.Info
		err = k.cdc.Unmarshal(v, &record)
		if err != nil {
			return nil, err
		}
		addr := types.ParseInfoAddrKey(itr.Key())

		if record.Referrer == "" {
			topLevel = append(topLevel, addr)
		}
	}
	return topLevel, nil
}

// AddTopLevelAccount adds accounts without parent and is supposed to be used during genesis
func (k Keeper) AddTopLevelAccount(ctx sdk.Context, acc string, status types.Status) (err error) {
	k.Logger(ctx).Debug("AddTopLevelAccount", "acc", acc)
	defer func() {
		if e := recover(); e != nil {
			k.Logger(ctx).Error("AddTopLevelAccount paniced", "err", e)
			if er, ok := e.(error); ok {
				err = errors.Wrap(er, "AddTopLevelAccount paniced")
			} else {
				err = errors.Errorf("AddTopLevelAccount paniced: %s", e)
			}
		}
	}()
	if k.exists(ctx, acc) {
		return sdkioerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("account %s already exists", acc),
		)
	}

	bu := newBunchUpdater(k, ctx)

	newItem := types.NewInfoWithStatus("", math.ZeroInt(), status)
	if err = bu.set(acc, newItem); err != nil {
		return err
	}
	if err = bu.commit(); err != nil {
		return err
	}
	return nil
}

func (k Keeper) setBackRelation(ctx sdk.Context, parentAcc, childAcc string) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetReferralsRelationKey(parentAcc, childAcc)
	return store.Set(key, []byte{0x01})
}

// AppendChild adds a new account to the referral structure. The parent account should already exist and the child one
// should not.
func (k Keeper) AppendChild(ctx sdk.Context, parentAcc string, childAcc string) error {
	return k.appendChild(ctx, parentAcc, childAcc, false, types.STATUS_NEW)
}
func (k Keeper) appendChild(ctx sdk.Context, parentAcc string, childAcc string, skipActivityCheck bool, status types.Status) error {
	if parentAcc == "" {
		return types.ErrParentNil
	}
	if k.exists(ctx, childAcc) {
		return sdkioerrors.Wrap(
			sdkerrors.ErrInvalidRequest,
			fmt.Sprintf("account %s already exists", childAcc),
		)
	}
	var (
		bu        = newBunchUpdater(k, ctx)
		anc       = parentAcc
		delegated = math.ZeroInt()
	)
	newItem := types.NewInfoWithStatus(parentAcc, delegated, status)

	err := bu.set(childAcc, newItem)
	if err != nil {
		return sdkioerrors.Wrap(err, "cannot set "+childAcc)
	}

	var registrationClosed bool
	err = bu.update(parentAcc, true, func(value *types.Info) error {
		newDelegated := (*value.TeamDelegated).Add(delegated)
		value.TeamDelegated = &newDelegated
		bu.addCallback(StakeChangedCallback, anc)
		err = k.setBackRelation(ctx, parentAcc, childAcc)
		if err != nil {
			return sdkioerrors.Wrap(err, "cannot set relation for "+parentAcc+" "+childAcc)
		}
		anc = value.Referrer
		if !skipActivityCheck {
			registrationClosed = false
		}
		return nil
	})
	if err != nil {
		return sdkioerrors.Wrap(err, "cannot update "+anc)
	}
	if registrationClosed {
		return types.ErrRegistrationClosed
	}

	for i := 1; i < 14; i++ {
		if anc == "" {
			break
		}
		err = bu.update(anc, true, func(value *types.Info) error {
			newDelegated := (*value.TeamDelegated).Add(delegated)
			value.TeamDelegated = &newDelegated
			bu.addCallback(StakeChangedCallback, anc)
			anc = value.Referrer
			return nil
		})
		if err != nil {
			return sdkioerrors.Wrap(err, "cannot update "+anc)
		}
	}

	if err := bu.commit(); err != nil {
		return sdkioerrors.Wrap(err, "cannot commit")
	}
	return nil
}

// GetParent returns a parent for an account
func (k Keeper) GetParent(ctx sdk.Context, acc string) (string, error) {
	data, err := k.Get(ctx, acc)
	if err != nil {
		return "", errors.Wrap(err, "cannot obtain data")
	}
	return data.Referrer, nil
}

func (k Keeper) GetChildren(ctx sdk.Context, acc string) (children []string, err error) {
	store := k.storeService.OpenKVStore(ctx)

	iteratorKey := types.GetReferralsChildIteratorKey(acc)

	itr, err := store.Iterator(iteratorKey, storetypes.PrefixEndBytes(iteratorKey))

	if err != nil {
		return nil, err
	}

	defer itr.Close()
	for ; itr.Valid(); itr.Next() {
		_, child, err := types.ParseReferralFromReleationKey(itr.Key())
		if err != nil {
			return nil, err
		}

		children = append(children, child)
	}

	return children, nil
}

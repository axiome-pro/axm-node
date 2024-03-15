package keeper

import (
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axiome-pro/axm-node/x/referral/types"
)

func (k Keeper) ExportToGenesis(ctx sdk.Context) (*types.GenesisState, error) {
	var (
		data       types.Info
		err        error
		params     types.Params
		topLevel   []string
		other      []types.Refs
		downgrades []types.Downgrade

		childrenAddr []string
		children     []*types.RefInfo
		thisLevel    []types.Refs
		nextLevel    []types.Refs
	)
	params = k.GetParams(ctx)
	topLevel, err = k.GetTopLevelAccounts(ctx)
	topLevelRefs := make([]*types.RefInfo, len(topLevel))
	if err != nil {
		return nil, err
	}

	for i, addr := range topLevel {
		data, err = k.Get(ctx, addr)
		if err != nil {
			return nil, err
		}
		if data.StatusDowngradeAt != nil {
			downgrades = append(downgrades, *types.NewDowngrade(addr, data.Status, *data.StatusDowngradeAt))
		}
		childrenAddr, err = k.GetChildren(ctx, addr)
		if err != nil {
			return nil, err
		}

		topLevelRefs[i] = types.NewRefInfo(addr, data.Status)

		if len(childrenAddr) == 0 {
			continue
		}

		children = make([]*types.RefInfo, len(childrenAddr))

		for i, childAddr := range childrenAddr {
			child, err := k.Get(ctx, childAddr)
			if err != nil {
				return nil, err
			}

			children[i] = types.NewRefInfo(childAddr, child.Status)
		}

		nextLevel = append(nextLevel, *types.NewRefs(addr, children))
	}
	for len(nextLevel) != 0 {
		other = append(other, nextLevel...)
		thisLevel = nextLevel
		nextLevel = nil
		for _, r := range thisLevel {
			for _, refInfo := range r.Referrals {
				data, err = k.Get(ctx, refInfo.Address)
				if err != nil {
					return nil, errors.Wrapf(err, "cannot obtain %s data", refInfo.Address)
				}
				if data.StatusDowngradeAt != nil {
					downgrades = append(downgrades, *types.NewDowngrade(refInfo.Address, data.Status, *data.StatusDowngradeAt))
				}
				childrenAddr, err = k.GetChildren(ctx, refInfo.Address)
				if err != nil {
					return nil, err
				}
				if len(childrenAddr) == 0 {
					continue
				}

				children = make([]*types.RefInfo, len(childrenAddr))

				for i, childAddr := range childrenAddr {
					child, err := k.Get(ctx, childAddr)
					if err != nil {
						return nil, err
					}

					children[i] = types.NewRefInfo(childAddr, child.Status)
				}

				nextLevel = append(nextLevel, *types.NewRefs(refInfo.Address, children))
			}
		}
	}

	return types.NewGenesisState(params, topLevelRefs, other, downgrades), nil
}

func (k Keeper) ImportFromGenesis(
	ctx sdk.Context,
	topLevel []*types.RefInfo,
	otherAccounts []types.Refs,
	downgrades []types.Downgrade,
) error {
	k.Logger(ctx).Info("... top level accounts")
	for _, top := range topLevel {
		if err := k.AddTopLevelAccount(ctx, top.Address, top.Status); err != nil {
			panic(errors.Wrapf(err, "cannot add %s", top.Address))
		}
		k.Logger(ctx).Debug("account added", "acc", top.Address, "parent", nil)
	}
	k.Logger(ctx).Info("... other accounts")
	for _, r := range otherAccounts {
		for _, ref := range r.Referrals {
			if err := k.appendChild(ctx, r.Referrer, ref.Address, true, ref.Status); err != nil {
				panic(errors.Wrapf(err, "cannot add %s", ref.Address))
			}

			k.Logger(ctx).Debug("account added", "acc", ref.Address, "parent", r.Referrer)
		}
	}
	bu := newBunchUpdater(k, ctx)
	k.Logger(ctx).Info("... status downgrades")
	for _, x := range downgrades {
		if err := bu.update(x.Account, false, func(value *types.Info) error {
			k.Logger(ctx).Debug("status downgrade", "acc", x.Account, "from", x.Current, "to", value.Status)
			if value.StatusDowngradeAt != nil {
				err1 := k.RemoveStatusDowngradeSchedule(ctx, x.Account, *value.StatusDowngradeAt)
				if err1 != nil {
					return err1
				}
			}
			value.StatusDowngradeAt = &x.Time
			return k.ScheduleStatusDowngrade(ctx, x.Account, *value.StatusDowngradeAt)
		}); err != nil {
			return err
		}
	}
	k.Logger(ctx).Info("... persisting")
	if err := bu.commit(); err != nil {
		return err
	}
	return nil
}

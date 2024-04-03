package keeper

import (
	"cosmossdk.io/collections"
	"fmt"
	"github.com/pkg/errors"

	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/axiome-pro/axm-node/util"
	"github.com/axiome-pro/axm-node/x/referral/types"
	storetypes "cosmossdk.io/core/store"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
)

// Keeper of the referral store
type Keeper struct {
	cdc                 codec.BinaryCodec
	storeService        storetypes.KVStoreService
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	stakingKeeper       types.StakingKeeper
	Params              collections.Item[types.Params]
	eventHooks          map[string][]func(ctx sdk.Context, acc string) error
	authority           sdk.AccAddress
	referralAccountName string
}

// NewKeeper creates a referral keeper
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	stakingKeeper types.StakingKeeper,
	authority sdk.AccAddress, referralAccountName string,
) *Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	keeper := Keeper{
		cdc:                 cdc,
		storeService:        storeService,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		stakingKeeper:       stakingKeeper,
		Params:              collections.NewItem(sb, types.ParamsKey, "params", codec.CollValue[types.Params](cdc)),
		eventHooks:          make(map[string][]func(ctx sdk.Context, acc string) error),
		authority:           authority,
		referralAccountName: referralAccountName,
	}
	return &keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetStatus returns a status for an account (i.e. lvl 1 "Lucky", lvl 2 "Leader", lvl 3 "Master" or so on)
func (k Keeper) GetStatus(ctx sdk.Context, acc string) (types.Status, error) {
	data, err := k.Get(ctx, acc)
	if err != nil {
		return 0, err
	}
	return data.Status, nil
}

// GetReferralFeesForDelegating returns a set of account-ratio pairs, describing what part of being delegated funds
// should go to what wallet. 0.15 total. The rest should be burned.
func (k Keeper) GetReferralFeesForDelegating(ctx sdk.Context, acc string) ([]types.ReferralFee, util.Fraction, error) {
	return k.getReferralFeesCore(
		ctx,
		acc,
		k.GetParams(ctx).DelegatingAward.Network,
	)
}

// AreStatusRequirementsFulfilled validates if the account suffices the status requirement.
// The actual account status doesn't matter and won't be updated.
func (k Keeper) AreStatusRequirementsFulfilled(ctx sdk.Context, acc string, s types.Status) (types.StatusCheckResult, error) {
	if s < types.MinimumStatus || s > types.MaximumStatus {
		return types.StatusCheckResult{Overall: false}, fmt.Errorf("there is no such status: %d", s)
	}
	data, err := k.Get(ctx, acc)
	if err != nil {
		return types.StatusCheckResult{Overall: false}, err
	}
	return checkStatusRequirements(s, data)
}

// GetDelegatedInNetwork returns total amount of delegated coins in a person's network
// Own coins inclusive.
func (k Keeper) GetDelegatedInNetwork(ctx sdk.Context, acc string) (math.Int, error) {
	data, err := k.Get(ctx, acc)
	if err != nil {
		return math.Int{}, err
	}

	return (*data.TeamDelegated).Add(*data.SelfDelegated), nil
}

func (k Keeper) OnBalanceChanged(ctx sdk.Context, acc string, dd math.Int) error {
	k.Logger(ctx).Debug("OnBalanceChanged", "acc", acc)
	var (
		bu             = newBunchUpdater(k, ctx)
		node           string
		changeActivity bool
		active         bool
	)
	if err := bu.update(acc, true, func(value *types.Info) error {
		if value.IsEmpty() {
			return types.ErrNotFound
		}

		newDelegated := value.SelfDelegated.Add(dd)

		// TODO: move active account threshold to params
		if newDelegated.Int64() >= 100_000_000 {
			if !value.Active {
				changeActivity = true
				active = true
			}
		} else {
			if value.Active {
				changeActivity = true
				active = false
			}
		}

		if !dd.IsZero() {
			bu.addCallback(StakeChangedCallback, acc)
		}

		node = value.Referrer

		value.SelfDelegated = &newDelegated
		return nil
	}); err != nil {
		if errors.Is(err, types.ErrNotFound) {
			k.Logger(ctx).Debug("account is out of the referral", "acc", acc)
			return nil
		} else {
			k.Logger(ctx).Error("OnBalanceChanged hook failed", "acc", acc, "step", 0, "error", err)
			return err
		}
	}

	for i := 1; i <= 14; i++ {
		if node == "" {
			break
		}

		if err := bu.update(node, true, func(value *types.Info) error {
			newLevelDelegated := (*value.TeamDelegated).Add(dd)
			value.TeamDelegated = &newLevelDelegated
			if !dd.IsZero() {
				bu.addCallback(StakeChangedCallback, node)
			}

			node = value.Referrer
			return nil
		}); err != nil {
			k.Logger(ctx).Error("OnBalanceChanged hook failed", "acc", acc, "step", i, "error", err)
			return err
		}
	}

	if changeActivity {
		err := k.SetActive(ctx, acc, active, true, bu)
		if err != nil {
			k.Logger(ctx).Error("OnBalanceChanged hook failed", "acc", acc, "step", "set active status true", "error", err)
			return err
		}
	}

	if err := bu.commit(); err != nil {
		k.Logger(ctx).Error("OnBalanceChanged hook failed", "acc", acc, "step", "commit", "error", err)
		return err
	}

	return nil
}

func ChangeTeamActive(aag *types.ActiveAggregations, teamSize uint64, delta int64) {
	changeTeamActive(aag, teamSize, delta)
}

func changeTeamActive(aag *types.ActiveAggregations, teamSize uint64, delta int64) {
	if teamSize < 15 {
		aag.Team0 += int32(delta)
	} else if teamSize < 50 {
		aag.Team15 += int32(delta)
	} else if teamSize < 100 {
		aag.Team50 += int32(delta)
	} else if teamSize < 300 {
		aag.Team100 += int32(delta)
	} else {
		aag.Team300 += int32(delta)
	}
}

func (k Keeper) SetActive(ctx sdk.Context, acc string, value, checkAncestorsForStatusUpdate bool, bu *bunchUpdater) error {
	var (
		valueIsAlreadySet = false

		parent     string
		delta      func(*uint64)
		deltaValue int64
	)
	k.Logger(ctx).Debug("Set active", "acc", acc, "value", value, "checkAncestorsForStatusUpdate", checkAncestorsForStatusUpdate)
	if value {
		delta = func(x *uint64) { *x += 1 }
		deltaValue = 1
	} else {
		delta = func(x *uint64) { *x -= 1 }
		deltaValue = -1
	}

	err := bu.update(acc, false, func(x *types.Info) error {
		if x.Active == value {
			valueIsAlreadySet = true
		} else {
			x.Active = value
			delta(&x.ActiveRefCounts[0])
			parent = x.Referrer

			if parent != "" {
				err2 := bu.update(parent, checkAncestorsForStatusUpdate, func(y *types.Info) error {
					// if now account active - increment ActiveAgregations
					activeTeam := x.GetActiveRefsCountFromLevelToLevel(1, 13)
					// core criteria
					changeTeamActive(y.ActiveCount, activeTeam, deltaValue)
					// xby0 criteria
					y.ActiveCount.FirstLine += int32(deltaValue)
					// xby3 criteria for current account (y)
					if x.ActiveCount.FirstLine >= StatusGuruMinXParameter {
						y.ActiveCount.FirstLineBy3 += int32(deltaValue)
					}
					// xby3 criteria for current account parent (y.Referrer)
					if y.Active {
						xby3ParentDelta := 0
						// we go through top of xby3 requirement
						if x.Active && y.ActiveCount.FirstLine == StatusGuruMinXParameter {
							xby3ParentDelta = 1
						}
						// we go through bottom of xby3 requirement
						if !x.Active && y.ActiveCount.FirstLine == StatusGuruMinXParameter-1 {
							xby3ParentDelta = -1
						}
						// if account is active, and it's border change - update parent data
						if xby3ParentDelta != 0 && y.Referrer != "" {
							return bu.update(y.Referrer, checkAncestorsForStatusUpdate, func(z *types.Info) error {
								z.ActiveCount.FirstLineBy3 += int32(xby3ParentDelta)

								return nil
							})
						}
					}

					return nil
				})
				if err2 != nil {
					return err2
				}
			}
		}
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "cannot update acc info")
	} else if valueIsAlreadySet {
		return nil
	}

	// update 1 to 14 level referrers - update n-teams rule
	for i := 0; i < 14; i++ {
		if parent == "" {
			break
		}

		err = bu.update(parent, checkAncestorsForStatusUpdate, func(x *types.Info) error {
			oldTeamSize := x.GetActiveRefsCountFromLevelToLevel(1, 13)
			delta(&x.ActiveRefCounts[i+1])
			newTeamSize := x.GetActiveRefsCountFromLevelToLevel(1, 13)
			parent = x.Referrer

			// if account is active and team size changed - update parent status aggregations
			if x.Active && parent != "" && oldTeamSize != newTeamSize {
				err2 := bu.update(parent, checkAncestorsForStatusUpdate, func(y *types.Info) error {

					if deltaValue > 0 {
						changeTeamActive(y.ActiveCount, oldTeamSize, -deltaValue)
						changeTeamActive(y.ActiveCount, newTeamSize, deltaValue)
					}

					if deltaValue < 0 {
						changeTeamActive(y.ActiveCount, oldTeamSize, deltaValue)
						changeTeamActive(y.ActiveCount, newTeamSize, -deltaValue)
					}

					return nil
				})
				if err2 != nil {
					return err2
				}
			}

			return nil
		})
		if err != nil {
			return errors.Wrapf(err, "cannot update ancestor's referral count (#%d)", i)
		}
	}

	return nil
}

func (k Keeper) MustSetActive(ctx sdk.Context, acc string, value bool, bu *bunchUpdater) {
	if err := k.SetActive(ctx, acc, value, true, bu); err != nil {
		panic(err)
	}
}

// MustSetActiveWithoutStatusUpdate updates active referrals but skips status update check after it. So this check MUST
// be performed from the outer code later. This is useful for massive updates like genesis init, because it allows to
// avoid excessive checks repeating again and again for the same account (every time any of referrals up to 14 lines
// down changes its activity).
func (k Keeper) MustSetActiveWithoutStatusUpdate(ctx sdk.Context, acc string, value bool, bu *bunchUpdater) {
	if err := k.SetActive(ctx, acc, value, false, bu); err != nil {
		panic(err)
	}
}

func (k Keeper) getReferralFeesCore(ctx sdk.Context, acc string, toAncestors []util.Fraction) ([]types.ReferralFee, util.Fraction, error) {
	if len(toAncestors) != 14 {
		return nil, util.Fraction{}, errors.Errorf("toAncestors param must have exactly 14 items (%d found)", len(toAncestors))
	}
	excess := util.Percent(0)
	result := make([]types.ReferralFee, 0, 14)

	ancestor, err := k.GetParent(ctx, acc)
	k.Logger(ctx).Info("Get starting at", "anc", ancestor)
	if err != nil {
		return nil, util.Fraction{}, err
	}
	for i := 0; i < 14; i++ {
		var (
			data types.Info
			err  error
		)

		if ancestor == "" {
			excess = excess.Add(toAncestors[i])
			continue
		}

		data, err = k.Get(ctx, ancestor)
		if err != nil {
			return nil, util.Fraction{}, err
		}

		if i < data.Status.LinesOpened() {
			if !toAncestors[i].IsZero() {
				result = append(result, types.ReferralFee{Beneficiary: ancestor, Ratio: toAncestors[i]})
			}
		} else {
			excess = excess.Add(toAncestors[i])
		}

		ancestor = data.Referrer
	}

	return result, excess, nil
}

// TODO: should we need to remove this function?
func (k Keeper) setStatus(ctx sdk.Context, target *types.Info, value types.Status, acc string) {
	if target.Status == value {
		return
	}
	target.Status = value
}

func (k Keeper) PayUpFees(ctx sdk.Context, acc string, totalAmount math.Int) (remain math.Int, err error) {
	fees, burn, err := k.GetReferralFeesForDelegating(ctx, acc)
	if err != nil {
		return totalAmount, err
	}

	cdc := k.accountKeeper.AddressCodec()
	accAddr, err := cdc.StringToBytes(acc)
	if err != nil {
		return totalAmount, err
	}

	bondDenom, err := k.stakingKeeper.BondDenom(ctx)
	if err != nil {
		return totalAmount, err
	}

	amountToBurn := math.NewInt(burn.MulInt64(totalAmount.Int64()).Int64())
	coinsToBurn := sdk.NewCoins(sdk.NewCoin(bondDenom, amountToBurn))

	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, accAddr, k.referralAccountName, coinsToBurn)
	if err != nil {
		return totalAmount, err
	}

	err = k.bankKeeper.BurnCoins(ctx, k.referralAccountName, coinsToBurn)
	if err != nil {
		return totalAmount, err
	}

	totalFee := int64(0)
	outputs := make([]banktypes.Output, 0, len(fees))
	for _, fee := range fees {
		x := fee.Ratio.MulInt64(totalAmount.Int64()).Int64()
		if x == 0 {
			continue
		}
		totalFee += x
		amount := sdk.NewCoins(sdk.NewCoin(bondDenom, math.NewInt(x)))
		outputs = append(outputs, banktypes.NewOutput(fee.GetBeneficiary(), amount))

		to, err := cdc.BytesToString(fee.GetBeneficiary())
		if err != nil {
			return totalAmount, err
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRefFee,
				sdk.NewAttribute(sdk.AttributeKeyAmount, amount.String()),
				sdk.NewAttribute(types.AttributeKeyFrom, acc),
				sdk.NewAttribute(types.AttributeKeyTo, to),
			),
		)
	}

	if totalFee != 0 {
		input := banktypes.NewInput(accAddr, sdk.NewCoins(sdk.NewCoin(util.ConfigMainDenom, math.NewInt(totalFee))))

		err = k.bankKeeper.InputOutputCoins(ctx, input, outputs)
		if err != nil {
			return totalAmount, err
		}
	}

	remain = totalAmount.Sub(amountToBurn).SubRaw(totalFee)

	return remain, nil
}

func (k Keeper) BurnCoins(ctx sdk.Context, acc sdk.AccAddress, amt sdk.Coins) error {
	err := k.bankKeeper.SendCoinsFromAccountToModule(ctx, acc, k.referralAccountName, amt)
	if err != nil {
		return err
	}

	return k.bankKeeper.BurnCoins(ctx, k.referralAccountName, amt)
}

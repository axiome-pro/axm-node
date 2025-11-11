package keeper

import (
	"context"
	"errors"
	"fmt"

	abci "github.com/cometbft/cometbft/abci/types"

	addresscodec "cosmossdk.io/core/address"
	storetypes "cosmossdk.io/core/store"
	"cosmossdk.io/log"
	"cosmossdk.io/math"

	"github.com/axiome-pro/axm-node/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Implements ValidatorSet interface
var _ types.ValidatorSet = Keeper{}

// Implements DelegationSet interface
var _ types.DelegationSet = Keeper{}

// Keeper of the x/staking store
type Keeper struct {
	storeService          storetypes.KVStoreService
	cdc                   codec.BinaryCodec
	authKeeper            types.AccountKeeper
	bankKeeper            types.BankKeeper
	hooks                 types.StakingHooks
	refHooks              types.RefStakingHooks
	authority             string
	validatorAddressCodec addresscodec.Codec
	consensusAddressCodec addresscodec.Codec
}

// NewKeeper creates a new staking Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	storeService storetypes.KVStoreService,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	authority string,
	validatorAddressCodec addresscodec.Codec,
	consensusAddressCodec addresscodec.Codec,
) *Keeper {
	// ensure bonded and not bonded module accounts are set
	if addr := ak.GetModuleAddress(types.BondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.BondedPoolName))
	}

	if addr := ak.GetModuleAddress(types.NotBondedPoolName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.NotBondedPoolName))
	}

	// ensure that authority is a valid AccAddress
	if _, err := ak.AddressCodec().StringToBytes(authority); err != nil {
		panic("authority is not a valid acc address")
	}

	if validatorAddressCodec == nil || consensusAddressCodec == nil {
		panic("validator and/or consensus address codec are nil")
	}

	return &Keeper{
		storeService:          storeService,
		cdc:                   cdc,
		authKeeper:            ak,
		bankKeeper:            bk,
		hooks:                 nil,
		authority:             authority,
		validatorAddressCodec: validatorAddressCodec,
		consensusAddressCodec: consensusAddressCodec,
	}
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx context.Context) log.Logger {
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	return sdkCtx.Logger().With("module", "x/"+types.ModuleName)
}

// Hooks gets the hooks for staking *Keeper {
func (k *Keeper) Hooks() types.StakingHooks {
	if k.hooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiStakingHooks{}
	}

	return k.hooks
}

// SetHooks sets the validator hooks.  In contrast to other receivers, this method must take a pointer due to nature
// of the hooks interface and SDK start up sequence.
func (k *Keeper) SetHooks(sh types.StakingHooks) {
	if k.hooks != nil {
		panic("cannot set validator hooks twice")
	}

	k.hooks = sh
}

// Hooks gets the hooks for staking *Keeper {
func (k *Keeper) RefHooks() types.RefStakingHooks {
	if k.refHooks == nil {
		// return a no-op implementation if no hooks are set
		return types.MultiRefStakingHooks{}
	}

	return k.refHooks
}

// SetRefHooks sets the validator hooks.  In contrast to other receivers, this method must take a pointer due to nature
// of the hooks interface and SDK start up sequence.
func (k *Keeper) SetRefHooks(sh types.RefStakingHooks) {
	if k.refHooks != nil {
		panic("cannot set validator hooks twice")
	}

	k.refHooks = sh
}

// GetLastTotalPower loads the last total validator power.
func (k Keeper) GetLastTotalPower(ctx context.Context) (math.Int, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.LastTotalPowerKey)
	if err != nil {
		return math.ZeroInt(), err
	}

	if bz == nil {
		return math.ZeroInt(), nil
	}

	ip := sdk.IntProto{}
	err = k.cdc.Unmarshal(bz, &ip)
	if err != nil {
		return math.ZeroInt(), err
	}

	return ip.Int, nil
}

// SetLastTotalPower sets the last total validator power.
func (k Keeper) SetLastTotalPower(ctx context.Context, power math.Int) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&sdk.IntProto{Int: power})
	if err != nil {
		return err
	}
	return store.Set(types.LastTotalPowerKey, bz)
}

// GetAuthority returns the x/staking module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// ValidatorAddressCodec returns the app validator address codec.
func (k Keeper) ValidatorAddressCodec() addresscodec.Codec {
	return k.validatorAddressCodec
}

// ConsensusAddressCodec returns the app consensus address codec.
func (k Keeper) ConsensusAddressCodec() addresscodec.Codec {
	return k.consensusAddressCodec
}

// SetValidatorUpdates sets the ABCI validator power updates for the current block.
func (k Keeper) SetValidatorUpdates(ctx context.Context, valUpdates []abci.ValidatorUpdate) error {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := k.cdc.Marshal(&types.ValidatorUpdates{Updates: valUpdates})
	if err != nil {
		return err
	}
	return store.Set(types.ValidatorUpdatesKey, bz)
}

// GetValidatorUpdates returns the ABCI validator power updates within the current block.
func (k Keeper) GetValidatorUpdates(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.ValidatorUpdatesKey)
	if err != nil {
		return nil, err
	}

	var valUpdates types.ValidatorUpdates
	err = k.cdc.Unmarshal(bz, &valUpdates)
	if err != nil {
		return nil, err
	}

	return valUpdates.Updates, nil
}

// SetStakeMoveVoting marks a delegator-validator pair as having an ongoing vote
func (k Keeper) SetStakeMoveVoting(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetStakeMoveVotingKey(delAddr, valAddr)
	// store presence with a single byte
	return store.Set(key, []byte{1})
}

// DeleteStakeMoveVoting removes the ongoing vote mark for a delegator-validator pair
func (k Keeper) DeleteStakeMoveVoting(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) error {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetStakeMoveVotingKey(delAddr, valAddr)
	return store.Delete(key)
}

// IsStakeMoveVoting checks if there is an ongoing vote for the delegator-validator pair
func (k Keeper) IsStakeMoveVoting(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (bool, error) {
	store := k.storeService.OpenKVStore(ctx)
	key := types.GetStakeMoveVotingKey(delAddr, valAddr)
	bz, err := store.Get(key)
	if err != nil {
		return false, err
	}
	return bz != nil, nil
}

// MoveDelegation transfers entire delegation shares from src delegator to dst delegator for a given validator.
// Steps:
// 1) Withdraw rewards for src delegator-validator; if dst already delegates to validator, withdraw its rewards too.
// 2) Call referral hooks to decrease src amount and increase dst amount.
// 3) Remove src delegation record.
// 4) Create or update dst delegation with moved shares and update points.
// 5) Invoke staking hooks Before/After delegation changes as appropriate.
func (k Keeper) MoveDelegation(ctx context.Context, srcDelegator, dstDelegator, validatorOper string) error {
	// convert addresses
	srcAcc, err := k.authKeeper.AddressCodec().StringToBytes(srcDelegator)
	if err != nil {
		return err
	}
	dstAcc, err := k.authKeeper.AddressCodec().StringToBytes(dstDelegator)
	if err != nil {
		return err
	}
	valAddr, err := k.validatorAddressCodec.StringToBytes(validatorOper)
	if err != nil {
		return err
	}

	// load validator and delegations
	val, err := k.GetValidator(ctx, valAddr)
	if err != nil {
		return err
	}

	srcDel, err := k.GetDelegation(ctx, srcAcc, valAddr)
	if err != nil {
		return err
	}

	// calculate source coins amount
	oldSrcCoins := val.TokensFromSharesTruncated(srcDel.GetShares()).TruncateInt()

	// check dst delegation existence and current coins
	dstDel, err := k.GetDelegation(ctx, dstAcc, valAddr)
	dstExists := (err == nil)
	if err != nil && !errors.Is(err, types.ErrNoDelegation) {
		return err
	}
	oldDstCoins := math.ZeroInt()
	if dstExists {
		oldDstCoins = val.TokensFromSharesTruncated(dstDel.GetShares()).TruncateInt()
	}

	// 1) withdraw rewards for src, and for dst if exists
	if err := k.Hooks().BeforeDelegationSharesModified(ctx, srcAcc, valAddr); err != nil {
		return err
	}
	if dstExists {
		if err := k.Hooks().BeforeDelegationSharesModified(ctx, dstAcc, valAddr); err != nil {
			return err
		}
	} else {
		// signal creation for distribution periods when creating new delegation
		if err := k.Hooks().BeforeDelegationCreated(ctx, dstAcc, valAddr); err != nil {
			return err
		}
	}

	// 2) referral hooks updates
	if err := k.RefHooks().DelegationCoinsModified(ctx, srcDelegator, validatorOper, oldSrcCoins, math.ZeroInt()); err != nil {
		return err
	}
	newDstCoins := oldDstCoins.Add(oldSrcCoins)
	if err := k.RefHooks().DelegationCoinsModified(ctx, dstDelegator, validatorOper, oldDstCoins, newDstCoins); err != nil {
		return err
	}

	// 3) remove src delegation (record only)
	if err := k.RemoveDelegation(ctx, srcDel); err != nil {
		return err
	}

	// 4) update/create dst delegation with moved shares and points
	if !dstExists {
		// create new delegation
		newDel := types.NewDelegation(dstDelegator, validatorOper, srcDel.GetShares())
		newDel.Points = val.Points
		if err := k.SetDelegation(ctx, newDel); err != nil {
			return err
		}
	} else {
		// update existing
		dstDel.Shares = dstDel.Shares.Add(srcDel.GetShares())
		dstDel.Points = val.Points
		if err := k.SetDelegation(ctx, dstDel); err != nil {
			return err
		}
	}

	// 5) staking hooks after modification for dst
	if err := k.Hooks().AfterDelegationModified(ctx, dstAcc, valAddr); err != nil {
		return err
	}

	return nil
}

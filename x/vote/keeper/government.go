package keeper

import (
	"github.com/axiome-pro/axm-node/x/vote/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/golang/protobuf/proto"
)

func (k Keeper) GetGovernment(ctx sdk.Context) types.Government {
	store := k.storeService.OpenKVStore(ctx)

	bz, _ := store.Get(types.KeyGovernment)

	var gov types.Government

	err := proto.Unmarshal(bz, &gov)
	if err != nil {
		panic(err)
	}

	return gov
}

func (k Keeper) SetGovernment(ctx sdk.Context, gov types.Government) {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := proto.Marshal(&gov)
	if err != nil {
		panic(err)
	}
	err = store.Set(types.KeyGovernment, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) RemoveGovernor(ctx sdk.Context, gov sdk.AccAddress) error {
	govs := k.GetGovernment(ctx)
	if !govs.Contains(gov) {
		return types.ErrProposalGovernorNotExists
	}

	if len(govs.GetMembers()) == 1 {
		return types.ErrProposalGovernorLast
	}

	govs.Remove(gov)
	k.SetGovernment(ctx, govs)
	return nil
}

func (k Keeper) AddGovernor(ctx sdk.Context, gov sdk.AccAddress) error {
	govs := k.GetGovernment(ctx)
	if govs.Contains(gov) {
		return types.ErrProposalGovernorExists
	}
	govs.Append(gov)
	k.SetGovernment(ctx, govs)
	return nil
}

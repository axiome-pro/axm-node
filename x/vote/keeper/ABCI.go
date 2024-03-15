package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlock(ctx sdk.Context) error {
	defer func() {
		if r := recover(); r != nil {
			k.Logger(ctx).Info("vote recovered from panic on begin blocker")
		}
	}()

	proposal := k.GetCurrentProposal(ctx)
	if proposal != nil {
		if proposal.EndTime.Before(ctx.BlockTime()) {
			_, agree := k.Validate(
				k.GetGovernment(ctx),
				k.GetAgreed(ctx),
				k.GetDisagreed(ctx),
			)

			k.EndProposal(ctx, *proposal, agree)
		}
	}

	poll, ok := k.GetCurrentPoll(ctx)
	if ok {
		if poll.EndTime.Before(ctx.BlockTime()) {
			k.EndPoll(ctx)
		}
	}

	return nil
}

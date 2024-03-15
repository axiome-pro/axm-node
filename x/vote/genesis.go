package vote

import (
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axiome-pro/axm-node/x/vote/keeper"
	"github.com/axiome-pro/axm-node/x/vote/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	k.Logger(ctx).Info("Starting from genesis...")
	k.SetParams(ctx, data.Params)
	k.SetGovernment(ctx, types.Government{Members: data.Government})
	if len(data.CurrentProposal.Messages) > 0 {
		k.SetCurrentProposal(ctx, data.CurrentProposal)
		k.SetStartBlock(ctx.WithBlockHeight(data.StartBlock))
	}
	if len(data.Agreed) > 0 {
		k.SetAgreed(ctx, types.Government{Members: data.Agreed})
	}
	if len(data.Disagreed) > 0 {
		k.SetDisagreed(ctx, types.Government{Members: data.Disagreed})
	}
	for _, record := range data.History {
		k.AddProposalHistoryRecord(ctx, record)
	}
	k.LoadPolls(ctx, data)
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data *types.GenesisState) {
	var currentProposal *types.Proposal
	var startBlock int64
	if currentProposal = k.GetCurrentProposal(ctx); currentProposal != nil {
		startBlock = k.GetStartBlock(ctx)
	}
	data = types.NewGenesisState(
		k.GetParams(ctx),
		k.GetGovernment(ctx),
		currentProposal,
		startBlock,
		k.GetAgreed(ctx),
		k.GetDisagreed(ctx),
		k.GetHistory(ctx, math.MaxInt32, 1),
	)
	if poll, ok := k.GetCurrentPoll(ctx); ok {
		data.CurrentPoll = &poll
		if err := k.IterateThroughCurrentPollAnswers(ctx, func(acc string, ans bool) (stop bool) {
			data.PollAnswers = append(data.PollAnswers, types.PollAnswer{Acc: acc, Ans: ans})
			return false
		}); err != nil {
			panic(err)
		}
	}
	data.PollHistory = k.GetPollHistoryAll(ctx)
	return data
}

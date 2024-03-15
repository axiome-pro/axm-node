package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) OneDay(ctx sdk.Context) time.Duration {
	//TODO: move to genesis to create fast forwards testnets
	return time.Duration(86_400_000_000_000)
}
func (k Keeper) OneWeek(ctx sdk.Context) time.Duration  { return 7 * k.OneDay(ctx) }
func (k Keeper) OneMonth(ctx sdk.Context) time.Duration { return 30 * k.OneDay(ctx) }

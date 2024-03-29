package keeper

import (
	"context"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"

	"github.com/axiome-pro/axm-node/x/staking/types"
	"github.com/cosmos/cosmos-sdk/telemetry"
)

// BeginBlocker will persist the current header and validator set as a historical entry
// and prune the oldest entry based on the HistoricalEntries parameter
func (k *Keeper) BeginBlocker(ctx context.Context) error {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)
	err := k.TrackHistoricalInfo(ctx)

	if err != nil {
		return err
	}

	return k.AllocateValidatorsPoints(ctx)
}

// EndBlocker called at every block, update validator set
func (k *Keeper) EndBlocker(ctx context.Context) ([]abci.ValidatorUpdate, error) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyEndBlocker)
	return k.BlockValidatorUpdates(ctx)
}

package vote

import (
	"github.com/axiome-pro/axm-node/x/vote/keeper"
	"github.com/axiome-pro/axm-node/x/vote/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	// functions aliases
	DefaultGenesisState = types.DefaultGenesisState
	ValidateGenesis     = types.ValidateGenesis
)

type (
	Keeper       = keeper.Keeper
	GenesisState = types.GenesisState
	Params       = types.Params
)

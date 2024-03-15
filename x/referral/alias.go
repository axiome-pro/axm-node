package referral

import (
	"github.com/axiome-pro/axm-node/x/referral/keeper"
	"github.com/axiome-pro/axm-node/x/referral/types"
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
	Keeper            = keeper.Keeper
	GenesisState      = types.GenesisState
	ReferralFee       = types.ReferralFee
	Params            = types.Params
	NetworkAward      = types.NetworkAward
	StatusCheckResult = types.StatusCheckResult
	Status            = types.Status
	DataRecord        = types.Info
)

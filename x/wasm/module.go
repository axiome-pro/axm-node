package wasm

import (
	modulev1 "github.com/axiome-pro/axm-node/api/axiome/wasm/module/v1"
	distributionkeeper "github.com/axiome-pro/axm-node/x/distribution/keeper"
	stakingKeeper "github.com/axiome-pro/axm-node/x/staking/keeper"

	//distributionkeeper "github.com/axiome-pro/axm-node/x/distribution/keeper"
	exportedtypes "github.com/axiome-pro/axm-node/x/wasm/types"
	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"github.com/CosmWasm/wasmd/x/wasm"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/CosmWasm/wasmd/x/wasm/simulation"
)

func init() {
	appmodule.Register(&modulev1.Module{}, appmodule.Provide(
		ProvideWasmKeeper,
		ProvideWasmModule,
	))
}

var (
	_ appmodule.AppModule         = wasm.AppModule{}
	_ module.AppModuleBasic       = wasm.AppModuleBasic{}
	_ exportedtypes.StakingKeeper = stakingKeeper.Keeper{}
	//_ exportedtypes.DistributionKeeper =  distributionkeeper.Keeper{}
	_ types.StakingKeeper = exportedtypes.WrappedStakingKeeper{}
)

type ModuleInputs struct {
	depinject.In

	Config       *modulev1.Module
	StoreService store.KVStoreService
	Cdc          codec.Codec
	//MsgServiceRouter baseapp.MsgServiceRouter
	MessageRouter baseapp.MessageRouter

	AccountKeeper types.AccountKeeper
	BankKeeper    types.BankKeeper

	Options servertypes.AppOptions `optional:"true"`

	DistrKeeper   distributionkeeper.Keeper
	StakingKeeper exportedtypes.StakingKeeper
	//SubSpace paramtypes.Subspace

	ValidatorSetSource wasmkeeper.ValidatorSetSource
	SimBankKeeper      simulation.BankKeeper
}

type ModuleOutputs struct {
	depinject.Out

	WasmKeper wasmkeeper.Keeper
	//Module    appmodule.AppModule
}

// ProvideWasmModule — строит модуль
func ProvideWasmKeeper(in ModuleInputs) ModuleOutputs {
	var homeDir = "."

	authority := govtypes.ModuleName
	if in.Config.Authority != "" {
		authority = in.Config.Authority
	}

	if in.Options != nil {
		if v := in.Options.Get(flags.FlagHome); v != nil {
			homeDir, _ = v.(string)
		}

		// Pass nil for StakingKeeper since we can't create a proper WrappedStakingKeeper instance
		// The type mismatch between our StakingKeeper and the one expected by wasmd prevents us from using it directly
		k := wasmkeeper.NewKeeper(
			in.Cdc,
			in.StoreService,
			in.AccountKeeper,
			in.BankKeeper,
			exportedtypes.WrappedStakingKeeper{Keeper: in.StakingKeeper},
			exportedtypes.WrappedDistributionKeeper{Keeper: distributionkeeper.NewQuerier(in.DistrKeeper)},
			nil,
			nil,
			nil,
			nil,
			nil,
			in.MessageRouter,
			nil,
			homeDir,
			types.DefaultWasmConfig(),
			"iterator,staking,distribution,cosmwasm_1_1,cosmwasm_1_2,cosmwasm_1_3,cosmwasm_1_4",
			authority,
		)

		return ModuleOutputs{
			WasmKeper: k,
		}
	}
	return ModuleOutputs{}
}

func ProvideWasmModule(
	keeper wasmkeeper.Keeper,
	cdc codec.Codec,
	validatorSetSource wasmkeeper.ValidatorSetSource,
	accountKeeper types.AccountKeeper,
	simBankKeeper simulation.BankKeeper,
) appmodule.AppModule {
	return wasm.NewAppModule(
		cdc,
		&keeper,
		validatorSetSource,
		accountKeeper,
		simBankKeeper,
		nil,
		nil,
	)
}

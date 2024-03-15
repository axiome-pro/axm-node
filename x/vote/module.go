package vote

import (
	"context"
	"cosmossdk.io/core/store"
	"cosmossdk.io/depinject"
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/baseapp"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	modulev1 "github.com/axiome-pro/axm-node/api/axiome/vote/module/v1"
	"github.com/axiome-pro/axm-node/x/vote/client/cli"
	"github.com/axiome-pro/axm-node/x/vote/keeper"
	"github.com/axiome-pro/axm-node/x/vote/types"
	"cosmossdk.io/core/address"
	"cosmossdk.io/core/appmodule"
)

// TypeCode check to ensure the interface is properly implemented
var (
	_ appmodule.AppModule       = AppModule{}
	_ appmodule.HasBeginBlocker = AppModule{}
	_ module.HasGenesis         = AppModule{}
	_ module.HasServices        = AppModule{}

	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the vote module.
type AppModuleBasic struct {
	cdc codec.Codec
	ac  address.Codec
}

// Name returns the vote module's name.
func (AppModuleBasic) Name() string {
	return types.ModuleName
}

// RegisterCodec registers the vote module's types for the given codec.
func (AppModuleBasic) RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	types.RegisterLegacyAminoCodec(cdc)
}

func (AppModuleBasic) RegisterInterfaces(registry codecTypes.InterfaceRegistry) {
	types.RegisterInterfaces(registry)
}

// DefaultGenesis returns default genesis state as raw bytes for the vote
// module.
func (AppModuleBasic) DefaultGenesis(cdc codec.JSONCodec) json.RawMessage {
	return cdc.MustMarshalJSON(types.DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the vote module.
func (AppModuleBasic) ValidateGenesis(cdc codec.JSONCodec, _ client.TxEncodingConfig, bz json.RawMessage) error {
	var data types.GenesisState
	err := cdc.UnmarshalJSON(bz, &data)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal x/vote genesis state")
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers the REST routes for the vote module.
func (AppModuleBasic) RegisterRESTRoutes(ctx client.Context, rtr *mux.Router) {}

// GetTxCmd returns the root tx command for the vote module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return cli.NewTxCmd()
}

// GetQueryCmd returns no root query command for the vote module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
	//return cli.NewQueryCmd()
}

func (AppModuleBasic) RegisterGRPCGatewayRoutes(clientCtx client.Context, mux *runtime.ServeMux) {
	err := types.RegisterQueryHandlerClient(context.Background(), mux, types.NewQueryClient(clientCtx))
	if err != nil {
		panic(err)
	}
}

//____________________________________________________________________________

// AppModule implements an application module for the vote module.
type AppModule struct {
	AppModuleBasic

	keeper         keeper.Keeper
	referralKeeper types.ReferralKeeper
}

func (am AppModule) BeginBlock(ctx context.Context) error {
	return am.keeper.BeginBlock(sdk.UnwrapSDKContext(ctx))
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, k keeper.Keeper,
	referralKeeper types.ReferralKeeper,
) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{cdc: cdc},
		keeper:         k,
		referralKeeper: referralKeeper,
	}
}

// Name returns the vote module's name.
func (AppModule) Name() string { return types.ModuleName }

// IsOnePerModuleType implements the depinject.OnePerModuleType interface.
func (am AppModule) IsOnePerModuleType() {}

// IsAppModule implements the appmodule.AppModule interface.
func (am AppModule) IsAppModule() {}

func (am AppModule) RegisterServices(cfg module.Configurator) {
	types.RegisterMsgServer(cfg.MsgServer(), keeper.MsgServer(am.keeper))
	types.RegisterQueryServer(cfg.QueryServer(), keeper.QueryServer{Keeper: am.keeper})
}

// RegisterInvariants registers the vote module invariants.
func (am AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// InitGenesis performs genesis initialization for the vote module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, mrshl codec.JSONCodec, data json.RawMessage) {
	var genesisState types.GenesisState
	mrshl.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, am.keeper, genesisState)
}

// ExportGenesis returns the exported genesis state as raw bytes for the vote
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context, mrshl codec.JSONCodec) json.RawMessage {
	gs := ExportGenesis(ctx, am.keeper)
	return mrshl.MustMarshalJSON(gs)
}

//// BeginBlock returns the begin blocker for the vote module.
//func (am AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
//}
//
//// EndBlock returns the end blocker for the vote module. It returns no validator
//// updates.
//func (AppModule) EndBlock(_ sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
//	return []abci.ValidatorUpdate{}
//}

//
// App Wiring Setup
//

func init() {
	appmodule.Register(&modulev1.Module{},
		appmodule.Provide(ProvideModule),
	)
}

type ModuleInputs struct {
	depinject.In

	Config         *modulev1.Module
	StoreService   store.KVStoreService
	Cdc            codec.Codec
	ReferralKeeper types.ReferralKeeper
	AccountKeeper  types.AccountKeeper

	MsgServiceRouter baseapp.MessageRouter
}

type ModuleOutputs struct {
	depinject.Out

	VoteKeeper keeper.Keeper
	Module     appmodule.AppModule
}

func ProvideModule(in ModuleInputs) ModuleOutputs {
	// default to governance authority if not provided
	authority := authtypes.NewModuleAddress(types.ModuleName)
	if in.Config.Authority != "" {
		authority = authtypes.NewModuleAddressOrBech32Address(in.Config.Authority)
	}

	k := keeper.NewKeeper(
		in.Cdc,
		in.StoreService,
		in.ReferralKeeper,
		authority,
		in.MsgServiceRouter,
		in.AccountKeeper,
	)

	m := NewAppModule(in.Cdc, k, in.ReferralKeeper)

	return ModuleOutputs{
		VoteKeeper: k,
		Module:     m,
	}
}

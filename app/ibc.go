package app

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	capability "github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
)

// ibcStoreKeys are the store keys added by IBC. They are listed here so that the
// upgrade handler and the wiring below cannot drift apart.
var ibcStoreKeys = []string{
	capabilitytypes.StoreKey,
	ibcexported.StoreKey,
	ibctransfertypes.StoreKey,
}

// setupIBC wires the IBC core module, the capability module and ICS-20 transfer
// into the application.
//
// These modules cannot be declared in app.yaml: ibc-go ships no depinject module
// configs in any released version, so they are constructed by hand and attached
// with the runtime hooks the SDK provides for exactly this case,
// runtime.App.RegisterStores and runtime.App.RegisterModules.
//
// It must be called after appBuilder.Build and before app.Load, because
// RegisterStores mounts the new stores and Load is what reads them.
func (app *AxmApp) setupIBC(appCodec codec.Codec, authority string) error {
	capabilityKey := storetypes.NewKVStoreKey(capabilitytypes.StoreKey)
	ibcKey := storetypes.NewKVStoreKey(ibcexported.StoreKey)
	transferKey := storetypes.NewKVStoreKey(ibctransfertypes.StoreKey)
	capabilityMemKey := storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey)

	if err := app.RegisterStores(capabilityKey, ibcKey, transferKey, capabilityMemKey); err != nil {
		return err
	}

	capabilityKeeper := capabilitykeeper.NewKeeper(appCodec, capabilityKey, capabilityMemKey)
	scopedIBCKeeper := capabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := capabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	capabilityKeeper.Seal()
	app.CapabilityKeeper = capabilityKeeper

	// The legacy param subspace is nil on purpose. It exists only so that chains
	// migrating from ibc-go v7 can read their old params; Axiome starts directly
	// on self-managed params, and the chain has no x/params module to hand over.
	app.IBCKeeper = ibckeeper.NewKeeper(
		appCodec,
		ibcKey,
		nil,
		newIBCStakingKeeper(app.StakingKeeper),
		app.UpgradeKeeper,
		scopedIBCKeeper,
		authority,
	)

	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		appCodec,
		transferKey,
		nil,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
		authority,
	)

	transferModule := ibctransfer.NewIBCModule(app.TransferKeeper)
	ibcRouter := porttypes.NewRouter()
	ibcRouter.AddRoute(ibctransfertypes.ModuleName, transferModule)
	app.IBCKeeper.SetRouter(ibcRouter)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	if err := app.RegisterModules(
		capability.NewAppModule(appCodec, *app.CapabilityKeeper, true),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
	); err != nil {
		return err
	}

	// x/upgrade captures the initial module version map through a depinject
	// invoker, which only sees modules declared in app.yaml. Without this the
	// three modules above are absent from the stored version map, and the next
	// upgrade would treat them as newly added and run their InitGenesis again,
	// wiping live IBC state. Re-publishing the full map from the module manager
	// fixes that; the SDK explicitly allows for it in x/upgrade's InitGenesis.
	app.UpgradeKeeper.SetInitVersionMap(app.ModuleManager.GetVersionMap())

	return nil
}

// ibcModuleAccountPermissions returns the module account the ICS-20 transfer
// module needs. Its constructor panics if the account is missing, so this is
// asserted at startup rather than discovered on the first transfer.
func ibcModuleAccountPermissions() map[string][]string {
	return map[string][]string{
		ibctransfertypes.ModuleName: {authtypes.Minter, authtypes.Burner},
	}
}

// IBCModuleBasics returns the module basics for the manually wired IBC modules.
//
// The root command builds its BasicManager straight from the app config, which
// only knows about modules declared in app.yaml. Without this, `axmd init`
// produces a genesis with no IBC, capability or transfer state, and the client
// cannot decode IBC messages.
func IBCModuleBasics() []module.AppModuleBasic {
	return []module.AppModuleBasic{
		capability.AppModuleBasic{},
		ibc.AppModuleBasic{},
		ibctransfer.AppModuleBasic{},
	}
}

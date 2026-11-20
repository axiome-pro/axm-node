package app

import (
	"context"
	"time"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	sdkstakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/cosmos/ibc-go/modules/capability"
	capabilitykeeper "github.com/cosmos/ibc-go/modules/capability/keeper"
	capabilitytypes "github.com/cosmos/ibc-go/modules/capability/types"
	ibctransfer "github.com/cosmos/ibc-go/v8/modules/apps/transfer"
	ibctransferkeeper "github.com/cosmos/ibc-go/v8/modules/apps/transfer/keeper"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	ibc "github.com/cosmos/ibc-go/v8/modules/core"
	porttypes "github.com/cosmos/ibc-go/v8/modules/core/05-port/types"
	ibcexported "github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v8/modules/core/keeper"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"

	stakingkeeper "github.com/axiome-pro/axm-node/x/staking/keeper"
	votetypes "github.com/axiome-pro/axm-node/x/vote/types"
)

// IBC modules are not wired through depinject (ibc-go v8 does not support it),
// so they are registered by hand after the depinject-built app, the same way
// the SDK 0.50 / ignite templates do it.

// ibcStakingKeeper adapts the Axiome staking keeper to the interface ibc-go
// expects. IBC only reads the header of the historical info (time, app hash,
// next validators hash) to build its own consensus state, so the validator set
// is not converted.
type ibcStakingKeeper struct {
	k *stakingkeeper.Keeper
}

func (s ibcStakingKeeper) GetHistoricalInfo(ctx context.Context, height int64) (sdkstakingtypes.HistoricalInfo, error) {
	hi, err := s.k.GetHistoricalInfo(ctx, height)
	if err != nil {
		return sdkstakingtypes.HistoricalInfo{}, err
	}
	return sdkstakingtypes.HistoricalInfo{Header: hi.Header}, nil
}

func (s ibcStakingKeeper) UnbondingTime(ctx context.Context) (time.Duration, error) {
	return s.k.UnbondingTime(ctx)
}

// sdkStakingQuery serves the standard cosmos.staking.v1beta1.Query/Params on
// top of the Axiome staking module (which has its own axiome.staking package).
// IBC relayers (Hermes, rly) read the unbonding period from there to build
// light clients; wallets and explorers use it too. Only Params is served.
type sdkStakingQuery struct {
	sdkstakingtypes.UnimplementedQueryServer
	k *stakingkeeper.Keeper
}

func (q sdkStakingQuery) Params(ctx context.Context, _ *sdkstakingtypes.QueryParamsRequest) (*sdkstakingtypes.QueryParamsResponse, error) {
	p, err := q.k.GetParams(ctx)
	if err != nil {
		return nil, err
	}
	return &sdkstakingtypes.QueryParamsResponse{Params: sdkstakingtypes.Params{
		UnbondingTime:     p.UnbondingTime,
		MaxValidators:     p.MaxValidators,
		MaxEntries:        p.MaxEntries,
		HistoricalEntries: p.HistoricalEntries,
		BondDenom:         p.BondDenom,
		MinCommissionRate: math.LegacyZeroDec(),
	}}, nil
}

// Axiome never had x/params, so the legacy param subspaces passed to the IBC
// keepers are nil: they are only read by the v7->v8 param migrations, which
// never run on a chain that starts IBC at v8.

func (app *AxmApp) registerIBCModules() error {
	if err := app.RegisterStores(
		storetypes.NewKVStoreKey(capabilitytypes.StoreKey),
		storetypes.NewKVStoreKey(ibcexported.StoreKey),
		storetypes.NewKVStoreKey(ibctransfertypes.StoreKey),
		storetypes.NewMemoryStoreKey(capabilitytypes.MemStoreKey),
	); err != nil {
		return err
	}

	authority := authtypes.NewModuleAddress(votetypes.ModuleName).String()

	app.CapabilityKeeper = capabilitykeeper.NewKeeper(
		app.appCodec,
		app.GetKey(capabilitytypes.StoreKey),
		app.GetMemKey(capabilitytypes.MemStoreKey),
	)
	scopedIBCKeeper := app.CapabilityKeeper.ScopeToModule(ibcexported.ModuleName)
	scopedTransferKeeper := app.CapabilityKeeper.ScopeToModule(ibctransfertypes.ModuleName)
	app.CapabilityKeeper.Seal()

	app.IBCKeeper = ibckeeper.NewKeeper(
		app.appCodec,
		app.GetKey(ibcexported.StoreKey),
		nil,
		ibcStakingKeeper{k: app.StakingKeeper},
		app.UpgradeKeeper,
		scopedIBCKeeper,
		authority,
	)

	app.TransferKeeper = ibctransferkeeper.NewKeeper(
		app.appCodec,
		app.GetKey(ibctransfertypes.StoreKey),
		nil,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.ChannelKeeper,
		app.IBCKeeper.PortKeeper,
		app.AccountKeeper,
		app.BankKeeper,
		scopedTransferKeeper,
		authority,
	)

	ibcRouter := porttypes.NewRouter().
		AddRoute(ibctransfertypes.ModuleName, ibctransfer.NewIBCModule(app.TransferKeeper))
	app.IBCKeeper.SetRouter(ibcRouter)

	app.ScopedIBCKeeper = scopedIBCKeeper
	app.ScopedTransferKeeper = scopedTransferKeeper

	sdkstakingtypes.RegisterQueryServer(app.GRPCQueryRouter(), &sdkStakingQuery{k: app.StakingKeeper})

	return app.RegisterModules(
		capability.NewAppModule(app.appCodec, *app.CapabilityKeeper, false),
		ibc.NewAppModule(app.IBCKeeper),
		ibctransfer.NewAppModule(app.TransferKeeper),
		ibctm.NewAppModule(),
	)
}

// GetMemKey returns the memory store key for the provided store key.
func (app *AxmApp) GetMemKey(storeKey string) *storetypes.MemoryStoreKey {
	key, _ := app.UnsafeFindStoreKey(storeKey).(*storetypes.MemoryStoreKey)
	return key
}

// RegisterIBC registers the IBC modules on the client side (CLI), since they
// are not part of the depinject app config.
func RegisterIBC(registry cdctypes.InterfaceRegistry) map[string]appmodule.AppModule {
	modules := map[string]appmodule.AppModule{
		capabilitytypes.ModuleName:  capability.AppModule{},
		ibcexported.ModuleName:      ibc.AppModule{},
		ibctransfertypes.ModuleName: ibctransfer.AppModule{},
		ibctm.ModuleName:            ibctm.AppModule{},
	}
	for name, m := range modules {
		module.CoreAppModuleBasicAdaptor(name, m).RegisterInterfaces(registry)
	}
	return modules
}

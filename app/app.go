package app

import (
	_ "embed"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/axiome-pro/axm-node/client/docs"
	"github.com/axiome-pro/axm-node/x/referral"
	"github.com/axiome-pro/axm-node/x/vote"
	_ "github.com/axiome-pro/axm-node/x/wasm"
	dbm "github.com/cosmos/cosmos-db"

	"cosmossdk.io/core/appconfig"
	"cosmossdk.io/depinject"
	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	upgradetypes "cosmossdk.io/x/upgrade/types"

	distrkeeper "github.com/axiome-pro/axm-node/x/distribution/keeper"
	"github.com/axiome-pro/axm-node/x/genutil"
	genutiltypes "github.com/axiome-pro/axm-node/x/genutil/types"

	slashigkeeper "github.com/axiome-pro/axm-node/x/slashing/keeper"
	stakingkeeper "github.com/axiome-pro/axm-node/x/staking/keeper"
	upgradekeeper "cosmossdk.io/x/upgrade/keeper"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/server/api"
	"github.com/cosmos/cosmos-sdk/server/config"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authzkeeper "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	consensuskeeper "github.com/cosmos/cosmos-sdk/x/consensus/keeper"

	_ "github.com/axiome-pro/axm-node/x/distribution" // import for side-effects
	_ "github.com/axiome-pro/axm-node/x/slashing"     // import for side-effects
	_ "github.com/axiome-pro/axm-node/x/staking"      // import for side-effects
	_ "cosmossdk.io/api/cosmos/tx/config/v1"        // import for side-effects
	feegrantmodule "cosmossdk.io/x/feegrant"
	_ "cosmossdk.io/x/upgrade"
	// CosmWasm imports
	_ "cosmossdk.io/x/feegrant/module" // import for side-effects
	"github.com/CosmWasm/wasmd/x/wasm"
	_ "github.com/CosmWasm/wasmd/x/wasm" // import for side-effects
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	_ "github.com/cosmos/cosmos-sdk/x/auth"           // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import for side-effects
	authz "github.com/cosmos/cosmos-sdk/x/authz"
	_ "github.com/cosmos/cosmos-sdk/x/authz/module" // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/bank"         // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/consensus"    // import for side-effects
	_ "github.com/cosmos/cosmos-sdk/x/mint"         // import for side-effects
)

// DefaultNodeHome default home directories for the application daemon
var DefaultNodeHome string

//go:embed app.yaml
var AppConfigYAML []byte

var (
	_ runtime.AppI            = (*AxmApp)(nil)
	_ servertypes.Application = (*AxmApp)(nil)
)

// AxmApp extends an ABCI application, but with most of its parameters exported.
// They are exported for convenience in creating helper functions, as object
// capabilities aren't needed for testing.
type AxmApp struct {
	*runtime.App
	legacyAmino       *codec.LegacyAmino
	appCodec          codec.Codec
	txConfig          client.TxConfig
	interfaceRegistry codectypes.InterfaceRegistry

	// keepers
	AccountKeeper         authkeeper.AccountKeeper
	BankKeeper            bankkeeper.Keeper
	StakingKeeper         *stakingkeeper.Keeper
	SlashingKeeper        slashigkeeper.Keeper
	DistrKeeper           distrkeeper.Keeper
	ConsensusParamsKeeper consensuskeeper.Keeper
	ReferralKeeper        referral.Keeper
	VoteKeeper            vote.Keeper
	UpgradeKeeper         *upgradekeeper.Keeper
	WasmKeeper            wasm.Keeper

	AuthzKeeper authzkeeper.Keeper

	// simulation manager
	sm *module.SimulationManager
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	DefaultNodeHome = filepath.Join(userHomeDir, ".axmd")
}

// AppConfig returns the default app config.
func AppConfig() depinject.Config {
	return depinject.Configs(
		appconfig.LoadYAML(AppConfigYAML),
		depinject.Supply(
			// supply custom module basics
			map[string]module.AppModuleBasic{
				genutiltypes.ModuleName: genutil.NewAppModuleBasic(genutiltypes.DefaultMessageValidator),
			},
		),
	)
}

// NewAxmApp returns a reference to an initialized AxmApp.
func NewAxmApp(
	logger log.Logger,
	db dbm.DB,
	traceStore io.Writer,
	loadLatest bool,
	appOpts servertypes.AppOptions,
	baseAppOptions ...func(*baseapp.BaseApp),
) (*AxmApp, error) {
	var (
		app        = &AxmApp{}
		appBuilder *runtime.AppBuilder
	)

	if err := depinject.Inject(
		depinject.Configs(
			AppConfig(),
			depinject.Supply(
				logger,
				appOpts,
			),
		),
		&appBuilder,
		&app.appCodec,
		&app.legacyAmino,
		&app.txConfig,
		&app.interfaceRegistry,
		&app.AccountKeeper,
		&app.BankKeeper,
		&app.StakingKeeper,
		&app.SlashingKeeper,
		&app.DistrKeeper,
		&app.ReferralKeeper,
		&app.VoteKeeper,
		&app.UpgradeKeeper,
		&app.ConsensusParamsKeeper,
		&app.WasmKeeper,
		&app.AuthzKeeper,
	); err != nil {
		return nil, err
	}

	app.App = appBuilder.Build(db, traceStore, baseAppOptions...)

	// register streaming services
	if err := app.RegisterStreamingServices(appOpts, app.kvStoreKeys()); err != nil {
		return nil, err
	}

	/****  Module Options ****/

	app.RegisterUpgradeHandlers()

	// Configure store loader for in-place store upgrades (adding new KV stores)
	upgradeInfo, err := app.UpgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		return nil, err
	}
	if upgradeInfo.Name == UpgradeNamev200 && !app.UpgradeKeeper.IsSkipHeight(upgradeInfo.Height) {
		// We are adding new stores in v2.0.0: wasm, authz, feegrant
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{
				wasmtypes.StoreKey,
				authz.ModuleName,
				feegrantmodule.ModuleName,
			},
		}
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(upgradeInfo.Height, &storeUpgrades))
	}

	// create the simulation manager and define the order of the modules for deterministic simulations
	// NOTE: this is not required apps that don't use the simulator for fuzz testing transactions
	app.sm = module.NewSimulationManagerFromAppModules(app.ModuleManager.Modules, make(map[string]module.AppModuleSimulation, 0))
	app.sm.RegisterStoreDecoders()

	if err := app.Load(loadLatest); err != nil {
		return nil, err
	}

	return app, nil
}

// LegacyAmino returns AxmApp's amino codec.
func (app *AxmApp) LegacyAmino() *codec.LegacyAmino {
	return app.legacyAmino
}

// GetKey returns the KVStoreKey for the provided store key.
func (app *AxmApp) GetKey(storeKey string) *storetypes.KVStoreKey {
	sk := app.UnsafeFindStoreKey(storeKey)
	kvStoreKey, ok := sk.(*storetypes.KVStoreKey)
	if !ok {
		return nil
	}
	return kvStoreKey
}

func (app *AxmApp) kvStoreKeys() map[string]*storetypes.KVStoreKey {
	keys := make(map[string]*storetypes.KVStoreKey)
	for _, k := range app.GetStoreKeys() {
		if kv, ok := k.(*storetypes.KVStoreKey); ok {
			keys[kv.Name()] = kv
		}
	}

	return keys
}

// SimulationManager implements the SimulationApp interface
func (app *AxmApp) SimulationManager() *module.SimulationManager {
	return app.sm
}

// RegisterAPIRoutes registers all application module routes with the provided
// API server.
func (app *AxmApp) RegisterAPIRoutes(apiSvr *api.Server, apiConfig config.APIConfig) {
	app.App.RegisterAPIRoutes(apiSvr, apiConfig)

	// register swagger API in app.go so that other applications can override easily
	if apiConfig.Swagger {
		root, err := fs.Sub(docs.SwaggerUI, "swagger-ui")
		if err != nil {
			panic(err)
		}

		staticServer := http.FileServer(http.FS(root))
		apiSvr.Router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", staticServer))
	}
}

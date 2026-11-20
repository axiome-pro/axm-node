// upgrade-sim applies the v2.3.0-ibc upgrade offline to a COPY of a node home
// (application.db + wasm) and checks the result. It does what the chain does
// at the upgrade height — store loader adds the IBC stores, the upgrade
// handler runs migrations (InitGenesis of capability/ibc/transfer) — and then
// commits one version, so it shows whether the upgrade works on real state.
//
// Never point it at a live node's home: it writes a new version into the DB.
//
//	upgrade-sim --home /path/to/copy --chain-id axiome-1
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/store/rootmulti"
	upgradetypes "cosmossdk.io/x/upgrade/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	dbm "github.com/cosmos/cosmos-db"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	host "github.com/cosmos/ibc-go/v8/modules/core/24-host"

	"github.com/axiome-pro/axm-node/app"
	"github.com/axiome-pro/axm-node/app/params"
)

type opts map[string]interface{}

func (o opts) Get(k string) interface{} { return o[k] }

func must(err error, what string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "FAIL %s: %v\n", what, err)
		os.Exit(1)
	}
}

func open(home string) (dbm.DB, *app.AxmApp) {
	db, err := dbm.NewGoLevelDB("application", filepath.Join(home, "data"), nil)
	must(err, "open application.db")
	a, err := app.NewAxmApp(log.NewLogger(os.Stderr, log.LevelOption(0)), db, nil, true,
		opts{flags.FlagHome: home})
	must(err, "load app")
	return db, a
}

func main() {
	home := flag.String("home", "", "node home (a copy!)")
	chainID := flag.String("chain-id", "axiome-1", "chain id")
	flag.Parse()
	params.SetAddressPrefixes()
	if *home == "" {
		fmt.Fprintln(os.Stderr, "--home required")
		os.Exit(2)
	}

	// 1. latest committed version of the real state
	ldb, err := dbm.NewGoLevelDB("application", filepath.Join(*home, "data"), nil)
	must(err, "open application.db")
	h := rootmulti.GetLatestVersion(ldb)
	ldb.Close()
	up := h + 1
	fmt.Printf("state height %d, upgrade at %d\n", h, up)

	// 2. what the upgrade module writes to disk when the chain halts
	info, _ := json.Marshal(map[string]interface{}{"name": app.UpgradeNameV230, "height": up})
	must(os.WriteFile(filepath.Join(*home, "data", "upgrade-info.json"), info, 0o644), "upgrade-info")

	// 3. load with the IBC store loader and run the upgrade handler
	db, a := open(*home)
	before := a.LastCommitID()
	fmt.Printf("before: version %d app hash %X\n", before.Version, before.Hash)

	hdr := cmtproto.Header{ChainID: *chainID, Height: up, Time: time.Now().UTC()}
	ctx := a.NewUncachedContext(false, hdr)
	supplyBefore := a.BankKeeper.GetSupply(ctx, "uaxm")

	plan := upgradetypes.Plan{Name: app.UpgradeNameV230, Height: up}
	t0 := time.Now()
	must(a.UpgradeKeeper.ApplyUpgrade(ctx, plan), "apply upgrade")
	cid := a.CommitMultiStore().Commit()
	fmt.Printf("upgrade applied in %s, new version %d app hash %X\n", time.Since(t0).Round(time.Millisecond), cid.Version, cid.Hash)

	check(a, a.NewUncachedContext(false, hdr), supplyBefore)
	db.Close()

	// 4. reopen as a restarted node would (store loader no longer applies)
	db, a = open(*home)
	last := a.LastCommitID()
	fmt.Printf("reopened: version %d app hash %X\n", last.Version, last.Hash)
	if last.Version != cid.Version || string(last.Hash) != string(cid.Hash) {
		must(fmt.Errorf("reopened state differs from committed"), "reopen")
	}
	ctx = a.NewUncachedContext(false, cmtproto.Header{ChainID: *chainID, Height: up + 1, Time: time.Now().UTC()})
	// what capability's BeginBlocker does in the first block after a restart
	a.CapabilityKeeper.InitMemStore(ctx)
	check(a, ctx, supplyBefore)
	db.Close()
	fmt.Println("OK")
}

func check(a *app.AxmApp, ctx sdk.Context, supplyBefore sdk.Coin) {
	vm, err := a.UpgradeKeeper.GetModuleVersionMap(ctx)
	must(err, "version map")
	fmt.Printf("  module versions: capability=%d ibc=%d transfer=%d\n", vm["capability"], vm["ibc"], vm["transfer"])
	for _, m := range []string{"capability", "ibc", "transfer"} {
		if vm[m] == 0 {
			must(fmt.Errorf("module %s missing from version map", m), "version map")
		}
	}

	cp := a.IBCKeeper.ClientKeeper.GetParams(ctx)
	fmt.Printf("  ibc allowed clients: %v\n", cp.AllowedClients)
	tp := a.TransferKeeper.GetParams(ctx)
	_, bound := a.ScopedTransferKeeper.GetCapability(ctx, host.PortPath(ibctransfertypes.PortID))
	fmt.Printf("  transfer send=%v receive=%v, port %q bound=%v\n",
		tp.SendEnabled, tp.ReceiveEnabled, a.TransferKeeper.GetPort(ctx), bound)
	if !bound {
		must(fmt.Errorf("transfer port not bound"), "port")
	}
	if a.AccountKeeper.GetModuleAddress(ibctransfertypes.ModuleName) == nil {
		must(fmt.Errorf("no transfer module address"), "module account")
	}
	s := a.BankKeeper.GetSupply(ctx, "uaxm")
	fmt.Printf("  uaxm supply %s (before %s)\n", s.Amount, supplyBefore.Amount)
	if !s.Amount.Equal(supplyBefore.Amount) {
		must(fmt.Errorf("supply changed"), "supply")
	}
}

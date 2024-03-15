package testutil

import (
	appconfigurator "github.com/axiome-pro/axm-node/testutil/configurator"
	_ "github.com/axiome-pro/axm-node/x/distribution" // import as blank for app wiring
	_ "github.com/axiome-pro/axm-node/x/genutil"      // import as blank for app wiring
	_ "github.com/axiome-pro/axm-node/x/staking"      // import as blank for app wiring
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
	_ "github.com/cosmos/cosmos-sdk/x/auth"           // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/auth/tx/config" // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/bank"           // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/consensus"      // import as blank for app wiring
	_ "github.com/cosmos/cosmos-sdk/x/params"         // import as blank for app wiring
)

var AppConfig = configurator.NewAppConfig(
	configurator.AuthModule(),
	configurator.BankModule(),
	appconfigurator.StakingModule(),
	configurator.TxModule(),
	configurator.ConsensusModule(),
	configurator.ParamsModule(),
	appconfigurator.GenutilModule(),
	appconfigurator.DistributionModule(),
)

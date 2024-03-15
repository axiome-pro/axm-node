package configurator

import (
	distrmodulev1 "github.com/axiome-pro/axm-node/api/axiome/distribution/module/v1"
	genutilmodulev1 "github.com/axiome-pro/axm-node/api/axiome/genutil/module/v1"
	slashingmodulev1 "github.com/axiome-pro/axm-node/api/axiome/slashing/module/v1"
	stakingmodulev1 "github.com/axiome-pro/axm-node/api/axiome/staking/module/v1"
	appv1alpha1 "cosmossdk.io/api/cosmos/app/v1alpha1"
	"cosmossdk.io/core/appconfig"
	"github.com/cosmos/cosmos-sdk/testutil/configurator"
)

func StakingModule() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs["staking"] = &appv1alpha1.ModuleConfig{
			Name:   "staking",
			Config: appconfig.WrapAny(&stakingmodulev1.Module{}),
		}
	}
}

func SlashingModule() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs["slashing"] = &appv1alpha1.ModuleConfig{
			Name:   "slashing",
			Config: appconfig.WrapAny(&slashingmodulev1.Module{}),
		}
	}
}

func DistributionModule() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs["distribution"] = &appv1alpha1.ModuleConfig{
			Name:   "distribution",
			Config: appconfig.WrapAny(&distrmodulev1.Module{}),
		}
	}
}

func GenutilModule() configurator.ModuleOption {
	return func(config *configurator.Config) {
		config.ModuleConfigs["genutil"] = &appv1alpha1.ModuleConfig{
			Name:   "genutil",
			Config: appconfig.WrapAny(&genutilmodulev1.Module{}),
		}
	}
}

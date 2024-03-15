package simulation

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/axiome-pro/axm-node/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/types/simulation"
)

// Simulation parameter constants
const (
	unbondingTime     = "unbonding_time"
	maxValidators     = "max_validators"
	historicalEntries = "historical_entries"
	maxMonthlyPoints  = "max_monthly_points"
	minSelfDelegation = "min_self_delegation"
)

// genUnbondingTime returns randomized UnbondingTime
func genUnbondingTime(r *rand.Rand) (ubdTime time.Duration) {
	return time.Duration(simulation.RandIntBetween(r, 60, 60*60*24*3*2)) * time.Second
}

// genMaxValidators returns randomized MaxValidators
func genMaxValidators(r *rand.Rand) (maxValidators uint32) {
	return uint32(r.Intn(250) + 1)
}

// getHistEntries returns randomized HistoricalEntries between 0-100.
func getHistEntries(r *rand.Rand) uint32 {
	return uint32(r.Intn(int(types.DefaultHistoricalEntries + 1)))
}

// getMaxMonthlyPoints returns randomized HistoricalEntries between 0-100.
func getMaxMonthlyPoints(r *rand.Rand) uint64 {
	return uint64(r.Intn(int(types.DefaultMaximumMonthlyPoints + 1)))
}

// getHistEntries returns randomized HistoricalEntries between 0-100.
func getMinSelfDelegation(r *rand.Rand) sdkmath.Int {
	return sdkmath.NewInt(r.Int63n(types.DefaultMinSelfDelegation.AddRaw(1).Int64()))
}

// RandomizedGenState generates a random GenesisState for staking
func RandomizedGenState(simState *module.SimulationState) {
	// params
	var (
		unbondTime            time.Duration
		maxVals               uint32
		histEntries           uint32
		validatorEmissionRate sdkmath.LegacyDec
		maximumMonthlyPoints  uint64
		minimumSelfDelegation sdkmath.Int
	)

	simState.AppParams.GetOrGenerate(unbondingTime, &unbondTime, simState.Rand, func(r *rand.Rand) { unbondTime = genUnbondingTime(r) })

	simState.AppParams.GetOrGenerate(maxValidators, &maxVals, simState.Rand, func(r *rand.Rand) { maxVals = genMaxValidators(r) })

	simState.AppParams.GetOrGenerate(historicalEntries, &histEntries, simState.Rand, func(r *rand.Rand) { histEntries = getHistEntries(r) })

	simState.AppParams.GetOrGenerate(maxMonthlyPoints, &maximumMonthlyPoints, simState.Rand, func(r *rand.Rand) { maximumMonthlyPoints = getMaxMonthlyPoints(r) })

	simState.AppParams.GetOrGenerate(maxMonthlyPoints, &minimumSelfDelegation, simState.Rand, func(r *rand.Rand) { minimumSelfDelegation = getMinSelfDelegation(r) })

	// NOTE: the slashing module need to be defined after the staking module on the
	// NewSimulationManager constructor for this to work
	simState.UnbondTime = unbondTime
	params := types.NewParams(simState.UnbondTime, simState.UnbondTime, maxVals, 7, histEntries,
		simState.BondDenom, maximumMonthlyPoints, validatorEmissionRate, []*types.EmissionRange{&types.DefaultEmissionRange},
		minimumSelfDelegation)

	// validators & delegations
	var (
		validators  []types.Validator
		delegations []types.Delegation
	)

	valAddrs := make([]sdk.ValAddress, simState.NumBonded)

	for i := 0; i < int(simState.NumBonded); i++ {
		valAddr := sdk.ValAddress(simState.Accounts[i].Address)
		valAddrs[i] = valAddr

		validator, err := types.NewValidator(valAddr.String(), simState.Accounts[i].ConsKey.PubKey(), types.Description{})
		if err != nil {
			panic(err)
		}
		validator.Tokens = simState.InitialStake
		validator.DelegatorShares = sdkmath.LegacyNewDecFromInt(simState.InitialStake)

		delegation := types.NewDelegation(simState.Accounts[i].Address.String(), valAddr.String(), sdkmath.LegacyNewDecFromInt(simState.InitialStake))

		validators = append(validators, validator)
		delegations = append(delegations, delegation)
	}

	stakingGenesis := types.NewGenesisState(params, validators, delegations)

	bz, err := json.MarshalIndent(&stakingGenesis.Params, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Selected randomly generated staking parameters:\n%s\n", bz)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(stakingGenesis)
}

package keeper_test

import (
	"github.com/axiome-pro/axm-node/app/params"
	"testing"
	"time"

	abci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"

	"github.com/axiome-pro/axm-node/x/distribution"
	"github.com/axiome-pro/axm-node/x/distribution/keeper"
	distrtestutil "github.com/axiome-pro/axm-node/x/distribution/testutil"
	disttypes "github.com/axiome-pro/axm-node/x/distribution/types"
	stakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
	"github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func TestAllocateTokensToValidatorWithCommission(t *testing.T) {
	ctrl := gomock.NewController(t)
	key := storetypes.NewKVStoreKey(disttypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(distribution.AppModuleBasic{})
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})

	bankKeeper := distrtestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := distrtestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := distrtestutil.NewMockAccountKeeper(ctrl)

	valCodec := address.NewBech32Codec(params.Bech32PrefixValAddr)

	accountKeeper.EXPECT().GetModuleAddress("distribution").Return(distrAcc.GetAddress())
	stakingKeeper.EXPECT().ValidatorAddressCodec().Return(valCodec).AnyTimes()

	distrKeeper := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		"fee_collector",
		authtypes.NewModuleAddress("gov").String(),
	)

	require.NoError(t, distrKeeper.Params.Set(ctx, disttypes.DefaultParams()))

	// create validator with 50% commission
	val, err := distrtestutil.CreateValidator(valConsPk0, math.NewInt(100))
	require.NoError(t, err)
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk0)).Return(val, nil).AnyTimes()
	bankKeeper.EXPECT().BurnCoins(gomock.Any(), "distribution", gomock.Any()).Return(nil).AnyTimes()

	// allocate tokens
	tokens := sdk.DecCoins{
		{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDec(10)},
	}
	require.NoError(t, distrKeeper.AllocateTokensToValidator(ctx, val, tokens))

	valBz, err := valCodec.StringToBytes(val.GetOperator())
	require.NoError(t, err)

	valCommission, err := distrKeeper.GetValidatorAccumulatedCommission(ctx, valBz)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(3, 2)}}, valCommission.Commission)

	// check current rewards
	currentRewards, err := distrKeeper.GetValidatorCurrentRewards(ctx, valBz)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(997, 2)}}, currentRewards.Rewards)
}

func TestAllocateTokensToManyValidators(t *testing.T) {
	ctrl := gomock.NewController(t)
	key := storetypes.NewKVStoreKey(disttypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(distribution.AppModuleBasic{})
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})

	bankKeeper := distrtestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := distrtestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := distrtestutil.NewMockAccountKeeper(ctrl)

	feeCollectorAcc := authtypes.NewEmptyModuleAccount("fee_collector")
	accountKeeper.EXPECT().GetModuleAddress("distribution").Return(distrAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAccount(gomock.Any(), "fee_collector").Return(feeCollectorAcc)
	stakingKeeper.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec(params.Bech32PrefixValAddr)).AnyTimes()
	bankKeeper.EXPECT().BurnCoins(gomock.Any(), "distribution", gomock.Any()).Return(nil).AnyTimes()

	distrKeeper := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		"fee_collector",
		authtypes.NewModuleAddress("gov").String(),
	)

	// reset fee pool & set params
	require.NoError(t, distrKeeper.Params.Set(ctx, disttypes.DefaultParams()))
	require.NoError(t, distrKeeper.FeePool.Set(ctx, disttypes.InitialFeePool()))

	// create validator with 50% commission
	valAddr0 := sdk.ValAddress(valConsAddr0)
	val0, err := distrtestutil.CreateValidator(valConsPk0, math.NewInt(100))
	require.NoError(t, err)
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk0)).Return(val0, nil).AnyTimes()

	// create second validator with 0% commission
	valAddr1 := sdk.ValAddress(valConsAddr1)
	val1, err := distrtestutil.CreateValidator(valConsPk1, math.NewInt(100))
	require.NoError(t, err)
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk1)).Return(val1, nil).AnyTimes()

	abciValA := abci.Validator{
		Address: valConsPk0.Address(),
		Power:   100,
	}
	abciValB := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   100,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	val0OutstandingRewards, err := distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0OutstandingRewards.Rewards.IsZero())

	val1OutstandingRewards, err := distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1OutstandingRewards.Rewards.IsZero())

	feePool, err := distrKeeper.FeePool.Get(ctx)
	require.NoError(t, err)
	require.True(t, feePool.CommunityPool.IsZero())

	val0Commission, err := distrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0Commission.Commission.IsZero())

	val1Commission, err := distrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1Commission.Commission.IsZero())

	val0CurrentRewards, err := distrKeeper.GetValidatorCurrentRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0CurrentRewards.Rewards.IsZero())

	val1CurrentRewards, err := distrKeeper.GetValidatorCurrentRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1CurrentRewards.Rewards.IsZero())

	// allocate tokens as if both had voted and second was proposer
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(100)))
	bankKeeper.EXPECT().GetAllBalances(gomock.Any(), feeCollectorAcc.GetAddress()).Return(fees)
	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), "fee_collector", disttypes.ModuleName, fees)

	votes := []abci.VoteInfo{
		{
			Validator: abciValA,
		},
		{
			Validator: abciValB,
		},
	}
	require.NoError(t, distrKeeper.AllocateTokens(ctx, 200, votes))

	// 98 outstanding rewards (100 less 2 to community pool)
	val0OutstandingRewards, err = distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(35, 0)}}, val0OutstandingRewards.Rewards)

	val1OutstandingRewards, err = distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(35, 0)}}, val1OutstandingRewards.Rewards)

	// 2 community pool coins
	feePool, err = distrKeeper.FeePool.Get(ctx)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins(nil), feePool.CommunityPool)

	// 50% commission for first proposer, (0.5 * 98%) * 100 / 2 = 23.25
	val0Commission, err = distrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr0)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(105, 3)}}, val0Commission.Commission)

	val1Commission, err = distrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(105, 3)}}, val1Commission.Commission)

	val0CurrentRewards, err = distrKeeper.GetValidatorCurrentRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(34895, 3)}}, val0CurrentRewards.Rewards)

	val1CurrentRewards, err = distrKeeper.GetValidatorCurrentRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.Equal(t, sdk.DecCoins{{Denom: sdk.DefaultBondDenom, Amount: math.LegacyNewDecWithPrec(34895, 3)}}, val1CurrentRewards.Rewards)
}

func TestAllocateTokensTruncation(t *testing.T) {
	ctrl := gomock.NewController(t)
	key := storetypes.NewKVStoreKey(disttypes.StoreKey)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	encCfg := moduletestutil.MakeTestEncodingConfig(distribution.AppModuleBasic{})
	ctx := testCtx.Ctx.WithBlockHeader(cmtproto.Header{Time: time.Now()})

	bankKeeper := distrtestutil.NewMockBankKeeper(ctrl)
	stakingKeeper := distrtestutil.NewMockStakingKeeper(ctrl)
	accountKeeper := distrtestutil.NewMockAccountKeeper(ctrl)

	feeCollectorAcc := authtypes.NewEmptyModuleAccount("fee_collector")
	accountKeeper.EXPECT().GetModuleAddress("distribution").Return(distrAcc.GetAddress())
	accountKeeper.EXPECT().GetModuleAccount(gomock.Any(), "fee_collector").Return(feeCollectorAcc)
	stakingKeeper.EXPECT().ValidatorAddressCodec().Return(address.NewBech32Codec(params.Bech32PrefixValAddr)).AnyTimes()
	bankKeeper.EXPECT().BurnCoins(gomock.Any(), "distribution", gomock.Any()).Return(nil).AnyTimes()

	distrKeeper := keeper.NewKeeper(
		encCfg.Codec,
		storeService,
		accountKeeper,
		bankKeeper,
		stakingKeeper,
		"fee_collector",
		authtypes.NewModuleAddress("gov").String(),
	)

	// reset fee pool
	require.NoError(t, distrKeeper.FeePool.Set(ctx, disttypes.InitialFeePool()))
	require.NoError(t, distrKeeper.Params.Set(ctx, disttypes.DefaultParams()))

	// create validator with 10% commission
	valAddr0 := sdk.ValAddress(valConsAddr0)
	val0, err := distrtestutil.CreateValidator(valConsPk0, math.NewInt(100))
	require.NoError(t, err)
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk0)).Return(val0, nil).AnyTimes()

	// create second validator with 10% commission
	valAddr1 := sdk.ValAddress(valConsAddr1)
	val1, err := distrtestutil.CreateValidator(valConsPk1, math.NewInt(100))
	require.NoError(t, err)
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk1)).Return(val1, nil).AnyTimes()

	// create third validator with 10% commission
	valAddr2 := sdk.ValAddress(valConsAddr2)
	val2, err := stakingtypes.NewValidator(sdk.ValAddress(valConsAddr2).String(), valConsPk1, stakingtypes.Description{})
	require.NoError(t, err)
	stakingKeeper.EXPECT().ValidatorByConsAddr(gomock.Any(), sdk.GetConsAddress(valConsPk2)).Return(val2, nil).AnyTimes()

	abciValA := abci.Validator{
		Address: valConsPk0.Address(),
		Power:   11,
	}
	abciValB := abci.Validator{
		Address: valConsPk1.Address(),
		Power:   10,
	}
	abciValC := abci.Validator{
		Address: valConsPk2.Address(),
		Power:   10,
	}

	// assert initial state: zero outstanding rewards, zero community pool, zero commission, zero current rewards
	val0OutstandingRewards, err := distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0OutstandingRewards.Rewards.IsZero())

	val1OutstandingRewards, err := distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1OutstandingRewards.Rewards.IsZero())

	feePool, err := distrKeeper.FeePool.Get(ctx)
	require.NoError(t, err)
	require.True(t, feePool.CommunityPool.IsZero())

	val0Commission, err := distrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0Commission.Commission.IsZero())

	val1Commission, err := distrKeeper.GetValidatorAccumulatedCommission(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1Commission.Commission.IsZero())

	val0CurrentRewards, err := distrKeeper.GetValidatorCurrentRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0CurrentRewards.Rewards.IsZero())

	val1CurrentRewards, err := distrKeeper.GetValidatorCurrentRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1CurrentRewards.Rewards.IsZero())

	// allocate tokens as if both had voted and second was proposer
	fees := sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, math.NewInt(634195840)))
	bankKeeper.EXPECT().GetAllBalances(gomock.Any(), feeCollectorAcc.GetAddress()).Return(fees)
	bankKeeper.EXPECT().SendCoinsFromModuleToModule(gomock.Any(), "fee_collector", disttypes.ModuleName, fees)

	votes := []abci.VoteInfo{
		{
			Validator: abciValA,
		},
		{
			Validator: abciValB,
		},
		{
			Validator: abciValC,
		},
	}
	require.NoError(t, distrKeeper.AllocateTokens(ctx, 31, votes))

	val0OutstandingRewards, err = distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr0)
	require.NoError(t, err)
	require.True(t, val0OutstandingRewards.Rewards.IsValid())

	val1OutstandingRewards, err = distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr1)
	require.NoError(t, err)
	require.True(t, val1OutstandingRewards.Rewards.IsValid())

	val2OutstandingRewards, err := distrKeeper.GetValidatorOutstandingRewards(ctx, valAddr2)
	require.NoError(t, err)
	require.True(t, val2OutstandingRewards.Rewards.IsValid())
}

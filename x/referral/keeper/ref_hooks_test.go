package keeper

import (
	"context"
	"testing"

	coreaddress "cosmossdk.io/core/address"
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/axiome-pro/axm-node/x/referral/types"
	sdkaddress "github.com/cosmos/cosmos-sdk/codec/address"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	"github.com/stretchr/testify/require"
)

type stubAccountKeeper struct {
	codec coreaddress.Codec
}

func (s stubAccountKeeper) AddressCodec() coreaddress.Codec {
	return s.codec
}

func (s stubAccountKeeper) GetAccount(context.Context, sdk.AccAddress) sdk.AccountI {
	return nil
}

type stubWasmKeeper struct {
	contracts map[string]bool
}

func (s stubWasmKeeper) HasContractInfo(_ context.Context, contractAddress sdk.AccAddress) bool {
	return s.contracts[string(contractAddress)]
}

func mustAddress(t *testing.T, codec coreaddress.Codec, raw []byte) (string, []byte) {
	t.Helper()

	addr, err := codec.BytesToString(raw)
	require.NoError(t, err)

	addrBytes, err := codec.StringToBytes(addr)
	require.NoError(t, err)

	return addr, addrBytes
}

func TestCheckDelegationAvailableRequiresReferralForNonContracts(t *testing.T) {
	key := storetypes.NewKVStoreKey(types.ModuleName)
	storeService := runtime.NewKVStoreService(key)
	testCtx := testutil.DefaultContextWithDB(t, key, storetypes.NewTransientStoreKey("transient_test"))
	ctx := testCtx.Ctx
	encCfg := moduletestutil.MakeTestEncodingConfig()
	accCodec := sdkaddress.NewBech32Codec("axm")

	contractAddr, contractAddrBytes := mustAddress(t, accCodec, []byte("contract-address-000"))
	userAddr, _ := mustAddress(t, accCodec, []byte("user-address-000000"))
	referredAddr, _ := mustAddress(t, accCodec, []byte("referred-address-00"))

	k := NewKeeper(
		encCfg.Codec,
		storeService,
		stubAccountKeeper{codec: accCodec},
		stubWasmKeeper{contracts: map[string]bool{string(contractAddrBytes): true}},
		nil,
		nil,
		sdk.AccAddress{},
		types.ReferralAccountName,
	)

	require.NoError(t, k.set(ctx, referredAddr, types.NewInfo("", math.ZeroInt())))

	err := k.Hooks().CheckDelegationAvailable(ctx, contractAddr, "axmvaloper1validator")
	require.NoError(t, err)

	err = k.Hooks().CheckDelegationAvailable(ctx, userAddr, "axmvaloper1validator")
	require.Error(t, err)
	require.ErrorIs(t, err, types.ErrNotFound)

	err = k.Hooks().CheckDelegationAvailable(ctx, referredAddr, "axmvaloper1validator")
	require.NoError(t, err)
}

package app

import (
	"testing"
	"time"

	"cosmossdk.io/math"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkstakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	clienttypes "github.com/cosmos/ibc-go/v8/modules/core/02-client/types"

	axmstakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
)

// The adapter exists only to satisfy this interface. If ibc-go widens it, this
// line fails to compile rather than failing at runtime on the first client.
var _ clienttypes.StakingKeeper = ibcStakingKeeper{}

// TestBondStatusValuesMatch guards the assumption the adapter relies on when it
// converts the status with a plain numeric cast. The two enums are declared in
// different proto packages, so nothing but this test keeps them aligned.
func TestBondStatusValuesMatch(t *testing.T) {
	cases := []struct {
		axm axmstakingtypes.BondStatus
		sdk sdkstakingtypes.BondStatus
	}{
		{axmstakingtypes.Unspecified, sdkstakingtypes.Unspecified},
		{axmstakingtypes.Unbonded, sdkstakingtypes.Unbonded},
		{axmstakingtypes.Unbonding, sdkstakingtypes.Unbonding},
		{axmstakingtypes.Bonded, sdkstakingtypes.Bonded},
	}

	for _, c := range cases {
		if int32(c.axm) != int32(c.sdk) {
			t.Fatalf("bond status drifted: axiome %d (%s) vs cosmos-sdk %d (%s)",
				c.axm, c.axm, c.sdk, c.sdk)
		}
	}
}

func TestConvertValidatorsEmpty(t *testing.T) {
	if got := convertValidators(nil); got != nil {
		t.Fatalf("expected nil for no validators, got %v", got)
	}
	if got := convertValidators([]axmstakingtypes.Validator{}); got != nil {
		t.Fatalf("expected nil for empty slice, got %v", got)
	}
}

func TestConvertValidatorsCopiesFields(t *testing.T) {
	pubkey, err := codectypes.NewAnyWithValue(&sdkstakingtypes.Description{})
	if err != nil {
		t.Fatalf("building a placeholder Any: %v", err)
	}
	unbondingAt := time.Date(2026, 8, 6, 12, 0, 0, 0, time.UTC)

	in := axmstakingtypes.Validator{
		OperatorAddress: "axmvaloper1test",
		ConsensusPubkey: pubkey,
		Jailed:          true,
		Status:          axmstakingtypes.Bonded,
		Tokens:          math.NewInt(1234567),
		DelegatorShares: math.LegacyNewDec(42),
		Description: axmstakingtypes.Description{
			Moniker:         "moniker",
			Identity:        "identity",
			Website:         "https://example.invalid",
			SecurityContact: "security@example.invalid",
			Details:         "details",
		},
		UnbondingHeight:         99,
		UnbondingTime:           unbondingAt,
		MinSelfDelegation:       math.NewInt(7),
		UnbondingOnHoldRefCount: 3,
		UnbondingIds:            []uint64{1, 2, 3},
	}

	out := convertValidators([]axmstakingtypes.Validator{in})
	if len(out) != 1 {
		t.Fatalf("expected one validator, got %d", len(out))
	}
	got := out[0]

	if got.OperatorAddress != in.OperatorAddress {
		t.Errorf("operator address: got %q, want %q", got.OperatorAddress, in.OperatorAddress)
	}
	if got.ConsensusPubkey != in.ConsensusPubkey {
		t.Error("consensus pubkey was not carried over")
	}
	if !got.Jailed {
		t.Error("jailed flag was lost")
	}
	if got.Status != sdkstakingtypes.Bonded {
		t.Errorf("status: got %s, want %s", got.Status, sdkstakingtypes.Bonded)
	}
	if !got.Tokens.Equal(in.Tokens) {
		t.Errorf("tokens: got %s, want %s", got.Tokens, in.Tokens)
	}
	if !got.DelegatorShares.Equal(in.DelegatorShares) {
		t.Errorf("delegator shares: got %s, want %s", got.DelegatorShares, in.DelegatorShares)
	}
	if got.Description.Moniker != in.Description.Moniker ||
		got.Description.Identity != in.Description.Identity ||
		got.Description.Website != in.Description.Website ||
		got.Description.SecurityContact != in.Description.SecurityContact ||
		got.Description.Details != in.Description.Details {
		t.Errorf("description was not copied faithfully: %+v", got.Description)
	}
	if got.UnbondingHeight != in.UnbondingHeight {
		t.Errorf("unbonding height: got %d, want %d", got.UnbondingHeight, in.UnbondingHeight)
	}
	if !got.UnbondingTime.Equal(in.UnbondingTime) {
		t.Errorf("unbonding time: got %s, want %s", got.UnbondingTime, in.UnbondingTime)
	}
	if !got.MinSelfDelegation.Equal(in.MinSelfDelegation) {
		t.Errorf("min self delegation: got %s, want %s", got.MinSelfDelegation, in.MinSelfDelegation)
	}
	if got.UnbondingOnHoldRefCount != in.UnbondingOnHoldRefCount {
		t.Errorf("unbonding on hold ref count: got %d, want %d",
			got.UnbondingOnHoldRefCount, in.UnbondingOnHoldRefCount)
	}
	if len(got.UnbondingIds) != len(in.UnbondingIds) {
		t.Errorf("unbonding ids: got %v, want %v", got.UnbondingIds, in.UnbondingIds)
	}

	// Axiome validators earn through the referral module and carry no
	// commission, so a zero commission is the honest representation.
	if !got.Commission.Rate.IsNil() && !got.Commission.Rate.IsZero() {
		t.Errorf("expected zero commission, got %s", got.Commission.Rate)
	}
}

// TestIBCStoreKeysAreDistinct is a cheap guard on the list the upgrade handler
// and the wiring share: a duplicate or an empty entry there would surface as a
// confusing mount error at the upgrade height.
func TestIBCStoreKeysAreDistinct(t *testing.T) {
	seen := make(map[string]struct{}, len(ibcStoreKeys))
	for _, k := range ibcStoreKeys {
		if k == "" {
			t.Fatal("empty store key in ibcStoreKeys")
		}
		if _, dup := seen[k]; dup {
			t.Fatalf("duplicate store key %q", k)
		}
		seen[k] = struct{}{}
	}
	if len(ibcStoreKeys) != 3 {
		t.Fatalf("expected capability, ibc and transfer store keys, got %v", ibcStoreKeys)
	}
}

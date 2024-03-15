package v2

import (
	"context"

	storetypes "cosmossdk.io/core/store"

	v1 "github.com/axiome-pro/axm-node/x/slashing/migrations/v1"
	"github.com/cosmos/cosmos-sdk/runtime"
	v2distribution "github.com/cosmos/cosmos-sdk/x/distribution/migrations/v2"
)

// MigrateStore performs in-place store migrations from v0.40 to v0.43. The
// migration includes:
//
// - Change addresses to be length-prefixed.
func MigrateStore(ctx context.Context, storeService storetypes.KVStoreService) error {
	store := runtime.KVStoreAdapter(storeService.OpenKVStore(ctx))
	v2distribution.MigratePrefixAddress(store, v1.ValidatorSigningInfoKeyPrefix)
	v2distribution.MigratePrefixAddressBytes(store, v1.ValidatorMissedBlockBitArrayKeyPrefix)
	v2distribution.MigratePrefixAddress(store, v1.AddrPubkeyRelationKeyPrefix)

	return nil
}

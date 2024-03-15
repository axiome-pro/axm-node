package types

import (
	referral "github.com/axiome-pro/axm-node/x/referral/types"
	"cosmossdk.io/core/address"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"context"
)

type ReferralKeeper interface {
	Get(ctx sdk.Context, acc string) (referral.Info, error)
}

// AccountKeeper defines the expected account keeper used for simulations (noalias)
type AccountKeeper interface {
	AddressCodec() address.Codec
	GetAccount(ctx context.Context, addr sdk.AccAddress) sdk.AccountI
}

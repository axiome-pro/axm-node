package types

import (
	"github.com/cosmos/cosmos-sdk/x/slashing/types"
)

// x/slashing module sentinel errors
var (
	ErrNoValidatorForAddress        = types.ErrNoValidatorForAddress
	ErrBadValidatorAddr             = types.ErrBadValidatorAddr
	ErrValidatorJailed              = types.ErrValidatorJailed
	ErrValidatorNotJailed           = types.ErrValidatorNotJailed
	ErrMissingSelfDelegation        = types.ErrMissingSelfDelegation
	ErrSelfDelegationTooLowToUnjail = types.ErrSelfDelegationTooLowToUnjail
	ErrNoSigningInfoFound           = types.ErrNoSigningInfoFound
	ErrValidatorTombstoned          = types.ErrValidatorTombstoned
)

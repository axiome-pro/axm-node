package types

import (
	"github.com/cosmos/cosmos-sdk/x/distribution/types"
)

// x/distribution module sentinel errors
var (
	ErrEmptyDelegatorAddr      = types.ErrEmptyDelegatorAddr
	ErrEmptyWithdrawAddr       = types.ErrEmptyWithdrawAddr
	ErrEmptyValidatorAddr      = types.ErrEmptyValidatorAddr
	ErrEmptyDelegationDistInfo = types.ErrEmptyDelegationDistInfo
	ErrNoValidatorDistInfo     = types.ErrNoValidatorDistInfo
	ErrNoValidatorCommission   = types.ErrNoValidatorCommission
	ErrSetWithdrawAddrDisabled = types.ErrSetWithdrawAddrDisabled
	ErrBadDistribution         = types.ErrBadDistribution
	ErrInvalidProposalAmount   = types.ErrInvalidProposalAmount
	ErrEmptyProposalRecipient  = types.ErrEmptyProposalRecipient
	ErrNoValidatorExists       = types.ErrNoValidatorExists
	ErrNoDelegationExists      = types.ErrNoDelegationExists
)

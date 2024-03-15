package types

import (
	"github.com/cosmos/cosmos-sdk/x/staking/types"
)

// x/staking module sentinel errors
var (
	ErrEmptyValidatorAddr              = types.ErrEmptyValidatorAddr
	ErrNoValidatorFound                = types.ErrNoValidatorFound
	ErrValidatorOwnerExists            = types.ErrValidatorOwnerExists
	ErrValidatorPubKeyExists           = types.ErrValidatorPubKeyExists
	ErrValidatorPubKeyTypeNotSupported = types.ErrValidatorPubKeyTypeNotSupported
	ErrValidatorJailed                 = types.ErrValidatorJailed
	ErrBadRemoveValidator              = types.ErrBadRemoveValidator
	ErrCommissionNegative              = types.ErrCommissionNegative
	ErrCommissionHuge                  = types.ErrCommissionHuge
	ErrCommissionGTMaxRate             = types.ErrCommissionGTMaxRate
	ErrCommissionUpdateTime            = types.ErrCommissionUpdateTime
	ErrCommissionChangeRateNegative    = types.ErrCommissionChangeRateNegative
	ErrCommissionChangeRateGTMaxRate   = types.ErrCommissionChangeRateGTMaxRate
	ErrCommissionGTMaxChangeRate       = types.ErrCommissionGTMaxChangeRate
	ErrSelfDelegationBelowMinimum      = types.ErrSelfDelegationBelowMinimum
	ErrMinSelfDelegationDecreased      = types.ErrMinSelfDelegationDecreased
	ErrEmptyDelegatorAddr              = types.ErrEmptyDelegatorAddr
	ErrNoDelegation                    = types.ErrNoDelegation
	ErrBadDelegatorAddr                = types.ErrBadDelegatorAddr
	ErrNoDelegatorForAddress           = types.ErrNoDelegatorForAddress
	ErrInsufficientShares              = types.ErrInsufficientShares
	ErrDelegationValidatorEmpty        = types.ErrDelegationValidatorEmpty
	ErrNotEnoughDelegationShares       = types.ErrNotEnoughDelegationShares
	ErrNotMature                       = types.ErrNotMature
	ErrNoUnbondingDelegation           = types.ErrNoUnbondingDelegation
	ErrMaxUnbondingDelegationEntries   = types.ErrMaxUnbondingDelegationEntries
	ErrNoRedelegation                  = types.ErrNoRedelegation
	ErrSelfRedelegation                = types.ErrSelfRedelegation
	ErrTinyRedelegationAmount          = types.ErrTinyRedelegationAmount
	ErrBadRedelegationDst              = types.ErrBadRedelegationDst
	ErrTransitiveRedelegation          = types.ErrTransitiveRedelegation
	ErrMaxRedelegationEntries          = types.ErrMaxRedelegationEntries
	ErrDelegatorShareExRateInvalid     = types.ErrDelegatorShareExRateInvalid
	ErrBothShareMsgsGiven              = types.ErrBothShareMsgsGiven
	ErrNeitherShareMsgsGiven           = types.ErrNeitherShareMsgsGiven
	ErrInvalidHistoricalInfo           = types.ErrInvalidHistoricalInfo
	ErrNoHistoricalInfo                = types.ErrNoHistoricalInfo
	ErrEmptyValidatorPubKey            = types.ErrEmptyValidatorPubKey
	ErrCommissionLTMinRate             = types.ErrCommissionLTMinRate
	ErrUnbondingNotFound               = types.ErrUnbondingNotFound
	ErrUnbondingOnHoldRefCountNegative = types.ErrUnbondingOnHoldRefCountNegative
	ErrInvalidSigner                   = types.ErrInvalidSigner
	ErrBadRedelegationSrc              = types.ErrBadRedelegationSrc
	ErrNoUnbondingType                 = types.ErrNoUnbondingType
)

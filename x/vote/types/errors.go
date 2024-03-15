package types

import (
	sdkerrors "cosmossdk.io/errors"
)

var (
	// Signer not in government list
	ErrSignerNotAllowed          = sdkerrors.Register(ModuleName, 1, "signer not in government list")
	ErrOtherActive               = sdkerrors.Register(ModuleName, 2, "other proposal is active")
	ErrAlreadyVoted              = sdkerrors.Register(ModuleName, 3, "already voted")
	ErrNoActiveProposals         = sdkerrors.Register(ModuleName, 4, "no active proposals to vote")
	ErrProposalGovernorExists    = sdkerrors.Register(ModuleName, 5, "candidate already in government list")
	ErrProposalGovernorNotExists = sdkerrors.Register(ModuleName, 6, "candidate not in government list")
	ErrProposalGovernorLast      = sdkerrors.Register(ModuleName, 7, "cannot remove the last governor")
	ErrNoActivePoll              = sdkerrors.Register(ModuleName, 8, "no active poll")
	ErrRespondentNotAllowed      = sdkerrors.Register(ModuleName, 9, "poll requirements don't match")
	ErrInvalidSigner             = sdkerrors.Register(ModuleName, 10, "invalid signer for proposed message")
	ErrUnroutableProposalMsg     = sdkerrors.Register(ModuleName, 11, "no route for proposed message")
)

package types

import "cosmossdk.io/collections"

const (
	// ModuleName is the name of the module
	ModuleName = "vote"

	// RouterKey to be used for routing msgs
	RouterKey = ModuleName
)

var (
	KeyGovernment       = []byte{0x01}
	KeyAgreedMembers    = []byte{0x02}
	KeyDisagreedMembers = []byte{0x03}
	KeyCurrentVote      = []byte{0x04}
	KeyTotalVotes       = []byte{0x05}
	KeyTotalAgreed      = []byte{0x06}
	KeyTotalDisagreed   = []byte{0x07}
	KeyStartBlock       = []byte{0x08}
	KeyHistoryPrefix    = []byte{0x09}

	KeyPollPrefix   = []byte{0x0A}
	KeyPollCurrent  = []byte{0x0B}
	KeyPollAnswers  = []byte{0x0C}
	KeyPollYesCount = []byte{0x0D}
	KeyPollNoCount  = []byte{0x0E}
	KeyPollHistory  = []byte{0x0F}

	KeyParams           = collections.NewPrefix([]byte{0x10})
	KeyProposalSchedule = []byte{0x11}
	KeyPollSchedule     = []byte{0x12}

	ValueYes = []byte{0x01}
	ValueNo  = []byte{0x00}
)

func GetPollPrefixedKey(key []byte) []byte {
	return append(KeyPollPrefix, key...)
}

func GetPollAnswersPrefixedKey(key []byte) []byte {
	return append(GetPollPrefixedKey(KeyPollAnswers), key...)
}

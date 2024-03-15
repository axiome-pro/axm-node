package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterCodec registers concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// Messages
	cdc.RegisterConcrete(MsgPropose{}, ModuleName+"/CreateProposal", nil)
	cdc.RegisterConcrete(MsgVote{}, ModuleName+"/ProposalVote", nil)
	cdc.RegisterConcrete(MsgStartPoll{}, ModuleName+"/StartPoll", nil)
	cdc.RegisterConcrete(MsgAnswerPoll{}, ModuleName+"/AnswerPoll", nil)
	// Other
	cdc.RegisterConcrete(Poll{}, ModuleName+"/Poll", nil)
	cdc.RegisterConcrete(Params{}, ModuleName+"/Params", nil)
}

func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgPropose{},
		&MsgVote{},
		&MsgStartPoll{},
		&MsgAnswerPoll{},
		&Proposal{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}

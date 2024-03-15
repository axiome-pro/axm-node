package types

import sdk "github.com/cosmos/cosmos-sdk/types"

var (
	_ sdk.Msg = (*MsgRegisterReferral)(nil)
)

func NewMsgRegisterReferral(referral, referrer sdk.AccAddress) *MsgRegisterReferral {
	return &MsgRegisterReferral{
		ReferralAddress: referral.String(),
		ReferrerAddress: referrer.String(),
	}
}

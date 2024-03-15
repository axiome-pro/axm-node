package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

// Query endpoints supported by the referral querier
const (
	QueryStatus             = "status"
	QueryReferrer           = "referrer"
	QueryDelegatedInNetwork = "delegated"
	QueryCheckStatus        = "check-status"
	QueryParams             = "params"
	QueryInfo               = "info"
)

type QueryResChildren []sdk.AccAddress

func (qr QueryResChildren) String() string {
	strs := make([]string, len(qr), len(qr))
	for i, adr := range qr {
		strs[i] = adr.String()
	}
	return strings.Join(strs[:], ", ")
}

type QueryResValidateTransition struct {
	Ok  bool   `json:"ok" yaml:"ok"`
	Err string `json:"err,omitempty" yaml:"err,omitempty"`
}

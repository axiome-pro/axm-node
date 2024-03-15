package types

import (
	"time"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewDowngrade(acc string, current Status, at time.Time) *Downgrade {
	return &Downgrade{
		Account: acc,
		Current: current,
		Time:    at,
	}
}

func (d Downgrade) GetAccount() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(d.Account)
	if err != nil {
		panic(err)
	}
	return addr
}

func NewRefs(parent string, children []*RefInfo) *Refs {
	return &Refs{
		Referrer:  parent,
		Referrals: children,
	}
}

func NewRefInfo(addr string, status Status) *RefInfo {
	return &RefInfo{
		Address: addr,
		Status:  status,
	}
}

// NewGenesisState creates a new GenesisState object
func NewGenesisState(
	params Params,
	topLevelAccounts []*RefInfo,
	otherAccounts []Refs,
	downgrades []Downgrade,
) *GenesisState {
	return &GenesisState{
		Params:           params,
		TopLevelAccounts: topLevelAccounts,
		OtherAccounts:    otherAccounts,
		Downgrades:       downgrades,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis validates the referral genesis parameters
func ValidateGenesis(data GenesisState) error {
	if data.TopLevelAccounts == nil {
		return errors.New("empty top level accounts set")
	}
	if err := data.Params.Validate(); err != nil {
		return errors.Wrap(err, "invalid params")
	}
	return nil
}

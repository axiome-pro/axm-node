package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/axiome-pro/axm-node/util"
	referral "github.com/axiome-pro/axm-node/x/referral/types"
)

var (
	_ codectypes.UnpackInterfacesMessage = Proposal{}
)

func (p Proposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	return sdktx.UnpackInterfaces(unpacker, p.Messages)
}

func (p Proposal) GetAuthor() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(p.Author)
	if err != nil {
		panic(err)
	}
	return addr
}

func (p Proposal) String() string {
	bz, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

func (p Proposal) Validate() error {
	if p.Name == "" {
		return errors.New("invalid name: empty string")
	}
	if _, err := sdk.AccAddressFromBech32(p.Author); err != nil {
		return errors.Wrap(err, "invalid author")
	}
	return nil
}

func (g Government) GetMembers() []sdk.AccAddress {
	addrz := make([]sdk.AccAddress, len(g.Members))
	for i, bech32 := range g.Members {
		addr, err := sdk.AccAddressFromBech32(bech32)
		if err != nil {
			panic(err)
		}
		addrz[i] = addr
	}
	return addrz
}

func (g Government) GetMember(i int) sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(g.Members[i])
	if err != nil {
		panic(err)
	}
	return addr
}

func (g Government) String() string {
	bz, err := yaml.Marshal(g.Members)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

func (g Government) Strings() []string { return g.Members }

func (g Government) Contains(addr sdk.AccAddress) bool {
	bech32 := addr.String()
	for _, elem := range g.Members {
		if elem == bech32 {
			return true
		}
	}

	return false
}

func (g *Government) Remove(addr sdk.AccAddress) {
	bech32 := addr.String()
	for index, elem := range g.Members {
		if elem == bech32 {
			g.Members = append(g.Members[:index], g.Members[index+1:]...)
			return
		}
	}
}

func (g *Government) Append(addr sdk.AccAddress) {
	g.Members = append(g.Members, addr.String())
}

func (r ProposalHistoryRecord) GetGovernment() *Government {
	return &Government{Members: r.Government}
}

func (r ProposalHistoryRecord) GetAgreed() *Government {
	return &Government{Members: r.Agreed}
}

func (r ProposalHistoryRecord) GetDisagreed() *Government {
	return &Government{Members: r.Disagreed}
}

func (r ProposalHistoryRecord) Validate() error {
	if err := r.Proposal.Validate(); err != nil {
		return errors.Wrap(err, "invalid proposal")
	}
	if r.Government == nil {
		return errors.New("invalid government: empty list")
	}
	for i, bech32 := range r.Government {
		if _, err := sdk.AccAddressFromBech32(bech32); err != nil {
			return errors.Wrapf(err, "invalid government (item #%d)", i)
		}
	}
	for i, bech32 := range r.Agreed {
		if _, err := sdk.AccAddressFromBech32(bech32); err != nil {
			return errors.Wrapf(err, "invalid agreed (item #%d)", i)
		}
	}
	for i, bech32 := range r.Disagreed {
		if _, err := sdk.AccAddressFromBech32(bech32); err != nil {
			return errors.Wrapf(err, "invalid disagreed (item #%d)", i)
		}
	}
	if r.Started <= 0 {
		return errors.New("invalid started: must be positive")
	}
	if r.Finished <= 0 {
		return errors.New("invalid finished: must be positive")
	}
	return nil
}

func NewPollValidators(author sdk.AccAddress, name, text string, quorum util.Fraction) Poll {
	return Poll{
		Name:         name,
		Author:       author.String(),
		Question:     text,
		Quorum:       &quorum,
		Requirements: &Poll_CanValidate{CanValidate: &Poll_Unit{}},
	}
}

func NewPollStatus(author sdk.AccAddress, name, text string, quorum util.Fraction, status referral.Status) Poll {
	return Poll{
		Name:         name,
		Author:       author.String(),
		Question:     text,
		Quorum:       &quorum,
		Requirements: &Poll_MinStatus{MinStatus: status},
	}
}

func (p Poll) String() string {
	bz, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

func (p Poll) Validate() error {
	if len(p.Name)+len(p.Question) == 0 {
		return errors.New("neither name nor question specified")
	}
	if _, err := sdk.AccAddressFromBech32(p.Author); err != nil {
		return errors.Wrap(err, "cannot parse author")
	}
	if p.Quorum != nil && (p.Quorum.IsNegative() || p.Quorum.GT(util.FractionInt(1))) {
		return errors.New("quorum must be nil or in [0; 1]")
	}
	switch r := p.Requirements.(type) {
	case *Poll_CanValidate:
		// pass
	case *Poll_MinStatus:
		if r.MinStatus < referral.MinimumStatus || r.MinStatus > referral.MaximumStatus {
			return errors.New("min_status is out of range")
		}
	}
	if p.StartTime != nil && p.EndTime != nil && !p.EndTime.After(*p.StartTime) {
		return errors.New("start_time after end_time")
	}
	return nil
}

func (u *Poll_Unit) Equal(other *Poll_Unit) bool { return (u == nil) == (other == nil) }

package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/axiome-pro/axm-node/util"
)

func (s Status) LinesOpened() int {
	switch s {
	case STATUS_NEW:
		return 0
	case STATUS_STARTER:
		return 2
	case STATUS_LEADER:
		return 4
	case STATUS_GURU:
		return 6
	case STATUS_BOSS:
		return 8
	case STATUS_PRO:
		return 10
	case STATUS_TOP:
		return 12
	case STATUS_MEGA:
		return 14
	default:
		return 14
	}
}

const MinimumStatus = STATUS_NEW
const MaximumStatus = STATUS_MEGA

func NewInfo(referrer string, delegated math.Int) Info {
	zero := math.ZeroInt()
	return Info{
		Status:          STATUS_NEW,
		Referrer:        referrer,
		Active:          false,
		SelfDelegated:   &delegated,
		TeamDelegated:   &zero,
		ActiveCount:     &ActiveAggregations{},
		ActiveRefCounts: make([]uint64, 15),
	}
}

func NewInfoWithStatus(referrer string, delegated math.Int, status Status) Info {
	zero := math.ZeroInt()
	return Info{
		Status:          status,
		Referrer:        referrer,
		Active:          false,
		SelfDelegated:   &delegated,
		TeamDelegated:   &zero,
		ActiveCount:     &ActiveAggregations{},
		ActiveRefCounts: make([]uint64, 15),
	}
}

func (r Info) RegistrationClosed(ctx sdk.Context) bool {
	return false
}

func (r Info) GetReferrer() sdk.AccAddress {
	if r.Referrer == "" {
		return nil
	}
	addr, err := sdk.AccAddressFromBech32(r.Referrer)
	if err != nil {
		panic(err)
	}
	return addr
}

func (r *Info) Normalize() {
	for len(r.ActiveRefCounts) < 15 {
		r.ActiveRefCounts = append(r.ActiveRefCounts, uint64(0))
	}
}

func (r Info) IsEmpty() bool {
	return r.Status == STATUS_UNSPECIFIED
}

func (r Info) GetActiveRefsCountFromLevelToLevel(from, to int) (sum uint64) {
	for i := from; i <= to; i++ {
		sum += r.ActiveRefCounts[i]
	}
	return sum
}

type ReferralFee struct {
	Beneficiary string        `json:"beneficiary" yaml:"beneficiary"`
	Ratio       util.Fraction `json:"ratio" yaml:"ratio"`
}

func (fee ReferralFee) GetBeneficiary() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(fee.Beneficiary)
	if err != nil {
		panic(err)
	}
	return addr
}

type ReferralValidatorFee struct {
	Beneficiary string        `json:"beneficiary" yaml:"beneficiary"`
	Ratio       util.Fraction `json:"ratio" yaml:"ratio"`
}

func (fee ReferralValidatorFee) GetBeneficiary() sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(fee.Beneficiary)
	if err != nil {
		panic(err)
	}
	return addr
}

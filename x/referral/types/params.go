package types

import (
	"fmt"
	"gopkg.in/yaml.v3"

	"github.com/axiome-pro/axm-node/util"
)

var (
	DefaultDelegatingAward = NetworkAward{
		Network: []util.Fraction{
			util.Percent(5),  // 1
			util.Percent(1),  // 2
			util.Percent(2),  // 3
			util.Percent(1),  // 4
			util.Percent(1),  // 5
			util.Percent(1),  // 6
			util.Percent(1),  // 7
			util.Percent(1),  // 8
			util.Permille(7), // 9
			util.Permille(5), // 10
			util.Permille(3), // 11
			util.Permille(2), // 12
			util.Permille(2), // 13
			util.Permille(1), // 14
		},
	}

	// DefaultStatusDowngradePeriod 7 days
	DefaultStatusDowngradePeriod int32 = 7 * 24 * 60 * 60
)

func (na NetworkAward) Validate() error { return validateNetworkAward(na) }

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

func (p Params) Validate() error {
	if err := validateNetworkAward(p.DelegatingAward); err != nil {
		return nil
	}
	return nil
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return Params{
		DelegatingAward:       DefaultDelegatingAward,
		StatusDowngradePeriod: DefaultStatusDowngradePeriod,
	}
}

func validateNetworkAward(i interface{}) error {
	na, ok := i.(NetworkAward)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	total := util.NewFraction(0, 1)
	for i := 0; i < 14; i++ {
		if na.Network[i].IsNegative() {
			return fmt.Errorf("level %d award must be non-negative", i+1)
		}
		total = total.Add(na.Network[i])
	}
	if total.GTE(util.Percent(100)) {
		return fmt.Errorf("total network award must be less than 100%%")
	}
	return nil
}

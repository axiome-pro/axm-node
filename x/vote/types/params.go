package types

import (
	"fmt"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

// Parameter store keys
var (
	DefaultvotePeriod int32 = 24 * 60
)

// NewParams creates a new Params object
func NewParams(votePeriod, pollPeriod int32) Params {
	return Params{
		VotePeriod: votePeriod,
		PollPeriod: pollPeriod,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	bz, err := yaml.Marshal(p)
	if err != nil {
		panic(err)
	}
	return string(bz)
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultvotePeriod, DefaultvotePeriod)
}

func (p Params) Validate() error {
	if err := validatevotePeriod(p.VotePeriod); err != nil {
		return errors.Wrap(err, "invalid vote_period")
	}
	if err := validatevotePeriod(p.PollPeriod); err != nil {
		return errors.Wrap(err, "invalid poll_period")
	}
	return nil
}

func validatevotePeriod(i interface{}) error {
	v, ok := i.(int32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v < 1 {
		return fmt.Errorf("validating period must be at least 1 minute: %d", v)
	}

	return nil
}

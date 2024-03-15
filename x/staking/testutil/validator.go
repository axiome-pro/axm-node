package testutil

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/axiome-pro/axm-node/x/staking/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewValidator is a testing helper method to create validators in tests
func NewValidator(tb testing.TB, operator sdk.ValAddress, pubKey cryptotypes.PubKey) types.Validator {
	tb.Helper()
	v, err := types.NewValidator(operator.String(), pubKey, types.Description{})
	require.NoError(tb, err)
	return v
}

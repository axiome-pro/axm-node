package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// validator struct to define the fields of the validator
type validator struct {
	Amount   sdk.Coin
	PubKey   cryptotypes.PubKey
	Moniker  string
	Identity string
	Website  string
	Security string
	Details  string
}

func parseAndValidateValidatorJSON(cdc codec.Codec, path string) (validator, error) {
	type internalVal struct {
		Amount   string          `json:"amount"`
		PubKey   json.RawMessage `json:"pubkey"`
		Moniker  string          `json:"moniker"`
		Identity string          `json:"identity,omitempty"`
		Website  string          `json:"website,omitempty"`
		Security string          `json:"security,omitempty"`
		Details  string          `json:"details,omitempty"`
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		return validator{}, err
	}

	var v internalVal
	err = json.Unmarshal(contents, &v)
	if err != nil {
		return validator{}, err
	}

	if v.Amount == "" {
		return validator{}, fmt.Errorf("must specify amount of coins to bond")
	}
	amount, err := sdk.ParseCoinNormalized(v.Amount)
	if err != nil {
		return validator{}, err
	}

	if v.PubKey == nil {
		return validator{}, fmt.Errorf("must specify the JSON encoded pubkey")
	}
	var pk cryptotypes.PubKey
	if err := cdc.UnmarshalInterfaceJSON(v.PubKey, &pk); err != nil {
		return validator{}, err
	}

	if v.Moniker == "" {
		return validator{}, fmt.Errorf("must specify the moniker name")
	}

	return validator{
		Amount:   amount,
		PubKey:   pk,
		Moniker:  v.Moniker,
		Identity: v.Identity,
		Website:  v.Website,
		Security: v.Security,
		Details:  v.Details,
	}, nil
}

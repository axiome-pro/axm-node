package cli

import (
	"cosmossdk.io/core/address"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"strings"

	"github.com/axiome-pro/axm-node/x/referral/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// NewTxCmd returns the transaction commands for this module
func NewTxCmd(ac address.Codec) *cobra.Command {
	referralTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Aliases:                    []string{"ref", "r"},
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	referralTxCmd.AddCommand(
		NewRegisterReferralCmd(ac),
		CmdBurnPartnerFee(),
	)

	return referralTxCmd
}

func CmdBurnPartnerFee() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "burn-partner-fee <amount>",
		Short: "Burn some uAXMs",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Burn some coins from account balance

Example:
$ %s tx referral burn-partner-fee 1000000uaxm --from mykey
`,
				version.AppName,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			coins, err := sdk.ParseCoinsNormalized(args[0])
			if err != nil {
				return err
			}

			if len(coins) == 0 {
				return fmt.Errorf("invalid coins")
			}

			msg := types.MsgBurnPartnerFee{
				FromAddress: clientCtx.GetFromAddress().String(),
				Amount:      coins,
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// NewRegisterReferralCmd returns a CLI command handler for creating a MsgRegisterReferral transaction.
func NewRegisterReferralCmd(ac address.Codec) *cobra.Command {
	bech32PrefixAccAddr := sdk.GetConfig().GetBech32AccountAddrPrefix()

	cmd := &cobra.Command{
		Use:   "register-referral [referrer-address]",
		Short: "register referral account in referral marketing module",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Register referral in referral marketing module

Example:
$ %s tx referral register-referral %s1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey
`,
				version.AppName, bech32PrefixAccAddr,
			),
		),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			referral := clientCtx.GetFromAddress()
			referrer, err := ac.StringToBytes(args[0])
			if err != nil {
				return err
			}

			msg := types.NewMsgRegisterReferral(referral, referrer)

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

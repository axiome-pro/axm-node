package cli

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"regexp"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"

	"github.com/axiome-pro/axm-node/util"
	referral "github.com/axiome-pro/axm-node/x/referral/types"
	"github.com/axiome-pro/axm-node/x/vote/types"
)

// GetTxCmd returns the transaction commands for this module
func NewTxCmd() *cobra.Command {
	voteTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	voteTxCmd.AddCommand(
		NewCmdSubmitProposal(),
		cmdVote(),
		cmdStartPoll(),
		cmdAnswerPoll(),
	)

	return voteTxCmd
}

func cmdVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote agree|disagree <voter_key_or_address>",
		Short: "Vote for/against the current proposal",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Set(flags.FlagFrom, args[1]); err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			voter := clientCtx.GetFromAddress().String()
			agree := strings.ToLower(args[0]) == "agree"
			if !agree && strings.ToLower(args[0]) != "disagree" {
				return errors.New("cannot parse aggree/disagree flag")
			}

			msg := &types.MsgVote{
				Voter: voter,
				Agree: agree,
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdStartPoll() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "start-poll <author_key_or_address> validators|status:<status> <name> <text> [quorum]",
		Aliases: []string{"start_poll", "sp"},
		Short:   "Start a public poll",
		Example: `start-poll ivan validators Halving "Should we decrease all awards by a half next Monday?" 2/3`,
		Args:    cobra.RangeArgs(4, 5),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Set(flags.FlagFrom, args[0]); err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			poll := types.Poll{
				Name:     args[2],
				Author:   clientCtx.GetFromAddress().String(),
				Question: args[3],
			}

			if req := args[1]; req == "validators" {
				poll.Requirements = &types.Poll_CanValidate{CanValidate: &types.Poll_Unit{}}
			} else if m := regexp.MustCompile(`/^status:(\d+)|([A-Za-z_]+)$/`).FindStringSubmatch(req); m != nil {
				var status referral.Status
				if len(m[1]) > 0 {
					if s, err := strconv.Atoi(m[1]); err != nil {
						return errors.Wrap(err, "cannot parse status")
					} else {
						status = referral.Status(s)
					}
				} else {
					name := strings.ToUpper(m[2])
					if !strings.HasPrefix(name, "STATUS_") {
						name = "STATUS_" + name
					}
					if s, ok := referral.Status_value[name]; !ok {
						return errors.New("cannot parse status")
					} else {
						status = referral.Status(s)
					}
				}

				poll.Requirements = &types.Poll_MinStatus{MinStatus: status}
			} else {
				return errors.New("cannot parse requirements")
			}

			if len(args) > 4 {
				if q, err := util.ParseFraction(args[4]); err != nil {
					return errors.Wrap(err, "cannot parse quorum")
				} else {
					poll.Quorum = &q
				}
			}

			msg := types.MsgStartPoll{Poll: poll}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func cmdAnswerPoll() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "answer yes|no <respondent_key_or_address>",
		Aliases: []string{"ans", "a", "answer-poll", "answer_poll"},
		Short:   "Answer the current public poll",
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmd.Flags().Set(flags.FlagFrom, args[1]); err != nil {
				return err
			}
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var yes bool
			if ans := strings.ToLower(args[0]); ans == "yes" {
				yes = true
			} else if ans != "no" {
				return errors.New("cannot parse answer")
			}

			msg := types.MsgAnswerPoll{
				Respondent: clientCtx.GetFromAddress().String(),
				Yes:        yes,
			}
			if err = msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

func NewCmdSubmitProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submit-proposal [path/to/proposal.json]",
		Short: "Submit a proposal along with some messages, metadata and deposit",
		Args:  cobra.ExactArgs(1),
		Long: strings.TrimSpace(
			fmt.Sprintf(`Submit a proposal along with some messages, metadata and deposit.
They should be defined in a JSON file.

Example:
$ %s tx gov submit-proposal path/to/proposal.json

Where proposal.json contains:

{
  // array of proto-JSON-encoded sdk.Msgs
  "messages": [
    {
      "@type": "/cosmos.bank.v1beta1.MsgSend",
      "from_address": "cosmos1...",
      "to_address": "cosmos1...",
      "amount":[{"denom": "stake","amount": "10"}]
    }
  ],
  // metadata can be any of base64 encoded, raw text, stringified json, IPFS link to json
  // see below for example metadata
  "name": "Send coins proposal"
}

`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			proposal, msgs, err := parseSubmitProposal(clientCtx.Codec, args[0])
			if err != nil {
				return err
			}

			msg, err := types.NewMsgPropose(msgs, clientCtx.GetFromAddress().String(), proposal.Name)
			if err != nil {
				return fmt.Errorf("invalid message: %w", err)
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

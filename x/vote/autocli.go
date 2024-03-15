package vote

import (
	referralv1beta1 "github.com/axiome-pro/axm-node/api/axiome/vote/v1beta1"
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: referralv1beta1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Query the current vote parameters.",
				},
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service:           referralv1beta1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				//{
				//	RpcMethod: "SetWithdrawAddress",
				//	Use:       "set-withdraw-addr [withdraw-addr]",
				//	Short:     "Change the default withdraw address for rewards associated with an address",
				//	Example:   fmt.Sprintf("%s tx distribution set-withdraw-addr cosmos1gghjut3ccd8ay0zduzj64hwre2fxs9ld75ru9p --from mykey", version.AppName),
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{
				//		{ProtoField: "withdraw_address"},
				//	},
				//},
				//{
				//	RpcMethod: "WithdrawDelegatorReward",
				//	Use:       "withdraw-rewards [validator-addr]",
				//	Short:     "Withdraw rewards from a given delegation address",
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{
				//		{ProtoField: "validator_address"},
				//	},
				//},
				//{
				//	RpcMethod: "WithdrawValidatorCommission",
				//	Use:       "withdraw-validator-commission [validator-addr]",
				//	Short:     "Withdraw commissions from a validator address (must be a validator operator)",
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{
				//		{ProtoField: "validator_address"},
				//	},
				//},
				//{
				//	RpcMethod: "DepositValidatorRewardsPool",
				//	Use:       "fund-validator-rewards-pool [validator-addr] [amount]",
				//	Short:     "Fund the validator rewards pool with the specified amount",
				//	Example:   fmt.Sprintf("%s tx distribution fund-validator-rewards-pool cosmosvaloper1x20lytyf6zkcrv5edpkfkn8sz578qg5sqfyqnp 100uatom --from mykey", version.AppName),
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{
				//		{ProtoField: "validator_address"},
				//		{ProtoField: "amount", Varargs: true},
				//	},
				//},
				//{
				//	RpcMethod: "FundCommunityPool",
				//	Use:       "fund-community-pool [amount]",
				//	Short:     "Funds the community pool with the specified amount",
				//	Example:   fmt.Sprintf(`$ %s tx distribution fund-community-pool 100uatom --from mykey`, version.AppName),
				//	PositionalArgs: []*autocliv1.PositionalArgDescriptor{
				//		{ProtoField: "amount", Varargs: true},
				//	},
				//},
				//{
				//	RpcMethod: "UpdateParams",
				//	Skip:      true, // skipped because authority gated
				//},
				//{
				//	RpcMethod: "CommunityPoolSpend",
				//	Skip:      true, // skipped because authority gated
				//},
			},
			EnhanceCustomCommand: false, // use custom commands only until v0.51
		},
	}
}

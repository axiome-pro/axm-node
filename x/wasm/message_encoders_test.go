package wasm

import (
	"bytes"
	"testing"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	axmdisttypes "github.com/axiome-pro/axm-node/x/distribution/types"
	axmstakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func TestAxmStakingEncoderUsesLocalMsgTypes(t *testing.T) {
	sender := sdk.AccAddress(bytes.Repeat([]byte{0x1}, 20))
	encoders := axmMessageEncoders()

	tests := []struct {
		name          string
		msg           *wasmvmtypes.StakingMsg
		expectedType  sdk.Msg
		expectedCheck func(t *testing.T, msg sdk.Msg)
	}{
		{
			name: "delegate",
			msg: &wasmvmtypes.StakingMsg{
				Delegate: &wasmvmtypes.DelegateMsg{
					Validator: "axmvaloper1delegate",
					Amount:    wasmvmtypes.Coin{Denom: "uaxm", Amount: "42"},
				},
			},
			expectedType: &axmstakingtypes.MsgDelegate{},
			expectedCheck: func(t *testing.T, msg sdk.Msg) {
				delegateMsg := msg.(*axmstakingtypes.MsgDelegate)
				require.Equal(t, sender.String(), delegateMsg.DelegatorAddress)
				require.Equal(t, "axmvaloper1delegate", delegateMsg.ValidatorAddress)
				require.Equal(t, "uaxm", delegateMsg.Amount.Denom)
				require.Equal(t, "42", delegateMsg.Amount.Amount.String())
			},
		},
		{
			name: "undelegate",
			msg: &wasmvmtypes.StakingMsg{
				Undelegate: &wasmvmtypes.UndelegateMsg{
					Validator: "axmvaloper1undelegate",
					Amount:    wasmvmtypes.Coin{Denom: "uaxm", Amount: "7"},
				},
			},
			expectedType: &axmstakingtypes.MsgUndelegate{},
			expectedCheck: func(t *testing.T, msg sdk.Msg) {
				undelegateMsg := msg.(*axmstakingtypes.MsgUndelegate)
				require.Equal(t, sender.String(), undelegateMsg.DelegatorAddress)
				require.Equal(t, "axmvaloper1undelegate", undelegateMsg.ValidatorAddress)
				require.Equal(t, "uaxm", undelegateMsg.Amount.Denom)
				require.Equal(t, "7", undelegateMsg.Amount.Amount.String())
			},
		},
		{
			name: "redelegate",
			msg: &wasmvmtypes.StakingMsg{
				Redelegate: &wasmvmtypes.RedelegateMsg{
					SrcValidator: "axmvaloper1src",
					DstValidator: "axmvaloper1dst",
					Amount:       wasmvmtypes.Coin{Denom: "uaxm", Amount: "9"},
				},
			},
			expectedType: &axmstakingtypes.MsgBeginRedelegate{},
			expectedCheck: func(t *testing.T, msg sdk.Msg) {
				redelegateMsg := msg.(*axmstakingtypes.MsgBeginRedelegate)
				require.Equal(t, sender.String(), redelegateMsg.DelegatorAddress)
				require.Equal(t, "axmvaloper1src", redelegateMsg.ValidatorSrcAddress)
				require.Equal(t, "axmvaloper1dst", redelegateMsg.ValidatorDstAddress)
				require.Equal(t, "uaxm", redelegateMsg.Amount.Denom)
				require.Equal(t, "9", redelegateMsg.Amount.Amount.String())
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := encoders.Staking(sender, tc.msg)
			require.NoError(t, err)
			require.Len(t, got, 1)
			require.IsType(t, tc.expectedType, got[0])
			tc.expectedCheck(t, got[0])

			defaultMsgs, err := wasmkeeper.EncodeStakingMsg(sender, tc.msg)
			require.NoError(t, err)
			require.Len(t, defaultMsgs, 1)
			require.NotEqual(t, sdk.MsgTypeURL(defaultMsgs[0]), sdk.MsgTypeURL(got[0]))
		})
	}

	require.Equal(t, sdk.MsgTypeURL(&stakingtypes.MsgDelegate{}), "/cosmos.staking.v1beta1.MsgDelegate")
	require.Equal(t, sdk.MsgTypeURL(&axmstakingtypes.MsgDelegate{}), "/axiome.staking.v1beta1.MsgDelegate")
}

func TestAxmDistributionEncoderUsesLocalMsgTypes(t *testing.T) {
	sender := sdk.AccAddress(bytes.Repeat([]byte{0x2}, 20))
	encoders := axmMessageEncoders()

	msg := &wasmvmtypes.DistributionMsg{
		WithdrawDelegatorReward: &wasmvmtypes.WithdrawDelegatorRewardMsg{Validator: "axmvaloper1reward"},
	}

	got, err := encoders.Distribution(sender, msg)
	require.NoError(t, err)
	require.Len(t, got, 1)
	require.IsType(t, &axmdisttypes.MsgWithdrawDelegatorReward{}, got[0])

	withdrawMsg := got[0].(*axmdisttypes.MsgWithdrawDelegatorReward)
	require.Equal(t, sender.String(), withdrawMsg.DelegatorAddress)
	require.Equal(t, "axmvaloper1reward", withdrawMsg.ValidatorAddress)

	defaultMsgs, err := wasmkeeper.EncodeDistributionMsg(sender, msg)
	require.NoError(t, err)
	require.Len(t, defaultMsgs, 1)
	require.NotEqual(t, sdk.MsgTypeURL(defaultMsgs[0]), sdk.MsgTypeURL(got[0]))

	require.Equal(t, sdk.MsgTypeURL(&distributiontypes.MsgWithdrawDelegatorReward{}), "/cosmos.distribution.v1beta1.MsgWithdrawDelegatorReward")
	require.Equal(t, sdk.MsgTypeURL(&axmdisttypes.MsgWithdrawDelegatorReward{}), "/axiome.distribution.v1beta1.MsgWithdrawDelegatorReward")
}

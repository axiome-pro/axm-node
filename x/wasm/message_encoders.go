package wasm

import (
	"encoding/json"

	errorsmod "cosmossdk.io/errors"
	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	axmdisttypes "github.com/axiome-pro/axm-node/x/distribution/types"
	axmstakingtypes "github.com/axiome-pro/axm-node/x/staking/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// axmMessageEncoders overrides the staking messages and
// DistributionMsg::WithdrawDelegatorReward so CosmWasm routes into this chain's
// custom axiome protobuf types where needed.
func axmMessageEncoders() *wasmkeeper.MessageEncoders {
	return &wasmkeeper.MessageEncoders{
		Distribution: encodeAxmDistributionMsg,
		Staking:      encodeAxmStakingMsg,
		Custom:       unsupportedCustomMsg,
	}
}

func unsupportedCustomMsg(_ sdk.AccAddress, _ json.RawMessage) ([]sdk.Msg, error) {
	return nil, errorsmod.Wrap(wasmtypes.ErrUnknownMsg, "custom variant not supported")
}

func singleSDKMsg(msg sdk.Msg) []sdk.Msg {
	return []sdk.Msg{msg}
}

func unknownVariantError(kind string) error {
	return errorsmod.Wrapf(wasmtypes.ErrUnknownMsg, "unknown variant of %s", kind)
}

func encodeWithSingleCoin(amount wasmvmtypes.Coin, build func(sdk.Coin) sdk.Msg) ([]sdk.Msg, error) {
	coin, err := wasmkeeper.ConvertWasmCoinToSdkCoin(amount)
	if err != nil {
		return nil, err
	}

	return singleSDKMsg(build(coin)), nil
}

func encodeAxmDistributionMsg(sender sdk.AccAddress, msg *wasmvmtypes.DistributionMsg) ([]sdk.Msg, error) {
	switch {
	case msg.WithdrawDelegatorReward != nil:
		return singleSDKMsg(&axmdisttypes.MsgWithdrawDelegatorReward{
			DelegatorAddress: sender.String(),
			ValidatorAddress: msg.WithdrawDelegatorReward.Validator,
		}), nil

	default:
		return wasmkeeper.EncodeDistributionMsg(sender, msg)
	}
}

func encodeAxmStakingMsg(sender sdk.AccAddress, msg *wasmvmtypes.StakingMsg) ([]sdk.Msg, error) {
	senderAddress := sender.String()

	switch {
	case msg.Delegate != nil:
		return encodeWithSingleCoin(msg.Delegate.Amount, func(coin sdk.Coin) sdk.Msg {
			return &axmstakingtypes.MsgDelegate{
				DelegatorAddress: senderAddress,
				ValidatorAddress: msg.Delegate.Validator,
				Amount:           coin,
			}
		})

	case msg.Redelegate != nil:
		return encodeWithSingleCoin(msg.Redelegate.Amount, func(coin sdk.Coin) sdk.Msg {
			return &axmstakingtypes.MsgBeginRedelegate{
				DelegatorAddress:    senderAddress,
				ValidatorSrcAddress: msg.Redelegate.SrcValidator,
				ValidatorDstAddress: msg.Redelegate.DstValidator,
				Amount:              coin,
			}
		})

	case msg.Undelegate != nil:
		return encodeWithSingleCoin(msg.Undelegate.Amount, func(coin sdk.Coin) sdk.Msg {
			return &axmstakingtypes.MsgUndelegate{
				DelegatorAddress: senderAddress,
				ValidatorAddress: msg.Undelegate.Validator,
				Amount:           coin,
			}
		})

	default:
		return nil, unknownVariantError("Staking")
	}
}

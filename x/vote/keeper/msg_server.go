package keeper

import (
	"context"

	"github.com/axiome-pro/axm-node/util"
	"github.com/axiome-pro/axm-node/x/vote/types"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type MsgServer Keeper

var _ types.MsgServer = MsgServer{}

func (ms MsgServer) Propose(ctx context.Context, msg *types.MsgPropose) (*types.MsgProposeResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = Keeper(ms)
	)
	k.Logger(sdkCtx).Info("Message type", "type", msg.Type(), "name", msg.Name)

	_, err := sdktx.GetMsgs(msg.Messages, "sdk.Msg")
	if err != nil {
		k.Logger(sdkCtx).Info("Error", "err", err)
		return nil, err
	}

	if err := k.Propose(sdkCtx, *msg); err != nil {
		return nil, err
	}
	util.TagTx(sdkCtx, types.ModuleName, msg)
	return &types.MsgProposeResponse{}, nil
}

func (ms MsgServer) Vote(ctx context.Context, msg *types.MsgVote) (*types.MsgVoteResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = Keeper(ms)
	)
	if err := k.Vote(sdkCtx, msg.GetVoter(), msg.Agree); err != nil {
		return nil, err
	}
	util.TagTx(sdkCtx, types.ModuleName, msg)
	return &types.MsgVoteResponse{}, nil
}

func (ms MsgServer) StartPoll(ctx context.Context, msg *types.MsgStartPoll) (*types.MsgStartPollResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = Keeper(ms)
	)
	if err := k.StartPoll(sdkCtx, msg.Poll); err != nil {
		return nil, err
	}
	util.TagTx(sdkCtx, types.ModuleName, msg)
	return &types.MsgStartPollResponse{}, nil
}

func (ms MsgServer) AnswerPoll(ctx context.Context, msg *types.MsgAnswerPoll) (*types.MsgAnswerPollResponse, error) {
	var (
		sdkCtx = sdk.UnwrapSDKContext(ctx)
		k      = Keeper(ms)
	)
	if err := k.Answer(sdkCtx, msg.Respondent, msg.Yes); err != nil {
		return nil, err
	}
	util.TagTx(sdkCtx, types.ModuleName, msg)
	return &types.MsgAnswerPollResponse{}, nil
}

func (ms MsgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	var k = Keeper(ms)

	if err := k.validateAuthority(msg.Authority); err != nil {
		return nil, err
	}

	if err := msg.Params.Validate(); err != nil {
		return nil, err
	}

	if err := k.Params.Set(ctx, msg.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

func (ms MsgServer) changeGovernment(ctx context.Context, authority, governor string, add bool) error {
	var (
		k      = Keeper(ms)
		sdkCtx = sdk.UnwrapSDKContext(ctx)
	)

	if err := k.validateAuthority(authority); err != nil {
		return err
	}

	governorAddr, err := k.accountKeeper.AddressCodec().StringToBytes(governor)
	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid governor address: %s", err)
	}

	if add {
		err = k.AddGovernor(sdkCtx, governorAddr)
	} else {
		err = k.RemoveGovernor(sdkCtx, governorAddr)
	}

	return err
}

func (ms MsgServer) AddGovernor(ctx context.Context, msg *types.MsgAddGovernor) (*types.MsgAddGovernorResponse, error) {
	if err := ms.changeGovernment(ctx, msg.Authority, msg.Governor, true); err != nil {
		return nil, err
	}
	return &types.MsgAddGovernorResponse{}, nil
}

func (ms MsgServer) RemoveGovernor(ctx context.Context, msg *types.MsgRemoveGovernor) (*types.MsgRemoveGovernorResponse, error) {
	if err := ms.changeGovernment(ctx, msg.Authority, msg.Governor, false); err != nil {
		return nil, err
	}

	return &types.MsgRemoveGovernorResponse{}, nil
}
func (k *Keeper) validateAuthority(authority string) error {
	_, err := k.accountKeeper.AddressCodec().StringToBytes(authority)

	if err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
	}

	if sdk.AccAddress.String(k.authority) != authority {
		return errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, authority)
	}

	return nil
}

package keeper

import (
	"context"
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"

	"github.com/axiome-pro/axm-node/x/referral/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	Keeper
}

func NewMsgServer(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) RegisterReferral(ctx context.Context, msg *types.MsgRegisterReferral) (*types.MsgRegisterReferralResponse, error) {
	_, err := k.accountKeeper.AddressCodec().StringToBytes(msg.ReferrerAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid referrer address: %s", err)
	}

	_, err = k.accountKeeper.AddressCodec().StringToBytes(msg.ReferralAddress)
	if err != nil {
		return nil, sdkerrors.ErrInvalidAddress.Wrapf("invalid referrer address: %s", err)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	if !k.exists(sdkCtx, msg.ReferrerAddress) {
		return nil, sdkerrors.ErrNotFound.Wrapf("parent not found: %s", msg.ReferralAddress)
	}

	err = k.AppendChild(sdkCtx, msg.ReferrerAddress, msg.ReferralAddress)

	if err != nil {
		return nil, errors.Wrap(err, "unable to append new child")
	}

	return &types.MsgRegisterReferralResponse{}, nil
}

func (k msgServer) UpdateParams(ctx context.Context, msg *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
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

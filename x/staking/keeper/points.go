package keeper

import (
	"context"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) GetEmissionRatioFromBondedRatio(ctx context.Context, bondedRatio math.LegacyDec) (math.LegacyDec, error) {
	emissionTable, err := k.EmissionTable(ctx)
	if err != nil {
		return math.LegacyDec{}, err
	}

	start := math.LegacyZeroDec()
	rate := math.LegacyZeroDec()

	for _, emissionRange := range emissionTable {
		if emissionRange.Start.GTE(start) && emissionRange.Start.LTE(bondedRatio) {
			rate = emissionRange.Rate
		}
	}

	return rate, nil
}

func (k *Keeper) AllocateValidatorsPoints(ctx context.Context) error {
	logger := k.Logger(ctx)

	sdkCtx := sdk.UnwrapSDKContext(ctx)

	histInfo, err := k.GetHistoricalInfo(ctx, sdkCtx.BlockHeight()-1)

	if err != nil {
		logger.Info("No historical entry, emission points not calculated", "height", sdkCtx.BlockHeight()-1)
		return nil
	}

	ratio, err := k.StakedRatio(ctx)
	if err != nil {
		return err
	}

	timeElapsed := (sdkCtx.BlockTime().UnixMicro() - histInfo.Header.Time.UnixMicro())
	emissionRate, err := k.GetEmissionRatioFromBondedRatio(ctx, ratio)
	pointsToAdd := emissionRate.MulInt64(timeElapsed).TruncateInt().Uint64()
	pointsToAddInt := math.NewInt(int64(pointsToAdd))

	logger.Debug("Point routine", "bondedRatio", ratio, "emissionRate", emissionRate, "pointToAdd", pointsToAdd)

	monthlyPoints, err := k.MaximumMonthlyPoints(ctx)
	if err != nil {
		return err
	}

	monthlyPointsInt := math.NewInt(int64(monthlyPoints))

	rate, err := k.ValidatorEmissionRate(ctx)
	if err != nil {
		return err
	}

	for _, vote := range sdkCtx.VoteInfos() {
		validator, err := k.GetValidatorByConsAddr(ctx, vote.Validator.Address)
		if err != nil {
			return err
		}

		validator.Points += pointsToAdd

		validator.Emission = validator.Emission.Add(rate.MulInt(validator.GetTokens().Mul(pointsToAddInt).Quo(monthlyPointsInt)))

		err = k.SetValidator(ctx, validator)

		if err != nil {
			return err
		}

		logger.Debug("Validators points", validator.GetOperator(), validator.GetPoints())
	}

	return nil
}

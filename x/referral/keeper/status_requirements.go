package keeper

import (
	"fmt"
	"github.com/pkg/errors"

	"github.com/axiome-pro/axm-node/x/referral/types"
)

func checkStatusRequirements(status types.Status, value types.Info) (types.StatusCheckResult, error) {
	if status == types.STATUS_UNSPECIFIED {
		return types.StatusCheckResult{Overall: true}, nil
	}
	return statusRequirements[status](value)
}

const (
	StatusGuruMinXCriteria = 3
	StatusGuruMinXParameter
	StatusLeaderMinXCriteria
)

var statusRequirements = map[types.Status]func(value types.Info) (types.StatusCheckResult, error){
	types.STATUS_NEW: func(_ types.Info) (types.StatusCheckResult, error) {
		return types.StatusCheckResult{Overall: true}, nil
	},
	types.STATUS_STARTER: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsSelfStake(value, 100)
	},
	types.STATUS_LEADER: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsXByX(value, 250, 20_000, StatusLeaderMinXCriteria, 0)
	},
	types.STATUS_GURU: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsXByX(value, 600, 50_000, StatusGuruMinXCriteria, StatusGuruMinXParameter)
	},
	types.STATUS_BOSS: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsCore(value, 1500, 150_000, 15)
	},
	types.STATUS_PRO: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsCore(value, 4_000, 300_000, 50)
	},
	types.STATUS_TOP: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsCore(value, 10_000, 800_000, 100)
	},
	types.STATUS_MEGA: func(value types.Info) (types.StatusCheckResult, error) {
		return statusRequirementsCore(value, 30_000, 2_000_000, 300)
	},
}

func statusRequirementsXByX(value types.Info, selfCoins, coins int64, count int, size int) (types.StatusCheckResult, error) {
	var (
		result    = types.StatusCheckResult{Overall: true}
		criterion types.StatusCheckResult_Criterion
	)

	if coins > 0 {
		criterion = types.StatusCheckResult_Criterion{
			Rule:        types.RULE_N_COINS_IN_STRUCTURE,
			TargetValue: uint64(coins),
			ActualValue: value.TeamDelegated.QuoRaw(1_000_000).Uint64(),
		}
		if criterion.ActualValue >= criterion.TargetValue {
			criterion.Met = true
			criterion.ActualValue = criterion.TargetValue
		}
		result.Criteria = append(result.Criteria, criterion)
		result.Overall = result.Overall && criterion.Met
	}

	if selfCoins > 0 {
		criterion = types.StatusCheckResult_Criterion{
			Rule:        types.RULE_SELF_STAKE,
			TargetValue: uint64(selfCoins),
			ActualValue: value.SelfDelegated.QuoRaw(1_000_000).Uint64(),
		}
		if criterion.ActualValue >= criterion.TargetValue {
			criterion.Met = true
			criterion.ActualValue = criterion.TargetValue
		}
		result.Criteria = append(result.Criteria, criterion)
		result.Overall = result.Overall && criterion.Met
	}

	criterion = types.StatusCheckResult_Criterion{
		Rule:        types.RULE_N_REFERRALS_WITH_X_REFERRALS_EACH,
		TargetValue: uint64(count),
		ParameterX:  uint64(size),
	}

	if size == 0 {
		criterion.Met = value.ActiveCount.FirstLine >= int32(count)
		criterion.ActualValue = uint64(value.ActiveCount.FirstLine)
	} else if size == 3 {
		criterion.Met = value.ActiveCount.FirstLineBy3 >= int32(count)
		criterion.ActualValue = uint64(value.ActiveCount.FirstLineBy3)
	} else {
		return result, errors.New(fmt.Sprintf("statusRequirementsXByX incorrect size %d", size))
	}

	if criterion.ActualValue > criterion.TargetValue {
		criterion.ActualValue = criterion.TargetValue
	}

	result.Criteria = append(result.Criteria, criterion)
	result.Overall = result.Overall && criterion.Met
	return result, nil
}

func statusRequirementsSelfStake(value types.Info, selfCoins int64) (types.StatusCheckResult, error) {
	var (
		result    = types.StatusCheckResult{Overall: true}
		criterion types.StatusCheckResult_Criterion
	)

	if selfCoins > 0 {
		criterion = types.StatusCheckResult_Criterion{
			Rule:        types.RULE_SELF_STAKE,
			TargetValue: uint64(selfCoins),
			ActualValue: value.SelfDelegated.QuoRaw(1_000_000).Uint64(),
		}
		if criterion.ActualValue >= criterion.TargetValue {
			criterion.Met = true
			criterion.ActualValue = criterion.TargetValue
		}
		result.Criteria = append(result.Criteria, criterion)
		result.Overall = result.Overall && criterion.Met
	}

	return result, nil
}

func statusRequirementsCore(value types.Info, selfCoins, coins int64, leg uint64) (types.StatusCheckResult, error) {
	var (
		result    = types.StatusCheckResult{Overall: true}
		criterion types.StatusCheckResult_Criterion
	)

	if coins > 0 {
		criterion = types.StatusCheckResult_Criterion{
			Rule:        types.RULE_N_COINS_IN_STRUCTURE,
			TargetValue: uint64(coins),
			ActualValue: value.TeamDelegated.QuoRaw(1_000_000).Uint64(),
		}
		if criterion.ActualValue >= criterion.TargetValue {
			criterion.Met = true
			criterion.ActualValue = criterion.TargetValue
		}
		result.Criteria = append(result.Criteria, criterion)
		result.Overall = result.Overall && criterion.Met
	}

	if selfCoins > 0 {
		criterion = types.StatusCheckResult_Criterion{
			Rule:        types.RULE_SELF_STAKE,
			TargetValue: uint64(selfCoins),
			ActualValue: value.SelfDelegated.QuoRaw(1_000_000).Uint64(),
		}
		if criterion.ActualValue >= criterion.TargetValue {
			criterion.Met = true
			criterion.ActualValue = criterion.TargetValue
		}
		result.Criteria = append(result.Criteria, criterion)
		result.Overall = result.Overall && criterion.Met
	}

	criterion = types.StatusCheckResult_Criterion{
		Rule:        types.RULE_N_TEAMS_OF_X_PEOPLE_EACH,
		TargetValue: 3,
		ParameterX:  leg,
	}
	xByX := types.StatusCheckResult_Criterion{
		Rule:        types.RULE_N_REFERRALS_WITH_X_REFERRALS_EACH,
		TargetValue: 3,
		ParameterX:  3,
	}

	xByX.Met = value.ActiveCount.FirstLineBy3 >= 3
	xByX.ActualValue = uint64(value.ActiveCount.FirstLineBy3)

	if xByX.ActualValue > xByX.TargetValue {
		xByX.ActualValue = xByX.TargetValue
	}

	var actualValue uint64 = 0

	switch leg {
	case 15:
		actualValue = uint64(value.ActiveCount.Team15 + value.ActiveCount.Team50 + value.ActiveCount.Team100 + value.ActiveCount.Team300)
	case 50:
		actualValue = uint64(value.ActiveCount.Team50 + value.ActiveCount.Team100 + value.ActiveCount.Team300)
	case 100:
		actualValue = uint64(value.ActiveCount.Team100 + value.ActiveCount.Team300)
	case 300:
		actualValue = uint64(value.ActiveCount.Team300)
	default:
		return result, errors.New(fmt.Sprintf("statusRequirementsCore incorrect leg %d", leg))
	}

	criterion.Met = actualValue >= 3
	if actualValue > 3 {
		criterion.ActualValue = 3
	} else {
		criterion.ActualValue = actualValue
	}

	result.Criteria = append(result.Criteria, criterion, xByX)
	result.Overall = result.Overall && criterion.Met && xByX.Met

	return result, nil
}

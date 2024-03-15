package keeper

import (
	"bytes"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/axiome-pro/axm-node/util"
	"github.com/axiome-pro/axm-node/x/referral/types"
	storetypes "cosmossdk.io/store/types"
)

type kvRecord struct {
	key   []byte
	value []byte
}

type callback struct {
	event string
	acc   string
}

func (x callback) Eq(y callback) bool {
	return x.event == y.event && x.acc == y.acc
}

type callbacks []callback

func (cbz *callbacks) Len() int {
	return len(*cbz)
}
func (cbz *callbacks) Less(i, j int) bool {
	x := (*cbz)[i]
	y := (*cbz)[j]

	res := strings.Compare(x.acc, y.acc)
	if res < 0 {
		return true
	} else if res > 0 {
		return false
	} else {
		return x.event < y.event
	}
}
func (cbz *callbacks) Swap(i, j int) {
	tmp := (*cbz)[i]
	(*cbz)[i] = (*cbz)[j]
	(*cbz)[j] = tmp
}

type bunchUpdater struct {
	k         Keeper
	ctx       sdk.Context
	data      []kvRecord
	callbacks callbacks
}

func newBunchUpdater(k Keeper, ctx sdk.Context) *bunchUpdater {
	return &bunchUpdater{
		k:         k,
		ctx:       ctx,
		data:      nil,
		callbacks: nil,
	}
}

func (bu *bunchUpdater) set(acc string, value types.Info) error {
	keyBytes := types.GetInfoAddrKey(acc)
	valueBytes, err := bu.k.cdc.Marshal(&value)
	if err != nil {
		return err
	}
	for i, record := range bu.data {
		if bytes.Equal(record.key, keyBytes) {
			bu.data[i].value = valueBytes
			return nil
		}
	}
	bu.data = append(bu.data, kvRecord{
		key:   keyBytes,
		value: valueBytes,
	})
	return nil
}

// TODO: Refactor to mitigate string <-> AccAddress casting
func (bu *bunchUpdater) get(acc string) (value types.Info, err error) {
	var (
		keyBytes   = types.GetInfoAddrKey(acc)
		valueBytes = []byte(nil)
	)
	for _, record := range bu.data {
		if bytes.Equal(record.key, keyBytes) {
			valueBytes = record.value
			break
		}
	}
	if valueBytes == nil {
		store := bu.k.storeService.OpenKVStore(bu.ctx)
		valueBytes, err = store.Get(keyBytes)
		if err != nil {
			return value, err
		}
	}
	err = bu.k.cdc.Unmarshal(valueBytes, &value)
	return value, err
}

func (bu bunchUpdater) StatusDowngradeAfter() time.Duration {
	return time.Duration(bu.k.GetParams(bu.ctx).StatusDowngradePeriod) * time.Second
}

func (bu *bunchUpdater) update(acc string, checkForStatusUpdate bool, callback func(value *types.Info) error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			if _, ok := e.(storetypes.ErrorOutOfGas); ok {
				panic(e)
			} else if er, ok := e.(error); ok {
				err = errors.Wrap(er, "update paniced")
			} else {
				err = errors.Errorf("update paniced: %s", e)
			}
		}
	}()
	value, err := bu.get(acc)
	if err != nil {
		bu.k.Logger(bu.ctx).Info("Cannot update, no such account", "addr", acc)
		return nil
	}
	value.Normalize()
	err = callback(&value)
	if err != nil {
		return errors.Wrap(err, "callback failed")
	}
	if checkForStatusUpdate {
		checkResult, err := checkStatusRequirements(value.Status, value)
		if err != nil {
			return err
		}
		if !checkResult.Overall {
			if value.StatusDowngradeAt == nil {
				downgradeAt := bu.ctx.BlockTime().Add(bu.StatusDowngradeAfter())
				value.StatusDowngradeAt = &downgradeAt
				err = bu.k.ScheduleStatusDowngrade(bu.ctx, acc, downgradeAt)
				if err != nil {
					return err
				}
				util.EmitEvent(bu.ctx,
					&types.EventStatusWillBeDowngraded{
						Address: acc,
						Time:    downgradeAt,
					},
				)
			}
		} else {
			if value.StatusDowngradeAt != nil {
				err = bu.k.RemoveStatusDowngradeSchedule(bu.ctx, acc, *value.StatusDowngradeAt)
				if err != nil {
					return err
				}
				value.StatusDowngradeAt = nil
				util.EmitEvent(bu.ctx,
					&types.EventStatusDowngradeCanceled{
						Address: acc,
					},
				)
			}
			var nextStatus = value.Status
			for {
				if nextStatus == types.MaximumStatus {
					break
				}
				nextStatus++

				checkResult, err = checkStatusRequirements(nextStatus, value)
				if err != nil {
					return err
				}
				if !checkResult.Overall {
					nextStatus--
					break
				}
			}
			if nextStatus > value.Status {
				util.EmitEvent(bu.ctx,
					&types.EventStatusUpdated{
						Address: acc,
						Before:  value.Status,
						After:   nextStatus,
					},
				)
				bu.k.setStatus(bu.ctx, &value, nextStatus, acc)
				bu.addCallback(StatusUpdatedCallback, acc)
			}
		}
	}
	if err := bu.set(acc, value); err != nil {
		return err
	}
	return nil
}

func (bu *bunchUpdater) addCallback(eventName string, acc string) {
	bu.callbacks = append(bu.callbacks, callback{event: eventName, acc: acc})
}

func (bu *bunchUpdater) commit() error {
	store := bu.k.storeService.OpenKVStore(bu.ctx)
	for _, pair := range bu.data {
		err := store.Set(pair.key, pair.value)
		if err != nil {
			return err
		}
	}
	sort.Sort(&bu.callbacks)
	for i, cb := range bu.callbacks {
		if i > 0 && bu.callbacks[i-1].Eq(cb) {
			continue
		}
		if err := bu.k.callback(cb.event, bu.ctx, cb.acc); err != nil {
			return errors.Wrap(err, cb.event+" callback failed for "+cb.acc)
		}
	}
	return nil
}

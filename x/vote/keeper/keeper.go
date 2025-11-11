package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	corestore "cosmossdk.io/core/store"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/cachekv"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdktx "github.com/cosmos/cosmos-sdk/types/tx"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"cosmossdk.io/log"

	"github.com/axiome-pro/axm-node/util"
	"github.com/axiome-pro/axm-node/x/vote/types"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the vote store
type Keeper struct {
	cdc            codec.Codec
	storeService   corestore.KVStoreService
	Params         collections.Item[types.Params]
	referralKeeper types.ReferralKeeper
	authority      sdk.AccAddress
	accountKeeper  types.AccountKeeper

	router baseapp.MessageRouter
}

// NewKeeper creates a vote keeper
func NewKeeper(
	cdc codec.Codec, storeService corestore.KVStoreService,
	referralKeeper types.ReferralKeeper, authority sdk.AccAddress,
	router baseapp.MessageRouter, accountKeeper types.AccountKeeper,
) Keeper {
	sb := collections.NewSchemaBuilder(storeService)

	keeper := Keeper{
		cdc:            cdc,
		storeService:   storeService,
		referralKeeper: referralKeeper,
		accountKeeper:  accountKeeper,
		Params:         collections.NewItem(sb, types.KeyParams, "params", codec.CollValue[types.Params](cdc)),
		authority:      authority,
		router:         router,
	}
	return keeper
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (k Keeper) Router() baseapp.MessageRouter {
	return k.router
}

func (k Keeper) GetCurrentProposal(ctx sdk.Context) *types.Proposal {
	store := k.storeService.OpenKVStore(ctx)

	bz, err := store.Get(types.KeyCurrentVote)

	if err != nil {
		panic(err)
	}

	if bz == nil {
		return nil
	}

	var proposal types.Proposal
	err = proto.Unmarshal(bz, &proposal)
	if err != nil {
		panic(err)
	}

	err = proposal.UnpackInterfaces(k.cdc)
	if err != nil {
		panic(err)
	}

	return &proposal
}

func (k Keeper) SetCurrentProposal(ctx sdk.Context, proposal types.Proposal) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := proto.Marshal(&proposal)
	if err != nil {
		panic(err)
	}
	err = store.Set(types.KeyCurrentVote, bz)

	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetAgreed(ctx sdk.Context) (gov types.Government) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyAgreedMembers)
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.Government{}
	}

	if err := proto.Unmarshal(bz, &gov); err != nil {
		panic(err)
	}
	return gov
}

func (k Keeper) SetAgreed(ctx sdk.Context, agreed types.Government) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := proto.Marshal(&agreed)
	if err != nil {
		panic(err)
	}
	err = store.Set(types.KeyAgreedMembers, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetDisagreed(ctx sdk.Context) (gov types.Government) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyDisagreedMembers)

	if err != nil {
		panic(err)
	}

	if bz == nil {
		return types.Government{}
	}

	if err = proto.Unmarshal(bz, &gov); err != nil {
		panic(err)
	}
	return gov
}

func (k Keeper) SetDisagreed(ctx sdk.Context, disagreed types.Government) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := proto.Marshal(&disagreed)
	if err != nil {
		panic(err)
	}
	err = store.Set(types.KeyDisagreedMembers, bz)

	if err != nil {
		panic(err)
	}
}

func (k Keeper) Validate(gov types.Government,
	aGov types.Government,
	dGov types.Government,
) (complete bool, agreed bool) {

	if len(gov.Members) == (len(aGov.Members) + len(dGov.Members)) {
		complete = true
		agreed = len(aGov.Members)*3 >= len(gov.Members)*2
	}

	return complete, agreed
}

func (k Keeper) SaveProposalToHistory(ctx sdk.Context, store corestore.KVStore) {
	history := types.ProposalHistoryRecord{
		Proposal:   *k.GetCurrentProposal(ctx),
		Government: k.GetGovernment(ctx).Members,
		Agreed:     k.GetAgreed(ctx).Members,
		Disagreed:  k.GetDisagreed(ctx).Members,
		Started:    k.GetStartBlock(ctx),
		Finished:   ctx.BlockHeight(),
	}

	historyBz, err := proto.Marshal(&history)
	if err != nil {
		panic(err)
	}
	height := make([]byte, 8)
	binary.BigEndian.PutUint64(height, uint64(ctx.BlockHeight()))
	key := append(types.KeyHistoryPrefix, height...)
	store.Set(key, historyBz)
}

func (k Keeper) AddProposalHistoryRecord(ctx sdk.Context, record types.ProposalHistoryRecord) {
	store := k.storeService.OpenKVStore(ctx)
	historyBz, err := proto.Marshal(&record)
	if err != nil {
		panic(err)
	}
	height := make([]byte, 8)
	binary.BigEndian.PutUint64(height, uint64(record.Finished))
	key := append(types.KeyHistoryPrefix, height...)
	err = store.Set(key, historyBz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) SetStartBlock(ctx sdk.Context) {
	store := k.storeService.OpenKVStore(ctx)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, (uint64(ctx.BlockHeight())))
	err := store.Set(types.KeyStartBlock, bz)
	if err != nil {
		panic(err)
	}
}

func (k Keeper) GetStartBlock(ctx sdk.Context) int64 {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.KeyStartBlock)
	if err != nil {
		panic(err)
	}
	return int64(binary.BigEndian.Uint64(bz))
}

func safeExecuteHandler(ctx sdk.Context, msg sdk.Msg, handler baseapp.MsgServiceHandler,
) (res *sdk.Result, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("handling x/vote proposal msg [%s] PANICKED: %v", msg, r)
		}
	}()
	res, err = handler(ctx, msg)
	return
}

func (k Keeper) EndProposal(ctx sdk.Context, proposal types.Proposal, agreed bool) {
	store := k.storeService.OpenKVStore(ctx)

	// Save proposal data to history
	k.SaveProposalToHistory(ctx, store)

	// Delete all proposal info
	err := store.Delete(types.KeyCurrentVote)
	if err != nil {
		panic(err)
	}
	err = store.Delete(types.KeyAgreedMembers)
	if err != nil {
		panic(err)
	}
	err = store.Delete(types.KeyDisagreedMembers)
	if err != nil {
		panic(err)
	}
	err = store.Delete(types.KeyStartBlock)
	if err != nil {
		panic(err)
	}

	util.EmitEvent(ctx,
		&types.EventvoteFinished{
			Name:   proposal.Name,
			Agreed: agreed,
		},
	)

	if agreed {
		var (
			idx    int
			events sdk.Events
			msg    sdk.Msg
		)

		cacheCtx, writeCache := ctx.CacheContext()
		messages, err := sdktx.GetMsgs(proposal.Messages, "EndProposal")
		if err != nil {
			panic(err)
		}

		for idx, msg = range messages {
			handler := k.Router().Handler(msg)

			var res *sdk.Result
			k.Logger(cacheCtx).Error("handling msg", "msg", msg)
			res, err = safeExecuteHandler(cacheCtx, msg, handler)
			if err != nil {
				break
			}

			events = append(events, res.GetEvents()...)
		}

		if err == nil {
			// write state to the underlying multi-store
			writeCache()

			// propagate the msg events to the current context
			ctx.EventManager().EmitEvents(events)
		} else {
			k.Logger(ctx).Error("could not apply vote result due to error",
				"name", proposal.Name,
				"error", err,
				"result",
				fmt.Sprintf("passed, but msg %d (%s) failed on execution: %s", idx, sdk.MsgTypeURL(msg), err.Error()),
			)
		}
	}
}

func (k Keeper) GetHistory(ctx sdk.Context, limit int32, page int32) []types.ProposalHistoryRecord {
	store := k.storeService.OpenKVStore(ctx)

	iterator, err := store.Iterator(types.KeyHistoryPrefix, storetypes.PrefixEndBytes(types.KeyHistoryPrefix))
	if err != nil {
		panic(err)
	}
	defer iterator.Close()

	records := make([]types.ProposalHistoryRecord, 0)
	start := limit * (page - 1)
	end := limit * page

	for current := int32(0); iterator.Valid() && (current < end); iterator.Next() {
		if current < start {
			current++
			continue
		}
		current++
		var record types.ProposalHistoryRecord
		if err := proto.Unmarshal(iterator.Value(), &record); err != nil {
			panic(err)
		}
		records = append(records, record)
	}

	return records
}

func (k Keeper) Propose(ctx sdk.Context, msg types.MsgPropose) error {
	if k.GetCurrentProposal(ctx) != nil {
		return types.ErrOtherActive
	}

	var (
		proposal = types.Proposal{}
		gov      = k.GetGovernment(ctx)
	)
	if !gov.Contains(msg.GetAuthor()) {
		return errors.Wrap(types.ErrSignerNotAllowed, msg.Author)
	}

	proposal.Name = msg.Name
	proposal.Author = msg.Author

	params := k.GetParams(ctx)
	endTime := ctx.BlockTime().Add(time.Duration(params.VotePeriod) * time.Minute)
	proposal.EndTime = &endTime

	messages, err := sdktx.GetMsgs(msg.Messages, "sdk.Msg")
	if err != nil {
		return errors.Wrap(err, "Unable to create proposal")
	}

	// validate messages
	for _, msg := range messages {
		signers, _, err := k.cdc.GetMsgV1Signers(msg)
		if err != nil {
			return err
		}
		if len(signers) != 1 {
			return types.ErrInvalidSigner
		}

		// assert that the governance module account is the only signer of the messages
		if !bytes.Equal(signers[0], sdk.AccAddress(k.authority)) {
			return errorsmod.Wrapf(types.ErrInvalidSigner, sdk.AccAddress(signers[0]).String())
		}

		// use the msg service router to see that there is a valid route for that message.
		handler := k.router.Handler(msg)
		if handler == nil {
			return errorsmod.Wrap(types.ErrUnroutableProposalMsg, sdk.MsgTypeURL(msg))
		}

	}

	proposal.Messages = msg.Messages

	// Set proposal
	k.SetCurrentProposal(ctx, proposal)

	// Set empty lists of voters
	agreed, disagreed := types.Government{Members: []string{proposal.Author}}, types.Government{}
	k.SetAgreed(ctx, agreed)
	k.SetDisagreed(ctx, disagreed)
	k.SetStartBlock(ctx)

	util.EmitEvent(ctx,
		&types.EventProposalCreated{
			Name:   proposal.Name,
			Author: proposal.Author,
		},
	)

	if complete, agree := k.Validate(gov, agreed, disagreed); complete {
		k.EndProposal(ctx, proposal, agree)
	}
	return nil
}

func (k Keeper) Vote(ctx sdk.Context, voter sdk.AccAddress, agree bool) error {
	proposal := k.GetCurrentProposal(ctx)
	if proposal == nil {
		return types.ErrNoActiveProposals
	}

	gov := k.GetGovernment(ctx)
	if !gov.Contains(voter) {
		return errors.Wrap(types.ErrSignerNotAllowed, voter.String())
	}

	agreed := k.GetAgreed(ctx)
	if agreed.Contains(voter) {
		return errors.Wrap(types.ErrAlreadyVoted, voter.String())
	}

	disagreed := k.GetDisagreed(ctx)
	if disagreed.Contains(voter) {
		return errors.Wrap(types.ErrAlreadyVoted, voter.String())
	}

	if agree {
		agreed.Append(voter)
		k.SetAgreed(ctx, agreed)
	} else {
		disagreed.Append(voter)
		k.SetDisagreed(ctx, disagreed)
	}

	util.EmitEvent(ctx,
		&types.EventProposalVote{
			Voter:  voter.String(),
			Agreed: agree,
		},
	)

	if complete, agree := k.Validate(gov, agreed, disagreed); complete {
		k.EndProposal(ctx, *proposal, agree)
	}
	return nil
}

func (k Keeper) GetCurrentPoll(ctx sdk.Context) (poll types.Poll, ok bool) {
	store := k.storeService.OpenKVStore(ctx)
	bz, err := store.Get(types.GetPollPrefixedKey(types.KeyPollCurrent))
	if err != nil {
		panic(err)
	}
	if bz == nil {
		return types.Poll{}, false
	}

	k.cdc.MustUnmarshal(bz, &poll)
	return poll, true
}

func (k Keeper) GetPollStatus(ctx sdk.Context) (yes, no uint64) {
	store := k.storeService.OpenKVStore(ctx)
	if bz, _ := store.Get(types.GetPollPrefixedKey(types.KeyPollYesCount)); bz != nil {
		yes = binary.BigEndian.Uint64(bz)
	}
	if bz, _ := store.Get(types.GetPollPrefixedKey(types.KeyPollNoCount)); bz != nil {
		no = binary.BigEndian.Uint64(bz)
	}
	return
}

func (k Keeper) StartPoll(ctx sdk.Context, poll types.Poll) error {
	store := k.storeService.OpenKVStore(ctx)

	if _, err := store.Has(types.GetPollPrefixedKey(types.KeyPollCurrent)); err != nil {
		return types.ErrOtherActive
	}
	if !util.ContainsString(k.GetGovernment(ctx).Strings(), poll.Author) {
		return types.ErrSignerNotAllowed
	}

	start := ctx.BlockTime()
	end := start.Add(time.Duration(k.GetParams(ctx).PollPeriod) * time.Minute)
	poll.StartTime = &start
	poll.EndTime = &end

	return store.Set(types.KeyPollCurrent, k.cdc.MustMarshal(&poll))
}

func (k Keeper) Answer(ctx sdk.Context, acc string, yes bool) error {
	poll, ok := k.GetCurrentPoll(ctx)
	if !ok {
		return types.ErrNoActivePoll
	}

	_, err := sdk.AccAddressFromBech32(acc)
	if err != nil {
		panic(errors.Wrap(err, "cannot parse acc address"))
	}

	// TODO: add poll for validators
	switch r := poll.Requirements.(type) {
	case *types.Poll_MinStatus:
		info, err := k.referralKeeper.Get(ctx, acc)
		if err != nil {
			panic(errors.Wrap(err, "cannot obtain referral info"))
		}
		if info.Status < r.MinStatus {
			return types.ErrRespondentNotAllowed
		}
	}

	store := k.storeService.OpenKVStore(ctx)
	key := types.GetPollAnswersPrefixedKey([]byte(acc))
	if has, _ := store.Has(key); has {
		return types.ErrAlreadyVoted
	}

	var ans, countKey []byte
	if yes {
		ans = types.ValueYes
		countKey = types.KeyPollYesCount
	} else {
		ans = types.ValueNo
		countKey = types.KeyPollNoCount
	}
	err = store.Set(key, ans)
	if err != nil {
		panic(err)
	}

	var (
		bz    []byte
		value uint64
	)
	if bz, _ = store.Get(types.GetPollPrefixedKey(countKey)); bz != nil {
		value = binary.BigEndian.Uint64(bz)
	} else {
		bz = make([]byte, 8)
	}
	value += 1
	binary.BigEndian.PutUint64(bz, value)
	return store.Set(types.GetPollAnswersPrefixedKey(countKey), bz)
}

func (k Keeper) EndPollHandler(ctx sdk.Context, _ []byte, _ time.Time) { k.EndPoll(ctx) }
func (k Keeper) EndPoll(ctx sdk.Context) {
	store := cachekv.NewStore(
		prefix.NewStore(
			runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)),
			types.KeyPollPrefix,
		),
	)

	var (
		poll     types.Poll
		yes, no  uint64
		decision types.Decision
	)
	if bz := store.Get(types.KeyPollCurrent); bz != nil {
		k.cdc.MustUnmarshal(bz, &poll)
	} else {
		panic(types.ErrNoActivePoll)
	}
	if bz := store.Get(types.KeyPollYesCount); bz != nil {
		yes = binary.BigEndian.Uint64(bz)
	}
	if bz := store.Get(types.KeyPollNoCount); bz != nil {
		no = binary.BigEndian.Uint64(bz)
	}
	if poll.Quorum != nil {
		if yes != 0 && util.FractionInt(int64(yes)).GTE(poll.Quorum.Mul(util.FractionInt(int64(yes+no)))) {
			decision = types.DECISION_POSITIVE
		} else {
			decision = types.DECISION_NEGATIVE
		}
	}

	util.EmitEvent(ctx,
		&types.EventPollFinished{
			Name:     poll.Name,
			Yes:      yes,
			No:       no,
			Decision: decision,
		},
	)

	historyKey := make([]byte, len(types.KeyPollHistory)+8)
	copy(historyKey, types.KeyPollHistory)
	binary.BigEndian.PutUint64(historyKey[len(types.KeyPollHistory):], uint64(poll.EndTime.Unix()))
	store.Set(historyKey, k.cdc.MustMarshal(&types.PollHistoryItem{
		Poll:     poll,
		Yes:      yes,
		No:       no,
		Decision: decision,
	}))

	store.Delete(types.KeyPollCurrent)
	store.Delete(types.KeyPollYesCount)
	store.Delete(types.KeyPollNoCount)
	it := storetypes.KVStorePrefixIterator(store, types.KeyPollAnswers)
	for ; it.Valid(); it.Next() {
		store.Delete(it.Key())
	}
	it.Close()

	store.Write()
}

func (k Keeper) GetPollHistoryAll(ctx sdk.Context) []types.PollHistoryItem {
	return k.GetPollHistory(ctx, 0, 0)
}
func (k Keeper) GetPollHistory(ctx sdk.Context, limit int32, page int32) []types.PollHistoryItem {
	store := prefix.NewStore(
		runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPollPrefix,
	)
	var (
		it  storetypes.Iterator
		res []types.PollHistoryItem
	)
	if limit > 0 {
		it = storetypes.KVStorePrefixIteratorPaginated(store, types.KeyPollHistory, uint(page), uint(limit))
		res = make([]types.PollHistoryItem, 0, limit)
	} else {
		it = storetypes.KVStorePrefixIterator(store, types.KeyPollHistory)
	}
	for ; it.Valid(); it.Next() {
		var item types.PollHistoryItem
		k.cdc.MustUnmarshal(it.Value(), &item)
		res = append(res, item)
	}
	it.Close()
	return res
}

func (k Keeper) IterateThroughCurrentPollAnswers(ctx sdk.Context, callback func(acc string, ans bool) (stop bool)) (err error) {
	store := prefix.NewStore(
		runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPollPrefix,
	)
	if !store.Has(types.KeyPollCurrent) {
		return types.ErrNoActivePoll
	}

	it := storetypes.KVStorePrefixIterator(store, types.KeyPollAnswers)
	defer func() {
		it.Close()
		if e := recover(); e != nil {
			if er, ok := e.(error); ok {
				err = errors.Wrap(er, "callback paniced")
			} else {
				err = errors.Errorf("callback paniced: %s", er)
			}
		}
	}()

	for ; it.Valid(); it.Next() {
		acc := string(it.Key()[len(types.KeyPollAnswers):])
		ans := bytes.Equal(it.Value(), types.ValueYes)
		if stop := callback(acc, ans); stop {
			return nil
		}
	}
	return nil
}

func (k Keeper) LoadPolls(ctx sdk.Context, state types.GenesisState) {
	store := prefix.NewStore(
		runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPollPrefix,
	)

	if state.CurrentPoll != nil {
		store.Set(types.KeyPollCurrent, k.cdc.MustMarshal(state.CurrentPoll))
		for _, ans := range state.PollAnswers {
			if err := k.Answer(ctx, ans.Acc, ans.Ans); err != nil {
				panic(err)
			}
		}
	}

	key := make([]byte, len(types.KeyPollHistory)+8)
	copy(key, types.KeyPollHistory)
	for _, item := range state.PollHistory {
		binary.BigEndian.PutUint64(key[len(types.KeyPollHistory):], uint64(item.Poll.EndTime.Unix()))
		store.Set(key, k.cdc.MustMarshal(&item))
	}
}

package types

import (
	"bytes"
	"fmt"
	"time"

	"cosmossdk.io/collections"
	"encoding/binary"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/kv"
	"github.com/pkg/errors"
)

const (
	// ModuleName is the name of the module
	ModuleName          = "referral"
	ReferralAccountName = ModuleName

	EventTypeRefFee  = "ref_fee"
	AttributeKeyFrom = "from"
	AttributeKeyTo   = "to"
)

// Keys for referral store
// Items are stored with the following key: values
//
// - 0x00<accAddrLen (1 Byte)><accAddr_Bytes>: Info
//
// - 0x01<accAddrLen (1 Byte)><accAddr_Bytes><refAddrLen (1 Byte)><refAddr_Bytes>: Referral connection
//
// - 0x02<accAddrLen (1 Byte)><accAddr_Bytes><refAddrLen (1 Byte)><refAddr_Bytes>: Referral connection
//
// - 0x03: Params
var (
	InfoPrefix           = []byte{0x00}
	ReferralsPrefix      = []byte{0x01}
	ParamsKey            = collections.NewPrefix(2)
	DowngradeQueuePrefix = []byte{0x05}
)

// GetInfoAddrKey creates the key for a referral info record.
func GetInfoAddrKey(acc string) []byte {
	return append(InfoPrefix, []byte(acc)...)
}

// ParseInfoAddrKey creates the address from InfoAddrKey
func ParseInfoAddrKey(key []byte) string {
	kv.AssertKeyAtLeastLength(key, 2)
	return string(key[1:]) // remove prefix bytes and address length
}

// GetReferralsRelationKey create the key for a referral <=> referrer relation
func GetReferralsRelationKey(acc string, ref string) []byte {
	return append(append(ReferralsPrefix, address.MustLengthPrefix([]byte(acc))...), address.MustLengthPrefix([]byte(ref))...)
}

func GetReferralsChildIteratorKey(acc string) []byte {
	return append(ReferralsPrefix, address.MustLengthPrefix([]byte(acc))...)
}

// ParseInfoAddrKey creates the address from InfoAddrKey
func ParseReferralFromReleationKey(key []byte) (string, string, error) {
	prefixLength := len(ReferralsPrefix)
	if prefix := key[:prefixLength]; !bytes.Equal(prefix, ReferralsPrefix) {
		return "", "", fmt.Errorf("invalid prefix; expected: %X, got: %x", ReferralsPrefix, prefix)
	}

	key = key[prefixLength:] // remove the prefix byte
	if len(key) == 0 {
		return "", "", fmt.Errorf("no bytes left to parse: %X", key)
	}

	referrerAddLen := key[0]
	key = key[1:] // remove the length byte of referrer address.
	if len(key) == 0 {
		return "", "", fmt.Errorf("no bytes left to parse validator address: %X", key)
	}

	referrer := key[0:int(referrerAddLen)]

	key = key[int(referrerAddLen):] // remove the referrer address bytes
	if len(key) <= 1 {
		return "", "", fmt.Errorf("no bytes left to parse delegator address: %X", key)
	}

	referral := key[1:]

	return string(referrer), string(referral), nil
}

func GetDowngradeQueueKey(acc string, timestamp time.Time) []byte {
	bz := FormatTimeBytes(timestamp)
	return append(append(DowngradeQueuePrefix, bz...), []byte(acc)...)
}

func GetDowngradeQueueIteratorStartKey() []byte {
	return DowngradeQueuePrefix
}

func GetDowngradeQueueIteratorEndKey(timestamp time.Time) []byte {
	return GetDowngradeQueueKey("", timestamp.Add(time.Nanosecond))
}

func ExtractAccFromDowngradeQueueKey(key []byte) string {
	return string(key[len(DowngradeQueuePrefix)+8:])
}

func FormatTimeBytes(t time.Time) []byte {
	if t.Year() < 1970 || t.Year() > 2262 {
		panic(errors.Errorf("time is out of range: %s", t.String()))
	}
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(t.UnixNano()))
	return bz
}

// Package steamid implements helpers for handling SteamIDs.
//
// The implementation is based on and compliant with the details from
// https://developer.valvesoftware.com/wiki/SteamID.
package steamid

import (
	"fmt"
	"strconv"
)

// FromValues creates a SteamID from the given values.
func FromValues(universe uint8, accountInstance uint32, accountType AccountType, accountID uint32) SteamID {
	var id SteamID
	id.setBits(0, 1<<32-1, uint64(accountID))
	id.setBits(32, 1<<20-1, uint64(accountInstance))
	id.setBits(32+20, 1<<4-1, uint64(accountType))
	id.setBits(32+20+4, 1<<8-1, uint64(universe))
	return id
}

// SteamID represents a single SteamID, including its universe, instance, type,
// and ID.
//
// See
// https://developer.valvesoftware.com/wiki/SteamID#As_Represented_in_Computer_Programs
// for implementation details.
type SteamID uint64

func (id SteamID) getBits(offset uint16, mask uint64) uint64 {
	return (uint64(id) >> offset) & mask
}

func (id *SteamID) setBits(offset uint16, mask uint64, value uint64) {
	*id = SteamID((uint64(*id) & ^(mask << offset)) | ((value & mask) << offset))
}

// Universe returns the universe of this SteamID.
//
// See
// https://developer.valvesoftware.com/wiki/SteamID#Universes_Available_for_Steam_Accounts
// for available universes.
func (id SteamID) Universe() uint8 {
	return uint8(id.getBits(32+20+4, 1<<8-1))
}

// AccountID returns the account ID of this SteamID. This is the unique
// identifier for the account.
func (id SteamID) AccountID() uint32 {
	return uint32(id.getBits(0, 1<<32-1))
}

// AccountInstance returns the instance of the account. It is usually 1 for user
// accounts.
func (id SteamID) AccountInstance() uint32 {
	return uint32(id.getBits(32, 1<<20-1))
}

// AccountType returns the account type of this SteamID. For more information
// about account types, see the AccountType type.
func (id SteamID) AccountType() AccountType {
	return AccountType(id.getBits(32+20, 1<<4-1))
}

// SteamID2 returns the version 2 textual representation of this SteamID. For
// example, STEAM_0:0:1.
//
// When a SteamID type can not be turned into a version 2 ID, it returns an
// empty string.
func (id SteamID) SteamID2() string {
	if accType := id.AccountType(); accType == AccountTypeInvalid {
		return "UNKNOWN"
	} else if accType == AccountTypeIndividual {
		accountID := id.AccountID()
		if universe := id.Universe(); universe > UniversePublic {
			return "STEAM_" + strconv.FormatInt(int64(id.Universe()), 10) + ":" + strconv.FormatInt(int64(accountID&1), 10) + ":" + strconv.FormatInt(int64(accountID>>1), 10)
		}
		return "STEAM_0:" + strconv.FormatInt(int64(accountID&1), 10) + ":" + strconv.FormatInt(int64(accountID>>1), 10)
	} else if accType == AccountTypePending {
		return "STEAM_ID_PENDING"
	}

	return ""
}

// SteamID3 returns the version 3 textual representation of this SteamID. For
// example, [U:1:2].
//
// When a SteamID type can not be turned into a version 3 ID, it returns an
// empty string.
func (id SteamID) SteamID3() string {
	switch id.AccountType() {
	case AccountTypeInvalid:
		return "[I:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case AccountTypeIndividual:
		if accountInstance := id.AccountInstance(); accountInstance != 1 {
			// Not a desktop instance
			return "[U:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + ":" + strconv.FormatUint(uint64(accountInstance), 10) + "]"
		}
		return "[U:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case AccountTypeMultiseat:
		return "[M:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + ":" + strconv.FormatUint(uint64(id.AccountInstance()), 10) + "]"
	case AccountTypeGameServer:
		return "[G:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case AccountTypeAnonGameServer:
		return "[A:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + ":" + strconv.FormatUint(uint64(id.AccountInstance()), 10) + "]"
	case AccountTypePending:
		return "[P:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case AccountTypeContentServer:
		return "[C:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case AccountTypeClan:
		return "[g:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case AccountTypeChat:
		const instanceMask = 0x000FFFFF

		const (
			clanFlag             = (instanceMask + 1) >> 1
			lobbyFlag            = (instanceMask + 1) >> 2
			matchmakingLobbyFlag = (instanceMask + 1) >> 3
		)

		if accountInstance := id.AccountInstance(); accountInstance&clanFlag == clanFlag {
			return "[c:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		} else if accountInstance&lobbyFlag == lobbyFlag {
			return "[L:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		} else if accountInstance&matchmakingLobbyFlag == matchmakingLobbyFlag {
			return "[T:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		}
	}

	return ""
}

// String returns a string representation of this SteamID. It attempts to use a
// version 2 SteamID first, then a version 3 SteamID, and finally the raw
// representation of this SteamID. The last of those can not be parsed by
// Parse.
func (id SteamID) String() string {
	if version2 := id.SteamID2(); len(version2) > 0 {
		return version2
	}

	if version3 := id.SteamID3(); len(version3) > 0 {
		return version3
	}

	return strconv.FormatInt(int64(id), 10)
}

// Implementation of encoding.TextMarshaler & encoding.TextUnmarshaler

// MarshalText implements encoding.TextMarshaler
func (id SteamID) MarshalText() (text []byte, err error) {
	if version2 := id.SteamID2(); len(version2) > 0 {
		return []byte(version2), nil
	}

	if version3 := id.SteamID3(); len(version3) > 0 {
		return []byte(version3), nil
	}

	return nil, fmt.Errorf("Cannot marshal account of type %d", id.AccountType())
}

// UnmarshalText implements encoding.TextUnmarshaler
func (id *SteamID) UnmarshalText(text []byte) (err error) {
	*id, err = Parse(string(text))
	return
}

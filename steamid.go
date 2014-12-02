// Package steamid implements helpers for handling SteamIDs.
//
// The implementation is based on and compliant with the details from
// https://developer.valvesoftware.com/wiki/SteamID.
package steamid

import (
	"strconv"
)

// FromCommunityID creates a SteamID from a community ID.
//
// Note that this only works with AccountTypes that have modifiers. By default,
// these are only 1 and 7.
//
// For more information about account types, see
// https://developer.valvesoftware.com/wiki/SteamID#Types_of_Steam_Accounts
func FromCommunityID(id uint64, accountType uint8) SteamID {
	accType := accountTypes[accountType]
	authServer := uint8(id % 2)
	accountID := uint32((id - accType.Modifier - uint64(authServer)) / 2)

	return FromValues(0, 1, accType.Number, (accountID<<1)|uint32(authServer))
}

// FromValues creates a SteamID from the given values.
func FromValues(universe uint8, accountInstance uint32, accountType uint8, accountID uint32) SteamID {
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

func (id SteamID) accountTypeNumber() uint8 {
	return uint8(id.getBits(32+20, 1<<4-1))
}

// AccountType returns the account type number and name (if known). If the name
// is not known, an empty string is returned instead.
func (id SteamID) AccountType() (int, string) {
	n := int(id.accountTypeNumber())
	if n >= len(accountTypes) {
		return n, ""
	}
	return n, accountTypes[n].Name
}

// SteamID2 returns the version 2 textual representation of this SteamID. For
// example, STEAM_0:0:1.
func (id SteamID) SteamID2() string {
	if accTypeNum := id.accountTypeNumber(); accTypeNum == 1 {
		accountID := id.AccountID()
		if universe := id.Universe(); universe > UniversePublic {
			return "STEAM_" + strconv.FormatInt(int64(id.Universe()), 10) + ":" + strconv.FormatInt(int64(accountID&1), 10) + ":" + strconv.FormatInt(int64(accountID>>1), 10)
		}
		return "STEAM_0:" + strconv.FormatInt(int64(accountID&1), 10) + ":" + strconv.FormatInt(int64(accountID>>1), 10)
	} else if accTypeNum == 5 {
		return "STEAM_ID_PENDING"
	} else {
		return "INVALID"
	}
}

// SteamID3 returns the version 3 textual representation of this SteamID. For
// example, [U:1:2].
func (id SteamID) SteamID3() string {
	switch id.accountTypeNumber() {
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

	return "Unknown type"
}

// An alias of SteamID3
func (id SteamID) String() string {
	return id.SteamID3()
}

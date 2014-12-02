package steamid

import (
	"strconv"
)

// Create a SteamID from a community ID. Note that this only works with AccountTypes that has a Modifier (only Individual (1) and Clan (7) by default)
func FromCommunityID(id uint64, accountType *AccountType) SteamID {
	var authServer uint8 = uint8(id % 2)
	var accountID uint32 = uint32((id - accountType.Modifier - uint64(authServer)) / 2)

	return FromValues(0, 1, AccountTypes[1], (accountID<<1)|uint32(authServer))
}

// Create a SteamID from a set of values
func FromValues(universe uint8, accountInstance uint32, accountType *AccountType, accountID uint32) SteamID {
	var id SteamID
	id.setBits(0, 1<<32-1, uint64(accountID))
	id.setBits(32, 1<<20-1, uint64(accountInstance))
	id.setBits(32+20, 1<<4-1, uint64(accountType.Number))
	id.setBits(32+20+4, 1<<8-1, uint64(universe))
	return id
}

// Holds the SteamID information. See https://developer.valvesoftware.com/wiki/SteamID#As_Represented_in_Computer_Programs for implementation details.
type SteamID uint64

func (id SteamID) getBits(offset uint16, mask uint64) uint64 {
	return (uint64(id) >> offset) & mask
}

func (id *SteamID) setBits(offset uint16, mask uint64, value uint64) {
	*id = SteamID((uint64(*id) & ^(mask << offset)) | ((value & mask) << offset))
}

// Get the universe of the SteamID
func (id SteamID) Universe() uint8 {
	return uint8(id.getBits(32+20+4, 1<<8-1))
}

// Get the account's ID
func (id SteamID) AccountID() uint32 {
	return uint32(id.getBits(0, 1<<32-1))
}

// Get the instance of the account
func (id SteamID) AccountInstace() uint32 {
	return uint32(id.getBits(32, 1<<20-1))
}

func (id SteamID) accountTypeNumber() uint8 {
	return uint8(id.getBits(32+20, 1<<4-1))
}

// Get the account's type
func (id SteamID) AccountType() *AccountType {
	return AccountTypes[id.accountTypeNumber()]
}

// Get the version 2 textual representation
func (id SteamID) SteamID2() string {
	if accountType := id.AccountType(); accountType.Number == 1 {
		accountID := id.AccountID()
		if universe := id.Universe(); universe <= UniversePublic {
			return "STEAM_0:" + strconv.FormatInt(int64(accountID&1), 10) + ":" + strconv.FormatInt(int64(accountID>>1), 10)
		} else {
			return "STEAM_" + strconv.FormatInt(int64(id.Universe()), 10) + ":" + strconv.FormatInt(int64(accountID&1), 10) + ":" + strconv.FormatInt(int64(accountID>>1), 10)
		}
	} else if accountType.Number == 5 {
		return "STEAM_ID_PENDING"
	} else {
		return "INVALID"
	}
}

// Get the version 3 textual representation
func (id SteamID) SteamID3() string {
	switch accountType := id.AccountType(); accountType.Number {
	case 0: // Invalid
		return "[I:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case 1: // Individual
		if accountInstance := id.AccountInstace(); accountInstance == 1 { // Desktop instance
			return "[U:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		} else {
			return "[U:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + ":" + strconv.FormatUint(uint64(accountInstance), 10) + "]"
		}
	case 2: // Multiseat
		return "[M:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + ":" + strconv.FormatUint(uint64(id.AccountInstace()), 10) + "]"
	case 3: // GameServer
		return "[G:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case 4: // AnonGameServer
		return "[A:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + ":" + strconv.FormatUint(uint64(id.AccountInstace()), 10) + "]"
	case 5: // Pending
		return "[P:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case 6: // ContentServer
		return "[C:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case 7: // Clan
		return "[g:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
	case 8: // Chat
		const instanceMask = 0x000FFFFF

		const (
			clanFlag             = (instanceMask + 1) >> 1
			lobbyFlag            = (instanceMask + 1) >> 2
			matchmakingLobbyFlag = (instanceMask + 1) >> 3
		)

		if accountInstance := id.AccountInstace(); accountInstance&clanFlag == clanFlag {
			return "[c:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		} else if accountInstance&lobbyFlag == lobbyFlag {
			return "[L:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		} else if accountInstance&matchmakingLobbyFlag == matchmakingLobbyFlag {
			return "[T:" + strconv.FormatUint(uint64(id.Universe()), 10) + ":" + strconv.FormatUint(uint64(id.AccountID()), 10) + "]"
		}
	}

	return "Unknown type"
}

// Creates a URL to this SteamID. If this type of SteamID doesn't support URLs, returns an empty string.
func (id SteamID) URL() string {
	switch accountType := id.AccountType(); accountType.Number {
	case 1:
		accountID := id.AccountID()
		return "https://steamcommunity.com/profiles/[U:1:" + strconv.FormatUint(uint64((accountID>>1)*2+accountID&1), 10) + "]"
	case 7:
		accountID := id.AccountID()
		return "https://steamcommunity.com/gid/[g:1:" + strconv.FormatUint(uint64((accountID>>1)*2+accountID&1), 10) + "]"
	default:
		return ""
	}
}

// An alias of SteamID3
func (id SteamID) String() string {
	return id.SteamID3()
}

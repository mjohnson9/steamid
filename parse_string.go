package steamid

import (
	"errors"
	"strconv"
	"strings"
)

// ErrInvalidSteamID is returned when an attempt is made to parse an invalid
// SteamID.
var ErrInvalidSteamID = errors.New("invalid SteamID")

var invalidSteamID = FromValues(UniverseUnspecified, 0, 0, 0)

// Parse attempts to parse a SteamID from a string via automatic detection of
// the string's format. If it fails, it will return an invalid SteamID and
// ErrInvalidSteamID.
func Parse(steamID string) (SteamID, error) {
	if steamIDUpper := strings.ToUpper(steamID); steamIDUpper == "UNKNOWN" || strings.HasPrefix(steamIDUpper, "STEAM_") {
		return ParseV2(steamIDUpper)
	}

	if strings.HasPrefix(steamID, "[") && strings.HasSuffix(steamID, "]") {
		return ParseV3(steamID)
	}

	return invalidSteamID, ErrInvalidSteamID
}

// ParseV2 attempts to parse a SteamID from the version 2 textual
// representation. For example, STEAM_0:0:1.
func ParseV2(steamID string) (SteamID, error) {
	steamID = strings.ToUpper(steamID)

	if steamID == "STEAM_ID_PENDING" {
		return FromValues(UniverseUnspecified, 0, 5, 0), nil
	} else if steamID == "UNKNOWN" {
		return FromValues(UniverseUnspecified, 0, 0, 0), nil
	}

	if !strings.HasPrefix(steamID, "STEAM_") {
		return invalidSteamID, ErrInvalidSteamID
	}
	steamID = steamID[6:]

	split := strings.Split(steamID, ":")
	if len(split) != 3 {
		return invalidSteamID, ErrInvalidSteamID
	}

	universe, err := strconv.ParseInt(split[0], 10, 8)
	if err != nil {
		return invalidSteamID, err
	}

	authServer, err := strconv.ParseUint(split[1], 10, 1)
	if err != nil {
		return invalidSteamID, err
	}

	accountID, err := strconv.ParseUint(split[2], 10, 31)
	if err != nil {
		return invalidSteamID, err
	}

	return FromValues(uint8(universe), 1, 1, (uint32(accountID)<<1)|uint32(authServer)), nil
}

// ParseV3 attempts to parse a SteamID from the version 3 textual
// representation. For example, [U:0:2].
func ParseV3(steamID string) (SteamID, error) {
	steamID = strings.TrimPrefix(steamID, "[")
	steamID = strings.TrimSuffix(steamID, "]")

	split := strings.Split(steamID, ":")
	if len(split) < 3 || len(split) > 4 {
		return invalidSteamID, ErrInvalidSteamID
	}

	idType := split[0]
	if idType != "U" && len(split) != 3 {
		return invalidSteamID, ErrInvalidSteamID
	}

	universe, err := strconv.ParseInt(split[1], 10, 8)
	if err != nil {
		return invalidSteamID, err
	}

	accountID, err := strconv.ParseInt(split[2], 10, 32)
	if err != nil {
		return invalidSteamID, err
	}

	accountTypeNum := uint8(255)

	switch idType {
	case "I":
		accountTypeNum = AccountTypeInvalid

	case "U":
		if len(split) == 4 {
			accountInstance, err := strconv.ParseInt(split[3], 10, 20)
			if err != nil {
				return invalidSteamID, err
			}

			return FromValues(uint8(universe), uint32(accountInstance), AccountTypeIndividual, uint32(accountID)), nil
		}
		return FromValues(uint8(universe), 1, AccountTypeIndividual, uint32(accountID)), nil

	case "M":
		accountTypeNum = AccountTypeMultiseat

	case "G":
		accountTypeNum = AccountTypeGameServer

	case "A":
		accountTypeNum = AccountTypeAnonGameServer

	case "P":
		accountTypeNum = AccountTypePending

	case "C":
		accountTypeNum = AccountTypeContentServer

	case "g":
		accountTypeNum = AccountTypeClan

	case "c", "L", "T":
		var accountInstance uint32

		const (
			instanceMask         = 0x000FFFFF
			clanFlag             = (instanceMask + 1) >> 1
			lobbyFlag            = (instanceMask + 1) >> 2
			matchmakingLobbyFlag = (instanceMask + 1) >> 3
		)

		switch idType {
		case "c":
			accountInstance = clanFlag
		case "L":
			accountInstance = lobbyFlag
		case "T":
			accountInstance = matchmakingLobbyFlag
		default:
			return invalidSteamID, ErrInvalidSteamID
		}

		return FromValues(uint8(universe), accountInstance, AccountTypeChat, uint32(accountID)), nil
	}

	if accountTypeNum == 255 {
		return invalidSteamID, ErrInvalidSteamID
	}
	return FromValues(uint8(universe), 0, accountTypeNum, uint32(accountID)), nil
}

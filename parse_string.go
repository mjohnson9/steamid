package steamid

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var steamID2Regex = regexp.MustCompile("^STEAM_(\\d+):(\\d+):(\\d+)$")

// ErrInvalidSteamID is returned when an attempt is made to parse an invalid
// SteamID.
var ErrInvalidSteamID = errors.New("invalid SteamID")

var invalidSteamID = FromValues(UniverseUnspecified, 0, 0, 0)

// FromString attempts to parse a SteamID from a string. If it fails, it will
// return an invalid SteamID and ErrInvalidSteamID.
func FromString(steamID string) (SteamID, error) {
	steamID = strings.ToUpper(steamID)
	if steamID == "STEAM_ID_PENDING" {
		return FromValues(UniverseUnspecified, 0, 5, 0), nil
	}
	if steamID == "UNKNOWN" {
		return FromValues(UniverseUnspecified, 0, 0, 0), nil
	}

	if strings.HasPrefix(steamID, "STEAM_") {
		return FromSteamID2(steamID)
	}

	return invalidSteamID, ErrInvalidSteamID
}

// FromSteamID2 attempts to parse a SteamID from the version 2 textual
// representation. For example, STEAM_0:0:1.
func FromSteamID2(steamID string) (SteamID, error) {
	matches := steamID2Regex.FindStringSubmatch(steamID)
	if len(matches) < 4 {
		return invalidSteamID, ErrInvalidSteamID
	}

	var err error

	universe64, err := strconv.ParseInt(matches[1], 10, 8)
	if err != nil {
		return invalidSteamID, err
	}
	universe := uint8(universe64)

	authServer64, err := strconv.ParseUint(matches[2], 10, 2)
	if err != nil {
		return invalidSteamID, err
	}
	authServer := uint8(authServer64)

	accountID64, err := strconv.ParseUint(matches[3], 10, 31)
	if err != nil {
		return invalidSteamID, err
	}
	accountID := uint32(accountID64)

	return FromValues(universe, 1, 1, (accountID<<1)|uint32(authServer)), nil
}

package steamid

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var steamID2Regex = regexp.MustCompile("^STEAM_(\\d+):(\\d+):(\\d+)$")

var InvalidSteamIDError = errors.New("Invalid SteamID")

var invalidSteamID = FromValues(UniverseUnspecified, 0, AccountTypes[0], 0)

// Attempts to automatically parse a SteamID from a string using various methods
func FromString(steamID string) (SteamID, error) {
	if steamID == "STEAM_ID_PENDING" {
		return FromValues(UniverseUnspecified, 0, AccountTypes[5], 0), nil
	}

	if strings.HasPrefix(steamID, "STEAM_") {
		return FromSteamID2(steamID)
	}

	return invalidSteamID, InvalidSteamIDError
}

// Attempts to parse a SteamID from a verison 2 textual repesentation (i.e.: STEAM_0:0:11101)
func FromSteamID2(steamID string) (SteamID, error) {
	matches := steamID2Regex.FindStringSubmatch(steamID)
	if len(matches) < 4 {
		return invalidSteamID, InvalidSteamIDError
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

	return FromValues(universe, 1, AccountTypes[1], (accountID<<1)|uint32(authServer)), nil
}

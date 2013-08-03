package steamid

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var steamID2Regex = regexp.MustCompile("^STEAM_(\\d+):(\\d+):(\\d+)$")

var InvalidSteamIDError = errors.New("Invalid SteamID")

var invalidSteamID = FromValues(UniverseUnspecified, 0, 0, AccountTypes[0], 0)

func FromString(steamID string) (SteamID, error) {
	if steamID == "STEAM_ID_PENDING" {
		return FromValues(UniverseUnspecified, 0, 0, AccountTypes[5], 0), nil
	}

	if strings.HasPrefix(steamID, "STEAM_") {
		return FromSteamID2(steamID)
	}

	return invalidSteamID, InvalidSteamIDError
}

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

	authServer64, err := strconv.ParseInt(matches[2], 10, 2)
	if err != nil {
		return invalidSteamID, err
	}
	authServer := uint8(authServer64)

	accountID64, err := strconv.ParseInt(matches[3], 10, 31)
	if err != nil {
		return invalidSteamID, err
	}
	accountID := uint32(accountID64)

	return FromValues(universe, authServer, 1, AccountTypes[1], accountID), nil
}

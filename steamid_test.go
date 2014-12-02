package steamid_test

import (
	"testing"

	"github.com/nightexcessive/steamid"
)

var steamIDTests = []struct {
	Universe    uint8
	Instance    uint32
	Type        uint8
	ID          uint32
	ExpectedID2 string
	ExpectedID3 string
}{
	{steamid.UniverseUnspecified, 1, steamid.AccountTypeIndividual, 2, "STEAM_0:0:1", "[U:0:2]"},
	{steamid.UniversePublic, 1, steamid.AccountTypeIndividual, 35812865, "STEAM_0:1:17906432", "[U:1:35812865]"},
}

func TestSteamID(t *testing.T) {
	for _, details := range steamIDTests {
		steamID := steamid.FromValues(details.Universe, details.Instance, details.Type, details.ID)

		got2 := steamID.SteamID2()
		if got2 != details.ExpectedID2 {
			t.Errorf("SteamID2([Universe %d Instance %d Type %d ID %d]): got %q, expected %q", details.Universe, details.Instance, details.Type, details.ID, got2, details.ExpectedID2)
		}

		got3 := steamID.SteamID3()
		if got3 != details.ExpectedID3 {
			t.Errorf("SteamID3([Universe %d Instance %d Type %d ID %d]): got %q, expected %q", details.Universe, details.Instance, details.Type, details.ID, got3, details.ExpectedID3)
		}
	}
}

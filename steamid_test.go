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
	{steamid.UniverseUnspecified, 0, steamid.AccountTypeInvalid, 0, "UNKNOWN", "[I:0:0]"},
	{steamid.UniverseUnspecified, 1, steamid.AccountTypeIndividual, 2, "STEAM_0:0:1", "[U:0:2]"},
	{steamid.UniverseUnspecified, 2, steamid.AccountTypeIndividual, 2, "STEAM_0:0:1", "[U:0:2:2]"},
	{steamid.UniverseUnspecified, 1, steamid.AccountTypeIndividual, 35812865, "STEAM_0:1:17906432", "[U:0:35812865]"},
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

		gotMarshal, err := steamID.MarshalText()
		if err != nil {
			t.Errorf("MarshalText([Universe %d Instance %d Type %d ID %d]): %s", details.Universe, details.Instance, details.Type, details.ID, err)
		} else if string(gotMarshal) != details.ExpectedID2 && string(gotMarshal) != details.ExpectedID3 {
			t.Errorf("MarshalText([Universe %d Instance %d Type %d ID %d]): got %q, expected %q or %q", details.Universe, details.Instance, details.Type, details.ID, string(gotMarshal), details.ExpectedID2, details.ExpectedID3)
		}
	}
}

func TestParseID(t *testing.T) {
	for _, details := range steamIDTests {
		_, err := steamid.Parse(details.ExpectedID3)
		if err != nil {
			t.Errorf("Parse(%q): error parsing: %s", details.ExpectedID3, err)
		}

		_, err = steamid.Parse(details.ExpectedID2)
		if err != nil {
			t.Errorf("Parse(%q): error parsing: %s", details.ExpectedID2, err)
		}
	}
}

func TestParseID2(t *testing.T) {
	for _, details := range steamIDTests {
		steamID, err := steamid.ParseV2(details.ExpectedID2)
		if err != nil {
			t.Errorf("ParseV2(%q): error parsing: %s", details.ExpectedID2, err)
			continue
		}

		if u := steamID.Universe(); u != details.Universe {
			t.Errorf("Universe(%q): got %d, expected %d", details.ExpectedID2, u, details.Universe)
		}

		typ, _ := steamID.AccountType()

		if inst := steamID.AccountInstance(); inst != details.Instance &&
			!(typ == steamid.AccountTypeIndividual && details.Instance != 1) { // Version 2 IDs can't carry instance information
			t.Errorf("AccountInstance(%q): got %d, expected %d", details.ExpectedID2, inst, details.Instance)
		}

		if typ != int(details.Type) {
			t.Errorf("AccountType(%q): got %d, expected %d", details.ExpectedID2, typ, details.Type)
		}

		if id := steamID.AccountID(); id != details.ID {
			t.Errorf("AccountID(%q): got %d, expected %d", details.ExpectedID2, id, details.ID)
		}
	}
}

func TestParseID3(t *testing.T) {
	for _, details := range steamIDTests {
		steamID, err := steamid.ParseV3(details.ExpectedID3)
		if err != nil {
			t.Errorf("ParseV3(%q): error parsing: %s", details.ExpectedID3, err)
			continue
		}

		if u := steamID.Universe(); u != details.Universe {
			t.Errorf("Universe(%q): got %d, expected %d", details.ExpectedID3, u, details.Universe)
		}

		if inst := steamID.AccountInstance(); inst != details.Instance {
			t.Errorf("AccountInstance(%q): got %d, expected %d", details.ExpectedID3, inst, details.Instance)
		}

		if typ, _ := steamID.AccountType(); typ != int(details.Type) {
			t.Errorf("AccountType(%q): got %d, expected %d", details.ExpectedID3, typ, details.Type)
		}

		if id := steamID.AccountID(); id != details.ID {
			t.Errorf("AccountID(%q): got %d, expected %d", details.ExpectedID3, id, details.ID)
		}
	}
}

func BenchmarkParseID2(b *testing.B) {
	const testID = "STEAM_0:0:1"

	for i := 0; i < b.N; i++ {
		if _, err := steamid.ParseV2(testID); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseID3(b *testing.B) {
	const testID = "[U:0:2]"

	for i := 0; i < b.N; i++ {
		if _, err := steamid.ParseV3(testID); err != nil {
			b.Fatal(err)
		}
	}
}

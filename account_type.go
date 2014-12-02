package steamid

// accountType contains information about a given account type.
type accountType struct {
	Number uint8
	Name   string

	Modifier uint64
}

// A list of account types.
//
// See
// https://developer.valvesoftware.com/wiki/SteamID#Types_of_Steam_Accounts
// for more information.
var accountTypes = []*accountType{
	&accountType{
		Number: 0,
		Name:   "Invalid",
	},
	&accountType{
		Number: 1,
		Name:   "Individual",

		Modifier: 0x0110000100000000,
	},
	&accountType{
		Number: 2,
		Name:   "Multiseat",
	},
	&accountType{
		Number: 3,
		Name:   "GameServer",
	},
	&accountType{
		Number: 4,
		Name:   "AnonGameServer",
	},
	&accountType{
		Number: 5,
		Name:   "Pending",
	},
	&accountType{
		Number: 6,
		Name:   "ContentServer",
	},
	&accountType{
		Number: 7,
		Name:   "Clan",

		Modifier: 0x0170000000000000,
	},
	&accountType{
		Number: 8,
		Name:   "Chat",
	},
	&accountType{
		Number: 9,
		Name:   "P2P SuperSeeder",
	},
	&accountType{
		Number: 10,
		Name:   "AnonUser",
	},
}

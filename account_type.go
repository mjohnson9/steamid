package steamid

type AccountType struct {
	Number uint8
	Name   string

	Modifier uint64
}

var AccountTypes = [...]*AccountType{
	&AccountType{
		Number: 0,
		Name:   "Invalid",
	},
	&AccountType{
		Number: 1,
		Name:   "Individual",

		Modifier: 0x0110000100000000,
	},
	&AccountType{
		Number: 2,
		Name:   "Multiseat",
	},
	&AccountType{
		Number: 3,
		Name:   "GameServer",
	},
	&AccountType{
		Number: 4,
		Name:   "AnonGameServer",
	},
	&AccountType{
		Number: 5,
		Name:   "Pending",
	},
	&AccountType{
		Number: 6,
		Name:   "ContentServer",
	},
	&AccountType{
		Number: 7,
		Name:   "Clan",

		Modifier: 0x0170000000000000,
	},
	&AccountType{
		Number: 8,
		Name:   "Chat",
	},
	&AccountType{
		Number: 9,
		Name:   "P2P SuperSeeder",
	},
	&AccountType{
		Number: 10,
		Name:   "AnonUser",
	},
}

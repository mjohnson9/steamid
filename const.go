package steamid

// Universes
//
// See
// https://developer.valvesoftware.com/wiki/SteamID#Universes_Available_for_Steam_Accounts
// for more information.
const (
	UniverseUnspecified = 0
	UniversePublic      = 1
	UniverseBeta        = 2
	UniverseInternal    = 3
	UniverseDev         = 4
	UniverseRC          = 5
)

// Account types
//
// See
// https://developer.valvesoftware.com/wiki/SteamID#Types_of_Steam_Accounts for
// more information.
const (
	AccountTypeInvalid        = 0
	AccountTypeIndividual     = 1
	AccountTypeMultiseat      = 2
	AccountTypeGameServer     = 3
	AccountTypeAnonGameServer = 4
	AccountTypePending        = 5
	AccountTypeContentServer  = 6
	AccountTypeClan           = 7
	AccountTypeChat           = 8
	AccountTypeP2PSuperSeeder = 9
	AccountTypeAnonUser       = 10
)

package steamid

// The list of known account types
//
// See
// https://developer.valvesoftware.com/wiki/SteamID#Types_of_Steam_Accounts for
// more information.
const (
	AccountTypeInvalid        AccountType = 0
	AccountTypeIndividual     AccountType = 1
	AccountTypeMultiseat      AccountType = 2
	AccountTypeGameServer     AccountType = 3
	AccountTypeAnonGameServer AccountType = 4
	AccountTypePending        AccountType = 5
	AccountTypeContentServer  AccountType = 6
	AccountTypeClan           AccountType = 7
	AccountTypeChat           AccountType = 8
	AccountTypeP2PSuperSeeder AccountType = 9
	AccountTypeAnonUser       AccountType = 10
)

// AccountType represents a Steam account type.
//
// For more information about account types, see
// https://developer.valvesoftware.com/wiki/SteamID#Types_of_Steam_Accounts
type AccountType uint8

// Modifier returns the account type's modifier. If the account type does not
// have a modifier, it returns 0.
//
// Currently the only account types with modifiers are AccountTypeIndividual and
// AccountTypeClan.
func (a AccountType) Modifier() uint64 {
	switch a {
	case AccountTypeIndividual:
		return 0x0110000100000000
	case AccountTypeClan:
		return 0x0170000000000000
	default:
		return 0
	}
}

// String gives the name of this account type. If the name of this account type
// isn't known, it returns an empty string.
func (a AccountType) String() string {
	switch a {
	case AccountTypeInvalid:
		return "Invalid"
	case AccountTypeIndividual:
		return "Individual"
	case AccountTypeMultiseat:
		return "Multiseat"
	case AccountTypeGameServer:
		return "GameServer"
	case AccountTypeAnonGameServer:
		return "AnonGameServer"
	case AccountTypePending:
		return "Pending"
	case AccountTypeContentServer:
		return "ContentServer"
	case AccountTypeClan:
		return "Clan"
	case AccountTypeChat:
		return "Chat"
	case AccountTypeP2PSuperSeeder:
		return "P2P SuperSeeder"
	case AccountTypeAnonUser:
		return "AnonUser"
	default:
		return ""
	}
}

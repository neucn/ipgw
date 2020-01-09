package api

var (
	globalMissArgs      = "0 0"
	globalNotCompatible = "0 -1"
	globalFailLoad      = "0 -2"
	globalNetError      = "0 -3"

	loginNoPassword    = "1 1"
	loginNoStoredUser  = "1 2"
	loginWrongUP       = "1 3"
	loginCookieExpired = "1 4"
	loginBanned        = "1 5"
	loginFail          = "1 6"
)

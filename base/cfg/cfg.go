package cfg

// These are flags used by many commands.
var (
	Version string = "v0.1.0"

	CmdName string // "build", "install", "list", "mod tidy", etc.

	FullView bool // print details
)

package configuration

type CommandLineOption struct {
	Shorthand string
	Name      string
}

type Key struct {
	Name        string
	Description string
	CommandLine CommandLineOption
}

//nolint:gochecknoglobals // Const keys used as enums are acceptable
var (
	KeyColorizedLogging = Key{
		Name:        "log.color",
		Description: "Formatted & colorized logging on stdout",
		CommandLine: CommandLineOption{
			Name: "color",
		},
	}
	KeyAddress = Key{
		Name:        "server.address",
		Description: "address to listen to.",
		CommandLine: CommandLineOption{
			Shorthand: "a",
			Name:      "address",
		},
	}
	KeyPort = Key{
		Name:        "server.port",
		Description: "port to listen to.",
		CommandLine: CommandLineOption{
			Shorthand: "p",
			Name:      "port",
		},
	}
)

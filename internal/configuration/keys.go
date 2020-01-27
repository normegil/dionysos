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
	KeyDatabaseAddress = Key{
		Name:        "database.address",
		Description: "address to connect to database.",
		CommandLine: CommandLineOption{
			Name: "dbhost",
		},
	}
	KeyDatabasePort = Key{
		Name:        "database.port",
		Description: "port to connect to database.",
		CommandLine: CommandLineOption{
			Name: "dbport",
		},
	}
	KeyDatabaseUser = Key{
		Name:        "database.user",
		Description: "user to connect to database.",
		CommandLine: CommandLineOption{
			Name: "dbuser",
		},
	}
	KeyDatabasePassword = Key{
		Name:        "database.password",
		Description: "password to connect to database.",
		CommandLine: CommandLineOption{
			Name: "dbpassword",
		},
	}
	KeyDatabaseName = Key{
		Name:        "database.name",
		Description: "database name to use. The database will be created if required.",
		CommandLine: CommandLineOption{
			Name: "dbname",
		},
	}
)

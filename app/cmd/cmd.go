// Package cmd has all top-level commands dispatched by main's flag.Parse
// The entry point of each command is Execute function
package cmd

// CommonOptionsCommander extends flags.Commander with SetCommon
// All commands should implement this interfaces
type CommonOptionsCommander interface {
	SetCommon(commonOpts CommonOpts)
	Execute(args []string) error
}

// CommonOpts sets externally from main, shared across all commands
type CommonOpts struct {
	Revision   string
	DbPort     int    `long:"db_port" env:"DB_PORT" default:"5432" description:"port for database"`
	DbHost     string `long:"db_host" env:"DB_HOST" default:"localhost" description:"host for database"`
	DbUser     string `long:"db_user" env:"DB_USER" default:"postgres" description:"user for database"`
	DbPassword string `long:"db_password" env:"DB_PASSWORD" default:"password" description:"password for database"`
	DbName     string `long:"db_name" env:"DB_NAME" default:"postgres" description:"database name"`
}

// SetCommon satisfies CommonOptionsCommander interface and sets common option fields
// The method called by main for each command
func (c *CommonOpts) SetCommon(commonOpts CommonOpts) {
	c.Revision = commonOpts.Revision
}

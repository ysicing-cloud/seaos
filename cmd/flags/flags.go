package flags

import (
	flag "github.com/spf13/pflag"
)

// GlobalFlags is the flags that contains the global flags
type GlobalFlags struct {
	Debug      bool
	Silent     bool
	ConfigPath string
	Vars       []string
	Flags      *flag.FlagSet
}

// SetGlobalFlags applies the global flags
func SetGlobalFlags(flags *flag.FlagSet) *GlobalFlags {
	globalFlags := &GlobalFlags{
		Vars:  []string{},
		Flags: flags,
	}
	flags.BoolVar(&globalFlags.Debug, "debug", false, "Prints the stack trace if an error occurs")
	flags.BoolVar(&globalFlags.Silent, "silent", false, "Run in silent mode and prevents any qcadmin log output except panics & fatals")
	flags.StringVar(&globalFlags.ConfigPath, "config", "", "The qcadmin config file to use")
	return globalFlags
}

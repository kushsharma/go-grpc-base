package cmd

import (
	"github.com/spf13/cobra"
)

type Config struct {
	ServerPort int
}

// InitCommands initializes application cli interface
func InitCommands(appname, version string) *cobra.Command {
	conf := Config{
		ServerPort: 11000,
	}

	rootCmd := &cobra.Command{
		Use:     appname,
		Version: version,
	}
	rootCmd.AddCommand(initClient(conf))
	rootCmd.AddCommand(initServer(conf))
	return rootCmd
}

package cmd

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	appName = "awsctl"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "display version information",
	RunE:  runVersion,
}

var (
	version string
	date    string
	commit  string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(os.Stdout, "%s - %s@%s [%s].[%s].[%s]\n",
		appName,
		version,
		commit,
		date,
		runtime.GOOS,
		runtime.GOARCH,
	)
	return nil
}

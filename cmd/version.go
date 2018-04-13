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
	build   string
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) error {
	fmt.Fprintf(os.Stdout, "%s - version [%s]@[%s]  [%s-%s]\n",
		appName,
		version,
		build,
		runtime.GOOS,
		runtime.GOARCH,
	)
	return nil
}

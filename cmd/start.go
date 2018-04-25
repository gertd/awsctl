package cmd

import (
	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start instances",
	RunE:  runStart,
}

func init() {
	RootCmd.AddCommand(startCmd)
	startCmd.Flags().StringVarP(&name, "name", "n", defName, "AWS instance name")
	startCmd.Flags().StringVarP(&instanceID, "instance-id", "i", defInstanceID, "AWS instance ID")
	startCmd.Flags().BoolVar(&all, "all", defAll, "start all instances: default false")
}

func runStart(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region, profile, cmdLineCreds())
	defer client.Close()

	f := shared.NewInstanceFilter()

	if instanceID != defInstanceID {
		f.Add("instance-id", instanceID)
	} else if name != defName {
		f.Add("tag:Name", name)
	}

	client.StartInstance(f.Get())
	client.GetInstances(f.Get()).Print()

	return nil
}

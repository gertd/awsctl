package cmd

import (
	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop instances",
	RunE:  runStop,
}

func init() {
	RootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringVarP(&name, "name", "n", "", "AWS instance name")
	stopCmd.Flags().StringVarP(&instanceID, "instance-id", "i", "", "AWS instance ID")
	stopCmd.Flags().BoolVar(&all, "all", false, "stop all instances: default false")
}

func runStop(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region)
	defer client.Close()

	f := shared.NewInstanceFilter()

	if instanceID != defInstanceID {
		f.Add("instance-id", instanceID)
	} else if name != defName {
		f.Add("tag:Name", name)
	}

	client.StopInstance(f.Get())
	client.GetInstances(f.Get()).Print()

	return nil
}

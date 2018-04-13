package cmd

import (
	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list instance status",
	RunE:  runList,
}

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&name, "name", "n", defName, "AWS instance name")
	listCmd.Flags().StringVarP(&instanceID, "instance-id", "i", defInstanceID, "AWS instance ID")
}

func runList(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region)
	defer client.Close()

	f := shared.NewInstanceFilter()

	if instanceID != defInstanceID {
		f.Add("instance-id", instanceID)
	} else if name != defName {
		f.Add("tag:Name", name)
	}

	client.GetInstances(f.Get()).Print()

	return nil
}

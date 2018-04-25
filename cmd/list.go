package cmd

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list instance status",
	RunE:  runList,
}

const (
	defWindows = false
)

var (
	windows bool
)

func init() {
	RootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVarP(&name, "name", "n", defName, "AWS instance name")
	listCmd.Flags().StringVarP(&instanceID, "instance-id", "i", defInstanceID, "AWS instance ID")
	listCmd.Flags().BoolVar(&windows, "windows", defWindows, "Only return Windows nodes")
}

func runList(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region, profile, cmdLineCreds())
	defer client.Close()

	f := shared.NewInstanceFilter()

	if instanceID != defInstanceID {
		f.Add("instance-id", instanceID)
	} else if name != defName {
		f.Add("tag:Name", name)
	}

	var plaformFilter shared.PlaformFilterFunc = shared.All
	if windows {
		plaformFilter = shared.Windows
	}

	resolver := endpoints.DefaultResolver()
	partitions := filter(resolver.(endpoints.EnumPartitions).Partitions(), "aws")

	regions := []string{}

	for _, p := range partitions {

		if len(region) > 0 {
			if r, ok := p.Regions()[region]; ok {
				regions = append(regions, r.ID())
			}
		} else {
			for r := range p.Regions() {
				regions = append(regions, r)
			}
		}

		shared.Sort(
			len(regions),
			func(i, j int) {
				temp := regions[i]
				regions[i] = regions[j]
				regions[j] = temp
			},
			func(i, j int) bool {
				return regions[i] < regions[j]
			})

		shared.InstanceHeader()

		for _, region := range regions {
			client.SetRegion(region)

			client.GetInstances(f.Get()).
				PlatformFilter(client, plaformFilter).
				Print()
		}
	}

	return nil
}

func filter(partitions []endpoints.Partition, filter string) []endpoints.Partition {
	vps := make([]endpoints.Partition, 0)
	for _, v := range partitions {
		if strings.Compare(v.ID(), filter) == 0 {
			vps = append(vps, v)
		}
	}
	return vps
}

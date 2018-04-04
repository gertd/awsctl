package shared

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// TagFilter -
func TagFilter(tags []*ec2.Tag, key string) string {
	for _, v := range tags {
		if *v.Key == key {
			return *v.Value
		}
	}
	return ""
}

// InstanceFilter -
func InstanceFilter(filters ...Filter) []*ec2.Filter {

	instanceFilter := make([]*ec2.Filter, 0)
	for _, filter := range filters {
		instanceFilter = append(instanceFilter, &ec2.Filter{
			Name:   aws.String(filter.Name),
			Values: []*string{aws.String(filter.Value)},
		})
	}

	if len(instanceFilter) == 0 {
		instanceFilter = append(instanceFilter, &ec2.Filter{})
	}
	return instanceFilter
}

// Filter -
type Filter struct {
	Name  string
	Value string
}

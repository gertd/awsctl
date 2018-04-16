package shared

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// InstanceFilter --
type InstanceFilter struct {
	filterMap map[string][]*string
}

// NewInstanceFilter --
func NewInstanceFilter() *InstanceFilter {
	instanceFilter := InstanceFilter{
		filterMap: make(map[string][]*string),
	}
	return &instanceFilter
}

// Add --
func (i *InstanceFilter) Add(name, value string) {

	v, ok := i.filterMap[name]

	if ok {
		// ok == filter label exists in map, add fiter value to array
		v = append(v, &value)

	} else {
		// !ok == filter label does exists in map, add key and add fiter value to new array
		v = []*string{&value}
	}
	i.filterMap[name] = v
}

// Get --
func (i *InstanceFilter) Get() []*ec2.Filter {

	instanceFilter := []*ec2.Filter{}

	for k, v := range i.filterMap {
		instanceFilter = append(instanceFilter, &ec2.Filter{
			Name:   aws.String(k),
			Values: v,
		})
	}

	if len(instanceFilter) == 0 {
		instanceFilter = append(instanceFilter, &ec2.Filter{})
	}

	return instanceFilter
}

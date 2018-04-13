package shared

import (
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

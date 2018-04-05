package shared

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2Client --
type EC2Client struct {
	client *ec2.EC2
}

// NewEC2Client --
func NewEC2Client(region string) *EC2Client {

	client := EC2Client{}

	opts := session.Options{SharedConfigState: session.SharedConfigEnable}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	if len(region) > 0 {
		sess.Config.Region = aws.String(region)
	}

	client.client = ec2.New(sess)
	return &client
}

// Close --
func (c *EC2Client) Close() {

	c.client = nil
}

// GetInstances --
// TODO -- add ability to filter?
func (c *EC2Client) GetInstances() *Instances {

	describeInstanceInput := &ec2.DescribeInstancesInput{
		Filters: InstanceFilter(),
	}

	var err error
	insts, err := c.client.DescribeInstances(describeInstanceInput)
	if err != nil {
		log.Println(err)
		return nil
	}

	instances := Instances(*insts)

	return &instances
}

// StartInstance --
func (c *EC2Client) StartInstance(name string) {

	instanceFilters := []*ec2.Filter{}
	instanceFilters = append(instanceFilters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("stopped")},
	})

	if len(name) > 0 {
		instanceFilters = append(instanceFilters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String(name)},
		})
	}

	describeInstanceInput := &ec2.DescribeInstancesInput{
		Filters: instanceFilters,
	}

	instances, err := c.client.DescribeInstances(describeInstanceInput)
	if err != nil {
		log.Println(err)
		return
	}

	instancesToStart := []*string{}
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			instancesToStart = append(instancesToStart, instance.InstanceId)
		}
	}

	if len(instancesToStart) == 0 {
		return
	}

	fmt.Printf("starting stopped instances...\n")
	startResp, err := c.client.StartInstances(&ec2.StartInstancesInput{InstanceIds: instancesToStart})
	if err != nil {
		log.Println(err.Error())
	}

	for _, startingInstance := range startResp.StartingInstances {
		fmt.Printf("%s %s -> %s \n",
			*startingInstance.InstanceId,
			*startingInstance.PreviousState.Name,
			*startingInstance.CurrentState.Name,
		)
	}

	fmt.Printf("waiting for instances to be started...\n")
	err = c.client.WaitUntilInstanceRunning(&ec2.DescribeInstancesInput{InstanceIds: instancesToStart})
	if err != nil {
		log.Println(err)
	}

	return
}

// StopInstance --
func (c *EC2Client) StopInstance(name string) {

	instanceFilters := []*ec2.Filter{}
	instanceFilters = append(instanceFilters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("running")},
	})

	if len(name) > 0 {
		instanceFilters = append(instanceFilters, &ec2.Filter{
			Name:   aws.String("tag:Name"),
			Values: []*string{aws.String(name)},
		})
	}

	describeInstanceInput := &ec2.DescribeInstancesInput{
		Filters: instanceFilters,
	}

	instances, err := c.client.DescribeInstances(describeInstanceInput)
	if err != nil {
		log.Println(err)
		return
	}

	instancesToStop := []*string{}
	for _, reservation := range instances.Reservations {
		for _, instance := range reservation.Instances {
			instancesToStop = append(instancesToStop, instance.InstanceId)
		}
	}

	fmt.Printf("stopping running instances...\n")
	stopResp, err := c.client.StopInstances(&ec2.StopInstancesInput{InstanceIds: instancesToStop, Force: aws.Bool(true)})
	if err != nil {
		log.Println(err.Error())
	}
	for _, stoppingInstance := range stopResp.StoppingInstances {
		fmt.Printf("%s %s -> %s \n",
			*stoppingInstance.InstanceId,
			*stoppingInstance.PreviousState.Name,
			*stoppingInstance.CurrentState.Name,
		)
	}

	fmt.Printf("waiting for instances to be stopped...\n")
	err = c.client.WaitUntilInstanceStopped(&ec2.DescribeInstancesInput{InstanceIds: instancesToStop})
	if err != nil {
		log.Println(err)
	}

	return
}

// GetPasswordData --
func (c *EC2Client) GetPasswordData(instanceID string) string {

	input := ec2.GetPasswordDataInput{InstanceId: &instanceID}
	output, err := c.client.GetPasswordData(&input)
	if err != nil {
		log.Println(err)
	}
	return *output.PasswordData
}

// Instances --
type Instances ec2.DescribeInstancesOutput

// Print --
func (i *Instances) Print() {
	for _, v := range i.Reservations {

		fmt.Printf("%s - %-15s %-13s %s \n",
			*v.Instances[0].InstanceId,
			PStr(v.Instances[0].PublicIpAddress, "N/A"),
			*v.Instances[0].State.Name,
			TagFilter(v.Instances[0].Tags, "Name"),
		)
	}
}

// PStr -- convert string pointer to string, on nil pointer replacement value in first element or "" by default
func PStr(pstr *string, replacement ...string) string {
	if pstr == nil {
		if len(replacement) == 0 {
			return ""
		}
		return replacement[0]
	}
	return *pstr
}

package shared

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/service/ec2"
)

// EC2Client --
type EC2Client struct {
	session *session.Session
	client  *ec2.EC2
}

var (
	imageMap map[string]ec2.Image
	keyMap   map[string]ec2.KeyPairInfo
)

// NewEC2Client --
func NewEC2Client(region, profile string, cmdLineCreds func() credentials.Value) *EC2Client {

	client := EC2Client{}

	creds := credentials.NewChainCredentials(
		[]credentials.Provider{
			&credentials.SharedCredentialsProvider{},
			&credentials.EnvProvider{},
			&credentials.StaticProvider{Value: cmdLineCreds()},
		})

	config := aws.Config{
		Credentials:                   creds,
		CredentialsChainVerboseErrors: aws.Bool(true),
	}

	if len(region) > 0 {
		config.Region = aws.String(region)
	}

	opts := session.Options{
		Config:            config,
		SharedConfigState: session.SharedConfigEnable,
	}

	if len(profile) > 0 {
		opts.Profile = profile
	}

	sess, err := session.NewSessionWithOptions(opts)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	client.session = sess
	client.client = ec2.New(client.session)
	return &client
}

// Close --
func (c *EC2Client) Close() {
	c.session = nil
	c.client = nil
}

// SetRegion --
func (c *EC2Client) SetRegion(region string) {

	c.session.Config.Region = aws.String(region)
	c.client = ec2.New(c.session)
}

// GetInstances --
func (c *EC2Client) GetInstances(filters []*ec2.Filter) *Instances {

	describeInstanceInput := &ec2.DescribeInstancesInput{
		Filters: filters,
	}

	var err error
	out, err := c.client.DescribeInstances(describeInstanceInput)
	if err != nil {
		log.Println(err)
		return nil
	}

	instances := Instances{}
	for _, r := range out.Reservations {
		instances = append(instances, r.Instances...)
	}

	return &instances
}

// StartInstance --
func (c *EC2Client) StartInstance(filters []*ec2.Filter) {

	filters = append(filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("stopped")},
	})

	describeInstanceInput := &ec2.DescribeInstancesInput{
		Filters: filters,
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
}

// StopInstance --
func (c *EC2Client) StopInstance(filters []*ec2.Filter) {

	filters = append(filters, &ec2.Filter{
		Name:   aws.String("instance-state-name"),
		Values: []*string{aws.String("running")},
	})

	describeInstanceInput := &ec2.DescribeInstancesInput{
		Filters: filters,
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
type Instances []*ec2.Instance

// InstanceHeader --
func InstanceHeader() {

	fmt.Printf("%-19s - %-14s %-15s %-13s %s\n",
		"Instance ID",
		"Region",
		"Public IP Addr",
		"State",
		"Name",
	)
	fmt.Printf("%s\n", strings.Repeat("-", 79))
}

// Print --
func (in *Instances) Print() {
	if in == nil {
		return
	}

	for _, v := range *in {

		az := *v.Placement.AvailabilityZone

		fmt.Printf("%-19s - %-14s %-15s %-13s %s\n",
			*v.InstanceId,
			az[0:len(az)-1],
			PStr(v.PublicIpAddress, "N/A"),
			*v.State.Name,
			TagFilter(v.Tags, "Name"),
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

// GetImageByID --
func (c *EC2Client) GetImageByID(imageID *string) (*ec2.Image, error) {

	if imageID == nil || len(*imageID) == 0 {
		return nil, fmt.Errorf("GetImageByID(imageID) is nil or zero length")
	}

	if imageMap == nil {
		imageMap = make(map[string]ec2.Image)
	}

	// cache lookup
	i, ok := imageMap[*imageID]
	if ok {
		log.Printf("cache hit %v", *imageID)
		return &i, nil
	}

	describeImageInput := &ec2.DescribeImagesInput{ImageIds: []*string{imageID}}
	images, err := c.client.DescribeImages(describeImageInput)
	if err != nil {
		log.Printf("err %v", err)
		return nil, err
	}
	for _, v := range images.Images {
		imageMap[*v.ImageId] = *v
	}

	i, ok = imageMap[*imageID]
	if ok {
		return &i, nil
	}
	return nil, fmt.Errorf("image-id %s not found", *imageID)
}

// GetKeyPairInfoByName --
func (c *EC2Client) GetKeyPairInfoByName(keyName *string) (*ec2.KeyPairInfo, error) {

	if keyName == nil || len(*keyName) == 0 {
		return nil, fmt.Errorf("GetKeyPairInfoByName(keyName) is nil or zero length")
	}

	if keyMap == nil {
		keyMap = make(map[string]ec2.KeyPairInfo)
	}

	// cache lookup
	i, ok := keyMap[*keyName]
	if ok {
		log.Printf("cache hit %v", *keyName)
		return &i, nil
	}

	describeKeyPairsInput := ec2.DescribeKeyPairsInput{KeyNames: []*string{keyName}}
	keys, err := c.client.DescribeKeyPairs(&describeKeyPairsInput)
	if err != nil {
		log.Printf("err %v", err)
		return nil, err
	}

	for _, v := range keys.KeyPairs {
		keyMap[*v.KeyName] = *v
	}

	i, ok = keyMap[*keyName]
	if ok {
		return &i, nil
	}
	return nil, fmt.Errorf("key name %s not found", *keyName)
}

// PlaformFilterFunc --
type PlaformFilterFunc func(image *ec2.Image) bool

// Windows --
func Windows(image *ec2.Image) bool {
	return image != nil && image.Platform != nil && strings.ToLower(*image.Platform) == "windows"
}

// All --
func All(image *ec2.Image) bool { return true }

// PlatformFilter --
func (in *Instances) PlatformFilter(c *EC2Client, f PlaformFilterFunc) *Instances {

	instances := Instances{}

	for _, i := range *in {

		img, err := c.GetImageByID(i.ImageId)
		if err != nil {
			log.Print(err)
		}

		if f(img) {
			instances = append(instances, i)
		}
	}

	return &instances
}

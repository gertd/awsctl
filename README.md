# awsctl
Simple AWS command line for managing EC2 instances

## TL;DR;
Key capabilities are to quickly list, start, stop EC2 service instances based on name or instance-id and to retrieve the AWS Windows password using the key-pair used to create a Windows image

Currently this tools only supports using the AWS shared configuration state (located by default in ~/.aws)

## Installing AWSCTL
To install AWSCTL there are a couple of options:

Download the platform specific binary

Using go get (requires golang compiler 1.9.5 or higher to be installed)

	go get -x https://github.com/gertd/awsctl

## Help

	awsctl --help

	AWS instance manager
	
	Usage:
	  awsctl [command]
	
	Available Commands:
	  help        Help about any command
	  list        list running instances
	  pwd         get Windows password
	  start       start instances
	  stop        stop instances
	
	Flags:
	  -h, --help            help for awsctl
	      --region string   AWS region: default all
	
	Use "awsctl [command] --help" for more information about a command.

## List - list service status

	awsctl list
	i-04907c97017c96798 - N/A             stopped       jumpbox
	i-0908c2a42a46f67d6 - N/A             stopped
	i-0a1ef338556f8d2b9 - xxx.xxx.xxx.xxx running       win2016

## Start - start all or single service instance

### Start all 

	awsctl start
	starting stopped instances...
	i-0a1ef338556f8d2b9 stopped -> pending
	i-04907c97017c96798 stopped -> pending
	i-0908c2a42a46f67d6 stopped -> pending
	waiting for instances to be started...
	i-04907c97017c96798 - xxx.xxx.xxx.xxx running       jumpbox
	i-0908c2a42a46f67d6 - xxx.xxx.xxx.xxx running
	i-0a1ef338556f8d2b9 - xxx.xxx.xxx.xxx running       win2016

### Start single instance

	awsctl start --name jumpbox
	starting stopped instances...
	i-04907c97017c96798 stopped -> pending
	waiting for instances to be started...
	i-04907c97017c96798 - xxx.xxx.xxx.xxx running       jumpbox

## Stop - stop all or single service instance

### Stop all 
	awsctl stop
	stopping running instances...
	i-0a1ef338556f8d2b9 running -> stopping
	i-04907c97017c96798 running -> stopping
	i-0908c2a42a46f67d6 running -> stopping
	waiting for instances to be stopped...
	i-04907c97017c96798 - N/A             stopped       jumpbox
	i-0908c2a42a46f67d6 - N/A             stopped
	i-0a1ef338556f8d2b9 - N/A             stopped       win2016

### Stop single instance by name 

	awsctl stop --name jumpbox
	stopping running instances...
	i-04907c97017c96798 running -> stopping
	waiting for instances to be stopped...
	i-04907c97017c96798 - N/A             stopped       jumpbox

## Pwd - retrieve Windows password

Required parameters:

* --instance0d or --name 
* --keyfile

**NOTE**: If the keyfile contains a passphrase, the user will be prompted to provide the passphrase

Outputs:

* md5:  key file md5 fingerprint (used for AWS imported keys)
* sha1: key file sha1 fingerprint (used for AWS generated keys)
* pwd: **decrypted** password 

The fingerprints are provided to ease the correlation with the AWS KeyPair information, in order to identify the key used to creatd the Windows EC2 instance.


	awsctl pwd --instance-id i-0a1ef338556f8d2b9 --key-file ~/.ssh/id_rsa
	passphrase>>
	md5  bd:4d:40:e2:cf:8c:5e:c8:9a:c2:4a:02:ba:70:2b:a6
	sha1 a3:e6:21:6e:33:af:30:82:7f:95:9d:09:10:c6:ee:37:04:21:f3:bc
	pwd  xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx

## Todo List

Capabilities to add:

- [ ] command line login override (incl masked password entry)
- [ ] support for selecting configuration
- [ ] short parameter support for instance (-i) and name (-n)
- [ ] make start/stop all explicit using commandline parameter --all
- [ ] support list -instance-id or --name for singleton status
- [ ] add header to list output

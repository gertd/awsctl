[![Build Status](https://travis-ci.org/gertd/awsctl.svg?branch=dev)](https://travis-ci.org/gertd/awsctl)

# awsctl
Simple AWS command line for managing EC2 instances

## TL;DR;
Key capabilities are to quickly list, start, stop EC2 service instances based on name or instance-id and to retrieve the AWS Windows password using the key-pair used to create a Windows image

Currently this tools only supports using the AWS shared configuration state (by default located in ~/.aws)

## Installing AWSCTL
To install AWSCTL there are a couple of options:

Download the platform specific binary from the Github [release](https://github.com/gertd/awsctl/releases) page

### Linux

	curl -L https://github.com/gertd/awsctl/releases/download/v0.0.12/awsctl-linux-amd64 > ~/awsctl
	chmod +x ~/awsctl

### OSX

	curl -L https://github.com/gertd/awsctl/releases/download/v0.0.12/awsctl-darwin-amd64 > awsctl
	chmod +x ~/awsctl

### Windows (using PowerShell)

	[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12
	Invoke-WebRequest -Uri https://github.com/gertd/awsctl/releases/download/v0.0.12/awsctl-windows-amd64.exe -OutFile ~\awsctl.exe

### Go Get

Using [go get](https://golang.org/cmd/go/#hdr-Download_and_install_packages_and_dependencies), this requires the golang compiler 1.10.1 or higher to be installed. See [https://golang.org/dl/](https://golang.org/dl/)

	go get -x https://github.com/gertd/awsctl

Runing `go get` will install the process architecture specific binary in the $GOPATH/bin directory

## Help

	awsctl --help
	AWS instance manager
	
	Usage:
	  awsctl [command]
	
	Available Commands:
	  list        list instance status
	  pwd         get Windows password
	  start       start instances
	  stop        stop instances
	  version     display version information
	  help        Help about any command
	  
	Flags:
	      --access-key-id string       AWS AccessKeyID
	  -h, --help                       help for awsctl
	      --profile string             AWS profile
	      --region string              AWS region: default all
	      --secret-access-key string   AWS SecretAccessKey
	      		
	Use "awsctl [command] --help" for more information about a command.

## List - list instances in all regions

	awsctl list
	Instance ID         - Region         Public IP Addr  State         Name
	-------------------------------------------------------------------------------
	i-0908c2a42a46f67d6 - us-east-2      N/A             stopped
	i-04907c97017c96798 - us-west-2      xxx.xxx.xxx.xxx running       jumpbox
	i-0a1ef338556f8d2b9 - us-west-2      xxx.xxx.xxx.xxx running       win2016


## List - list Windows instances in region us-west-2

	awsctl list --windows --region us-west-2
	Instance ID         - Region         Public IP Addr  State         Name
	-------------------------------------------------------------------------------
	i-0a1ef338556f8d2b9 - us-west-2      xxx.xxx.xxx.xxx running       win2016

Parameters:

* --instance-id || --name (one of optional, defaults to all)
* --region
* --windows



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

Parameters:

* --instance-id || --name || --all (one of required)
* --region (optional)


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
	
Parameters:

* --instance-id || --name || --all (one of required)
* --region (optional)


## Pwd - retrieve Windows password

Parameters:

* --instance-id || --name (one of required) 
* --keyfile (required)

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


## Version - version info

	awsctl version
	awsctl - v0.0.12@eb97789 [2018-04-16T21:13:22+0000].[darwin].[amd64]

## Authentication

AWSCTL provides 3 ways to authenticate with AWS, in order of evaluation:

1. Shared configuration (~/.aws/config)
2. Environment variables (AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY and AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY)
3. Command line arguments --access-key-id and --secret-access-key

Command line example:

	awsctl list --access-key-id $MY_AWS_ACCESS_KEY_ID --secret-access-key $MY_AWS_SECRET_ACCESS_KEY


## Todo List

Capabilities to add:

- [X] command line login override (incl masked password entry)
- [X] support for selecting configuration
- [X] short parameter support for instance (-i) and name (-n)
- [X] make start/stop all explicit using commandline parameter --all
- [X] support list -instance-id or --name for singleton status
- [X] add header to list output
- [X] add error handling for when AWS config is not present
- [X] add windows filter to list command
- [ ] Proviode masked input prompt when --access-key-id or --secret-access-key empty

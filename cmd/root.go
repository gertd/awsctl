package cmd

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/spf13/cobra"
)

const (
	defName            = ""
	defInstanceID      = ""
	defAccessKeyID     = ""
	defSecretAccessKey = ""
	defRegion          = ""
	defProfile         = ""
	defKeyFile         = ""
	defAll             = false
)

var (
	region          string
	profile         string
	accessKeyID     string
	secretAccessKey string
	name            string
	keyFile         string
	instanceID      string
	all             bool
)

// RootCmd --
var RootCmd = &cobra.Command{
	Use:   "awsctl",
	Short: "AWS instance manager",
}

// Execute --
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	cobra.OnInitialize()
	cobra.EnableCommandSorting = false

	RootCmd.PersistentFlags().StringVar(&region, "region", defRegion, "AWS region: default all")
	RootCmd.PersistentFlags().StringVar(&accessKeyID, "access-key-id", defAccessKeyID, "AWS AccessKeyID")
	RootCmd.PersistentFlags().StringVar(&secretAccessKey, "secret-access-key", defSecretAccessKey, "AWS SecretAccessKey")
	RootCmd.PersistentFlags().StringVar(&profile, "profile", defProfile, "AWS profile")
}

func cmdLineCreds() func() credentials.Value {
	return func() credentials.Value {
		return credentials.Value{
			AccessKeyID:     accessKeyID,
			SecretAccessKey: secretAccessKey,
			SessionToken:    "",
		}
	}
}

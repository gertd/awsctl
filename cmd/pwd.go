package cmd

import (
	"encoding/base64"
	"fmt"
	"log"

	"github.com/gertd/awsctl/shared"
	"github.com/spf13/cobra"
)

var pwdCmd = &cobra.Command{
	Use:   "pwd",
	Short: "get Windows password",
	RunE:  runGetPwd,
}

func init() {
	RootCmd.AddCommand(pwdCmd)

	pwdCmd.Flags().StringVar(&instanceID, "instance-id", "", "AWS instance ID: default empty")
	pwdCmd.MarkFlagRequired("instance-id")

	pwdCmd.Flags().StringVar(&keyFile, "key-file", "", "key file")
	pwdCmd.MarkFlagRequired("key-file")
}

func runGetPwd(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region)
	defer client.Close()

	pwdData := client.GetPasswordData(instanceID)
	decData, err := base64.StdEncoding.DecodeString(pwdData)
	if err != nil {
		log.Println(err)
	}

	k, err := shared.ParseKeyFile(keyFile, shared.PassPhrasePrompt)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("md5 ", k.FingerPrintMD5())
	fmt.Println("sha1", k.FingerPrintSHA1())
	fmt.Println("pwd ", k.Decrypt(decData))

	return nil
}

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

	pwdCmd.Flags().StringVar(&instanceID, "instance-id", defInstanceID, "AWS instance ID: default empty")
	if err := pwdCmd.MarkFlagRequired("instance-id"); err != nil {
		log.Fatalln(err)
	}

	pwdCmd.Flags().StringVar(&keyFile, "key-file", defKeyFile, "key file")
	if err := pwdCmd.MarkFlagRequired("key-file"); err != nil {
		log.Fatalln(err)
	}
}

func runGetPwd(cmd *cobra.Command, args []string) error {

	client := shared.NewEC2Client(region, profile, cmdLineCreds())
	defer client.Close()

	pwdData := client.GetPasswordData(instanceID)
	decData, err := base64.StdEncoding.DecodeString(pwdData)
	if err != nil {
		return err
	}

	k, err := shared.ParseKeyFile(keyFile, shared.PassPhrasePrompt)
	if err != nil {
		return err
	}

	fmt.Println("md5 ", k.FingerPrintMD5())
	fmt.Println("sha1", k.FingerPrintSHA1())
	fmt.Println("pwd ", k.Decrypt(decData))

	return nil
}

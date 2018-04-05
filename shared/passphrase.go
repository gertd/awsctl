package shared

import (
	"fmt"
	"log"

	"golang.org/x/crypto/ssh/terminal"
)

// PassPhraseFunc -- callback function to set the passphrase
type PassPhraseFunc func() string

// NoPassPhrase -- return empty == no passphrase
func NoPassPhrase() string {
	return ""
}

// PassPhrase -- return argument as passphrase
func PassPhrase(pp string) PassPhraseFunc {
	return func() string { return pp }
}

// PassPhrasePrompt -- prompt for passphrase input
func PassPhrasePrompt() string {
	fmt.Printf("passphrase>> ")
	buf, err := terminal.ReadPassword(0)
	fmt.Printf("\n")
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(buf)
}

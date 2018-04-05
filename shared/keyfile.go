package shared

import (
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

// KeyFile -- key file structure
type KeyFile struct {
	FileName string
	Key      *rsa.PrivateKey
}

// ParseKeyFile -- return RSA key
func ParseKeyFile(keyFile string, passPhraseFunc PassPhraseFunc) (*KeyFile, error) {

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		return nil, fmt.Errorf("key file [%s] does not exist", keyFile)
	}

	k := KeyFile{FileName: keyFile}

	buf, err := ioutil.ReadFile(keyFile)
	if err != nil {
		return nil, err
	}

	block, rest := pem.Decode(buf)
	if len(rest) > 0 {
		return nil, fmt.Errorf("PEM decode, extra data included in key")
	}

	var keyBuf []byte
	if x509.IsEncryptedPEMBlock(block) {
		passPhrase := passPhraseFunc()
		keyBuf, err = x509.DecryptPEMBlock(block, []byte(passPhrase))
		if err != nil {
			return nil, err
		}
	} else {
		keyBuf = block.Bytes
	}

	k.Key, err = x509.ParsePKCS1PrivateKey(keyBuf)
	if err != nil {
		return nil, err
	}

	if err := k.Key.Validate(); err != nil {
		return nil, err
	}

	return &k, nil
}

// FingerPrintSHA1 -- SHA1 fingerprint private key of keyfile
// AWS generate - generate SHA1 has of PK8 key
// openssl pkcs8 -in path_to_private_key -inform PEM -outform DER -topk8 -nocrypt | openssl sha1 -c
func (k *KeyFile) FingerPrintSHA1() string {

	pkcs8Buf, err := x509.MarshalPKCS8PrivateKey(k.Key)
	if err != nil {
		log.Println(err)
		return ""
	}

	sha1Sum := sha1.Sum(pkcs8Buf)
	return fmt.Sprintf("%s", stringHash(sha1Sum[:]))
}

// FingerPrintMD5 -- MD5 fingerprint of public key of keyfile
// AWS imported keys - generate MD5 hash of public key
// openssl rsa -in path_to_private_key -pubout -outform DER | openssl md5 -c
func (k *KeyFile) FingerPrintMD5() string {

	rsaKey := k.Key
	pubBuf, err := x509.MarshalPKIXPublicKey(&rsaKey.PublicKey)
	if err != nil {
		log.Println(err)
		return ""
	}

	md5Sum := md5.Sum(pubBuf)
	return fmt.Sprintf("%s", stringHash(md5Sum[:]))
}

// Decrypt -- decrypt key using key of keyfile
func (k *KeyFile) Decrypt(chipherText []byte) string {
	b, err := rsa.DecryptPKCS1v15(nil, k.Key, chipherText)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(b)
}

func stringHash(hash []byte) string {
	out := ""
	for i := 0; i < len(hash); i++ {
		if i > 0 {
			out += ":"
		}
		out += fmt.Sprintf("%02x", hash[i]) // don't forget the leading zeroes
	}
	return out
}

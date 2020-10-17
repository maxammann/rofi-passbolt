package main

import (
	"errors"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"os"
	"strings"
)

func readKeyRing(filename string) (openpgp.EntityList, error) {
	in, err := os.Open(filename)
	if err != nil {
		return nil, nil
	}

	keyring, err := openpgp.ReadArmoredKeyRing(in)
	return keyring, in.Close()
}

func decryptArmoredMessage(armoredMessage string, keyring openpgp.EntityList) ([]byte, error) {
	block, err := armor.Decode(strings.NewReader(armoredMessage))
	if err != nil {
		return nil, err
	}

	if block.Type != "PGP MESSAGE" {
		return nil, errors.New("message does not have type PGP MESSAGE")
	}

	messageDigest, err := openpgp.ReadMessage(block.Body, keyring, nil, nil)
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(messageDigest.UnverifiedBody)
	return bytes, err
}

//func handlePrompt(keys []openpgp.Key, symmetric bool) ([]byte, error) {
//	fmt.Print("Enter Password: ")
//	//bytePassword, err := terminal.ReadPassword(syscall.Stdin)
//	var bytePasswordS string
//	_, err := fmt.Scanln(&bytePasswordS)
//	bytePassword := []byte(bytePasswordS)
//	println()
//
//	if err != nil {
//		return nil, err
//	}
//
//	for _, key := range keys {
//		err = key.PrivateKey.Decrypt(bytePassword)
//		return bytePassword, err
//	}
//
//	return bytePassword, err
//}

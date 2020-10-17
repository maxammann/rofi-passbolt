package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	//baseUrl := flag.String("base-url", "https://cloud.passbolt.com/", "base url of the passbolt instance")
	//fingerprint := flag.String("fingerprint", "0123456789ABCDEF", "fingerprint of your GPG key")
	//gpgKeyPath := flag.String("gpg-key", "key.asc", "path to your GPG key")
	//flag.Parse()

	if os.Getenv("ROFI_OUTSIDE") == "" {
		println("run this script in rofi")
		return
	}

	baseUrl := os.Getenv("PASSBOLT_BASE_URL")
	fingerprint := os.Getenv("PASSBOLT_FINGERPRINT")
	gpgKeyPath := os.Getenv("PASSBOLT_GPG_KEY_PATH")
	gpgKeyPassword := os.Getenv("PASSBOLT_GPG_KEY_PASSWORD")

	//fmt.Printf("%v\n", baseUrl)
	//fmt.Printf("%v\n", fingerprint)
	//fmt.Printf("%v\n", gpgKeyPath)

	keyring, err := readKeyRing(gpgKeyPath)
	if err != nil {
		log.Fatal(err)
		return
	}

	keys := keyring.DecryptionKeys()

	for _, key := range keys {
		key.PrivateKey.Decrypt([]byte(gpgKeyPassword))
	}

	auth, err := Login(baseUrl, fingerprint, keyring)
	if err != nil {
		log.Fatal(err)
		return
	}

	if len(os.Args) == 1 {
		fmt.Printf("Show secret\n")
		resources, err := FetchResources(baseUrl, *auth)
		if err != nil {
			log.Fatal(err)
			return
		}

		for _, resource := range resources {
			fmt.Printf("%v\x00meta\x1f%v\x1finfo\x1f%v\n", resource.Name, resource.Name, resource.Id)
		}
	} else if len(os.Args) == 2 {
		fmt.Printf("Show secret\n")
		info := os.Getenv("ROFI_INFO")
		fmt.Printf("%v\n", info)

		secret, err := FetchSecret(baseUrl, info, *auth, keyring)
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("%v\n", *secret)
	}
}

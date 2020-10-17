package main

// #include <stdio.h>
// #include <stdlib.h>
//
// #include <rofi/mode.h>
// #include <rofi/helper.h>
// #include <rofi/mode-private.h>
//
import "C"
import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"log"
	"os"
	"os/exec"
	"unsafe"
)

var currentRofiList RofiList

const OutputOptionAutotype = "autotype"
const OutputOptionPass = "pass"
const OutputOptionUser = "user"

var outputOptions = []string{OutputOptionAutotype, OutputOptionPass, OutputOptionUser}

var currentSecret string
var currentResource Resource
var xdotoolText string

type RofiList interface {
	getElement(index uint) string
	length() uint
	chooseElement(index uint) C.ModeMode
}

type ResourceList struct {
	baseUrl string
	auth    *Auth
	keyring openpgp.EntityList

	resources []Resource
}

type OutputOptionList struct {
	options []string
}

func (list ResourceList) getElement(index uint) string {
	return list.resources[index].Name
}

func (list ResourceList) getResourceSecret(index uint) (*string, error) {
	secretString, err := FetchSecret(list.baseUrl, list.resources[index].Id, *list.auth, list.keyring)
	if err != nil {
		log.Println(err)
	}

	return secretString, err
}

func (list ResourceList) chooseElement(index uint) C.ModeMode {
	fmt.Fprintf(os.Stderr, "You selected %+v\n", currentRofiList.getElement(index))

	var err error

	secret, err := list.getResourceSecret(index)

	if err != nil {
		log.Fatal(err)
		return C.MODE_EXIT
	}

	currentSecret = *secret
	currentResource = list.getResource(index)

	currentRofiList = OutputOptionList{outputOptions}
	return C.RESET_DIALOG
}

func (list ResourceList) getResource(index uint) Resource {
	return list.resources[index]
}

func (list ResourceList) length() uint {
	return uint(len(list.resources))
}

func (list OutputOptionList) getElement(index uint) string {
	return list.options[index]
}

func (list OutputOptionList) length() uint {
	return uint(len(list.options))
}

func (list OutputOptionList) chooseElement(index uint) C.ModeMode {
	element := currentRofiList.getElement(index)
	fmt.Fprintf(os.Stderr, "You selected %+v\n", element)

	if element == OutputOptionAutotype {
		xdotoolText = fmt.Sprintf("%v\t%v", currentResource.Username, currentSecret)
	} else if element == OutputOptionPass {
		xdotoolText = currentSecret
	} else if element == OutputOptionUser {
		xdotoolText = currentResource.Username
	}

	return C.MODE_EXIT
}

func FindArgumentString(key string) *string {
	args := os.Args
	for i := range args[:len(os.Args)-1] {
		if args[i] == key {
			return &args[i+1]
		}
	}
	return nil
}

//export rofi_init
func rofi_init(sw *C.Mode) int {
	baseUrl := *FindArgumentString("-base-url")
	fingerprint := *FindArgumentString("-fingerprint")
	gpgKeyPath := *FindArgumentString("-gpg-key")
	gpgKeyPassword := os.Getenv("PASSBOLT_GPG_KEY_PASSWORD")

	if gpgKeyPassword == "" {
		log.Fatal("GPG key password must be set")
		return 0
	}

	keyring, err := readKeyRing(gpgKeyPath)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	keys := keyring.DecryptionKeys()

	for _, key := range keys {
		err = key.PrivateKey.Decrypt([]byte(gpgKeyPassword))

		if err != nil {
			log.Fatal(err)
			return 1
		}
	}

	auth, err := Login(baseUrl, fingerprint, keyring)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	resources, err := FetchResources(baseUrl, *auth)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	currentRofiList = ResourceList{baseUrl, auth, keyring, resources}

	return 1
}

//export rofi_destroy
func rofi_destroy(sw *C.Mode) {
	//println("rofi_destroy")

	println(xdotoolText)

	command := exec.Command("/usr/bin/xdotool", "type", "--delay", "100", "--clearmodifiers", "--file", "-")
	buffer := bytes.Buffer{}
	buffer.Write([]byte(xdotoolText))
	command.Stdin = &buffer

	err := command.Run()
	if err != nil {
		log.Println(err)
	}
}

//export rofi_get_num_entries
func rofi_get_num_entries(sw *C.Mode) uint {
	//println("rofi_get_num_entries")
	//fmt.Fprintf(os.Stderr, "%+v\n", sw)
	return uint(currentRofiList.length())
}

//export rofi_get_display_value
func rofi_get_display_value(sw *C.Mode, selected_line uint, state *int, attr_list **C.GList, get_entry int) *C.char {
	//println("rofi_get_display_value")
	//fmt.Fprintf(os.Stderr, "%+v\n", sw)
	//fmt.Fprintf(os.Stderr, "%+v\n", selected_line)
	//fmt.Fprintf(os.Stderr, "%+v\n", state)
	//fmt.Fprintf(os.Stderr, "%+v\n", attr_list)
	//fmt.Fprintf(os.Stderr, "%+v\n", get_entry)

	if get_entry == 0 {
		return nil
	}

	return C.CString(currentRofiList.getElement(selected_line)) // This return value gets freed automatically by rofi https://github.com/davatorium/rofi/blob/011908e1ffda06d09c9f163867cb6a7d68fd6c20/source/view.c#L1023
}

//export rofi_result
func rofi_result(sw *C.Mode, mretv int, input **C.char, selected_line uint) C.ModeMode {
	//println("rofi_result")
	//fmt.Fprintf(os.Stderr, "%+v\n", sw)
	//fmt.Fprintf(os.Stderr, "%+v\n", mretv)
	//fmt.Fprintf(os.Stderr, "%+v\n", input)
	//fmt.Fprintf(os.Stderr, "%+v\n", selected_line)

	if (mretv & C.MENU_OK) != 0 {
		return currentRofiList.chooseElement(selected_line)
	}

	return C.MODE_EXIT
}

//export rofi_token_match
func rofi_token_match(sw *C.Mode, tokens **C.rofi_int_matcher, index uint) int {
	println("rofi_token_match")
	//fmt.Fprintf(os.Stderr, "%+v\n", sw)
	//firstToken := **tokens
	//fmt.Fprintf(os.Stderr, "%+v\n", firstToken)
	//fmt.Fprintf(os.Stderr, "%+v\n", index)

	//var regex *C.GRegex = firstToken.regex
	//fmt.Fprintf(os.Stderr, "%+v\n", *regex)

	current := unsafe.Pointer(tokens)
	size := unsafe.Sizeof(tokens)

	//fmt.Fprintf(os.Stderr, "%+v\n", **(**C.rofi_int_matcher) (current))
	//fmt.Fprintf(os.Stderr, "%+v\n", *(**C.rofi_int_matcher) (unsafe.Pointer(uintptr(current) + size)))
	//println(size)

	var match C.gboolean = 1

	input := currentRofiList.getElement(index)

	for *(**C.rofi_int_matcher)(current) != nil {
		token := **(**C.rofi_int_matcher)(current)
		current = unsafe.Pointer(uintptr(current) + size)
		//fmt.Fprintf(os.Stderr, "%+v\n", current)
		//fmt.Fprintf(os.Stderr, "%+v\n", token)

		cInput := C.CString(input)
		match = C.g_regex_match(token.regex, cInput, 0, nil)
		C.free(unsafe.Pointer(cInput))
		match ^= token.invert
	}

	if match == 0 {
		return 0
	} else {
		return 1
	}
}

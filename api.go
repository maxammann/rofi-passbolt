package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/openpgp"
	"net/http"
	"net/url"
	"strings"
)

type Auth struct {
	cookies []*http.Cookie
}

func FetchResources(baseUrl string, auth Auth) ([]Resource, error) {
	var resourcesResult ResourcesResult

	cookies := auth.cookies
	_, err := get(baseUrl, "/resources.json?api-version=v2", cookies, &resourcesResult)
	if err != nil {
		return nil, err
	}
	return resourcesResult.Body, nil
}

func FetchSecret(baseUrl string, resourceId string, auth Auth, keyring openpgp.EntityList) (*string, error) {
	var secretResult SecretResult

	cookies := auth.cookies
	_, err := get(baseUrl, fmt.Sprintf("/secrets/resource/%v.json?api-version=v2", resourceId), cookies, &secretResult)
	if err != nil {
		return nil, err
	}

	armoredSecret := secretResult.Body.Data
	secret, err := decryptArmoredMessage(armoredSecret, keyring)

	maybe_json := make(map[string]string)

    if err := json.Unmarshal(secret, &maybe_json); err != nil {
		secretString := string(secret)
		return &secretString, nil
    }

	ret := maybe_json["password"]

	return &ret, nil
}

func Login(baseUrl string, fingerprint string, keyring openpgp.EntityList) (*Auth, error) {
	var stage1Result Stage1Result

	resp, err := postForm(baseUrl, "/auth/login.json?api-version=v2", url.Values{"data[gpg_auth][keyid]": {fingerprint}}, &stage1Result)

	if err != nil {
		return nil, err
	}

	encodedGpgUserAuthToken := resp.Header.Get("X-GPGAuth-User-Auth-Token")

	gpgUserAuthToken, err := url.QueryUnescape(encodedGpgUserAuthToken)
	if err != nil {
		return nil, err
	}

	armoredMessage := strings.ReplaceAll(gpgUserAuthToken, "\\", "")

	bytes, err := decryptArmoredMessage(armoredMessage, keyring)
	if err != nil {
		return nil, err
	}
	token := string(bytes)

	//for name, headers := range resp.Header {
	//	print(name)
	//	for _, h := range headers {
	//		println(h)
	//	}
	//}

	var stage2Result Stage2Result

	resp, err = postForm(baseUrl, "/auth/login.json?api-version=v2", url.Values{"data[gpg_auth][keyid]": {fingerprint}, "data[gpg_auth][user_token_result]": {token}}, &stage2Result)
	if err != nil {
		return nil, err
	}

	auth := Auth{resp.Cookies()}
	return &auth, nil
}

func getCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, cookie := range cookies {
		if cookie.Name == name {
			return cookie
		}
	}

	return nil
}

func postForm(baseUrl string, path string, values url.Values, jsonResult interface{}) (*http.Response, error) {
	resp, err := http.PostForm(fmt.Sprintf("%v%v", baseUrl, path), values)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, err
	}

	dec := json.NewDecoder(resp.Body)
	//dec.DisallowUnknownFields()
	err = dec.Decode(&jsonResult)
	if err != nil {
		return resp, err
	}

	return resp, err
}

func get(baseUrl string, path string, cookies []*http.Cookie, jsonResult interface{}) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%v%v", baseUrl, path), nil)

	if err != nil {
		return nil, err
	}

	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	resp, err := client.Do(req)

	if err != nil {
		return resp, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp, errors.New("status code is not 200")
	}

	dec := json.NewDecoder(resp.Body)
	//dec.DisallowUnknownFields()
	err = dec.Decode(&jsonResult)

	return resp, err
}

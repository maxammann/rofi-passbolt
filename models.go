package main

import "go/types"

type Header struct {
	Id         string
	Status     string
	ServerTime int64
	Title      string
	Action     string
	Message    string
	Url        string
	Code       int32
}

type Stage1Result struct {
	Header Header
	Body   types.Nil
}

type Stage2Result struct {
	Header Header
	Body   struct {
		Active  bool
		Created string
		Deleted bool
		Gpgkey  struct {
			Armored_key string
			Bits        int
			Created     string
			Deleted     bool
			Expires     string
			Fingerprint string
			Id          string
			Key_created string
			Key_id      string
			Type        string
			Uid         string
			User_id     string
		}
		Groups_users   []interface{}
		Last_logged_in string
		Modified       string
		Profile        map[string]interface{}
		Role           map[string]interface{}
		Role_id        string
		Username       string
	}
	//	Body   map[string]interface{}
}

type SecretResult struct {
	Header Header
	Body   struct {
		Created string
		Data    string
	}
}

type Resource struct {
	Id          string
	Name        string
	Username    string
	Description string
}

type ResourcesResult struct {
	Header Header
	Body   []Resource
}

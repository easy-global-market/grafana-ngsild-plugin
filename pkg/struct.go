package main

import "encoding/json"

type Token struct {
	Access_token       string `json:"access_token"`
	Expires_in         int    `json:"expires_in"`
	Refresh_expires_in int    `json:"refresh_expires_in"`
	Refresh_token      string `json:"refresh_token"`
	Token_type         string `json:"token_type"`
	Session_state      string `json:"session_state"`
	Scope              string `json:"scope"`
}

type Attribute struct {
	Type       string          `json:"type"`
	CreatedAt  string          `json:"createdAt"`
	ModifiedAt string          `json:"modifiedAt"`
	Object     json.RawMessage `json:"object"`
	Value      json.RawMessage `json:"value"`
}
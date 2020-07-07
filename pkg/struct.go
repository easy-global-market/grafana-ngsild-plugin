package main

//Token struct
type Token struct {
	Access_token       string `json:"access_token"`
	Expires_in         int    `json:"expires_in"`
	Refresh_expires_in int    `json:"refresh_expires_in"`
	Refresh_token      string `json:"refresh_token"`
	Token_type         string `json:"token_type"`
	Session_state      string `json:"session_state"`
	Scope              string `json:"scope"`
}

//Name struc for Apiary
type Name struct {
	Type       string `json:"type"`
	CreatedAt  string `json:"createdAt"`
	Value      string `json:"value"`
	ModifiedAt string `json:"modifiedAt"`
}

//Apiary Struct
type Apiary struct {
	Id        string `json:"id"`
	Type      string `json:"type"`
	CreatedAt string `json:"createdAt"`
	Name      Name
}

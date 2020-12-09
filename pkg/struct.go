package main

import (
	"encoding/json"
)

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
	Object     string          `json:"object"`
	Value      json.RawMessage `json:"value"`
}
type Location struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type queryModel struct {
	QueryText  string `json:"queryText"`
	Format     string `json:"format"`
	MapMetric  string `json:"attribute"`
	Context    string `json:"context"`
	EntityType string `json:"entityType"`
}

type instanceSettings struct {
	authServerUrl    string
	resource         string
	clientId         string
	clientSecret     string
	contextBrokerUrl string
}

type settingsModel struct {
	AuthServerUrl    string `json:"authServerUrl"`
	Resource         string `json:"resource"`
	ClientId         string `json:"clientId"`
	ContextBrokerUrl string `json:"contextBrokerUrl"`
}

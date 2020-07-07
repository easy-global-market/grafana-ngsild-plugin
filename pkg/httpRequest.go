package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
)

func getToken() string {
	apiUrl := "sso.eglobalmark.com"
	resource := "/auth/realms/stellio/protocol/openid-connect/token"
	data := url.Values{}
	data.Set("client_id", "jenkins-integration")
	data.Set("client_secret", "8c79cec3-db58-47fc-9c84-24296b26cee8")
	data.Set("grant_type", "client_credentials")

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload

	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		log.DefaultLogger.Warn("err", err)
		log.DefaultLogger.Info("n:", n)
	}

	//log.DefaultLogger.Info("BUFF :", buf.String())

	//Getting json from string
	in := []byte(buf.String())

	var token Token
	errr := json.Unmarshal(in, &token)
	if errr != nil {
		panic(err)
	}
	//log.DefaultLogger.Info("TOKEN: ", string(iot.Access_token))
	return string(token.Access_token)

}

func getEntityById(entity string, token string) Apiary {

	bToken := "Bearer " + token
	apiUrl := "sso.eglobalmark.com"
	resource := "/ngsi-ld/v1/entities/" + entity

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)

	r.Header.Add("Authorization", bToken)

	resp, _ := client.Do(r)

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		log.DefaultLogger.Warn("err", err)
		log.DefaultLogger.Info("n:", n)
	}
	log.DefaultLogger.Info("response Status:", "request", resp.Status)
	log.DefaultLogger.Info("BUFF :", "request", buf.String())

	//Getting json from string
	in := []byte(buf.String())

	var apiary Apiary
	errr := json.Unmarshal(in, &apiary)
	if errr != nil {
		panic(err)
	}

	//log.DefaultLogger.Info("NAME.VALUE: ", "request", string(apiary.Name.Value))
	return apiary

}

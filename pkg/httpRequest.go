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

func getToken(instSetting *instanceSettings) string {
	authServerUrl := instSetting.authServerUrl
	resource := instSetting.resource
	data := url.Values{}
	data.Set("client_id", instSetting.clientId)
	data.Set("client_secret", instSetting.clientSecret)
	data.Set("grant_type", "client_credentials")

	uri, _ := url.ParseRequestURI(authServerUrl)
	uri.Path = resource
	urlStr := uri.String()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(req)

	buf := new(strings.Builder)
	info, err := io.Copy(buf, resp.Body)
	if err != nil {
		log.DefaultLogger.Warn("err", err)
		log.DefaultLogger.Info("info:", info)
	}

	in := []byte(buf.String())

	var token Token
	errr := json.Unmarshal(in, &token)
	if errr != nil {
		log.DefaultLogger.Info("unmarshal json error :", errr)
	}
	return string(token.Access_token)

}

func getEntityById(id string, token string, instSetting *instanceSettings) map[string]json.RawMessage {

	bToken := "Bearer " + token
	contextBrokerUrl := instSetting.contextBrokerUrl
	resource := "/ngsi-ld/v1/entities/" + id + "?options=sysAttrs"

	u, _ := url.ParseRequestURI(contextBrokerUrl + resource)
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
	log.DefaultLogger.Info("buffer :", "request", buf.String())

	in := []byte(buf.String())

	var e map[string]json.RawMessage

	if err := json.Unmarshal(in, &e); err != nil {
		log.DefaultLogger.Info("unmarshal json error :", err)
	}

	return e
}

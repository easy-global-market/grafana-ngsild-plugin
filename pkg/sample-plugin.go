package main

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

// newDatasource returns datasource.ServeOpts.
func newDatasource() datasource.ServeOpts {
	// creates a instance manager for your plugin. The function passed
	// into `NewInstanceManger` is called when the instance is created
	// for the first time or when a datasource configuration changed.
	im := datasource.NewInstanceManager(newDataSourceInstance)
	ds := &SampleDatasource{
		im: im,
	}

	return datasource.ServeOpts{
		QueryDataHandler:   ds,
		CheckHealthHandler: ds,
	}
}

// SampleDatasource is an example datasource used to scaffold
// new datasource plugins with an backend.
type SampleDatasource struct {
	// The instance manager can help with lifecycle management
	// of datasource instances in plugins. It's not a requirements
	// but a best practice that we recommend that you follow.
	im instancemgmt.InstanceManager
}

// QueryData handles multiple queries and returns multiple responses.
// req contains the queries []DataQuery (where each query contains RefID as a unique identifer).
// The QueryDataResponse contains a map of RefID to the response for each query, and each response
// contains Frames ([]*Frame).
func (td *SampleDatasource) QueryData(ctx context.Context, req *backend.QueryDataRequest) (*backend.QueryDataResponse, error) {
	log.DefaultLogger.Info("QueryData ", "request", req)

	log.DefaultLogger.Info("QueryData ", "request", req.Queries)

	token := getToken()
	//token := test()
	log.DefaultLogger.Info("TOKEN : ", token)

	// create response struct
	response := backend.NewQueryDataResponse()

	mydata := []string{}
	for _, q := range req.Queries {

		var qm queryModel
		response := backend.DataResponse{}
		response.Error = json.Unmarshal(q.JSON, &qm)
		mydata = append(mydata, qm.QueryText)
	}
	for i := 0; i < len(mydata); i++ {
		log.DefaultLogger.Info("DATA ", mydata[i])
	}

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := td.query(ctx, q, token)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}
	return response, nil
}

type queryModel struct {
	Format    string `json:"format"`
	QueryText string `json:"queryText"`
}

func (td *SampleDatasource) query(ctx context.Context, query backend.DataQuery, token string) backend.DataResponse {
	// Unmarshal the json into our queryModel
	var qm queryModel
	response := backend.DataResponse{}

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	// Log a warning if `Format` is empty.
	if qm.Format == "" {
		log.DefaultLogger.Warn("format is empty. defaulting to time series")
	}

	//entityID demandé
	log.DefaultLogger.Info("Query text ", "request", qm.QueryText)

	//getEntityById(qm.QueryText, token)
	// create data frame response
	//frame := data.NewFrame("response")
	frame := data.NewFrame(qm.QueryText)

	// add the time dimension
	//frame.Fields = append(frame.Fields,
	//	data.NewField("time", nil, []time.Time{query.TimeRange.From, query.TimeRange.To}),
	//)

	frame.Fields = append(frame.Fields,
		//data.NewField("log", nil, []string{"test", "test2"}),
		data.NewField(qm.QueryText, nil, []string{"name", "createdAt", "Apiary"}),
	)

	// add values
	// frame.Fields = append(frame.Fields,
	// 	data.NewField("values", nil, []int64{5, 10}),
	// )

	// add the frames to the response
	response.Frames = append(response.Frames, frame)
	return response
}

// CheckHealth handles health checks sent from Grafana to the plugin.
// The main use case for these health checks is the test button on the
// datasource configuration page which allows users to verify that
// a datasource is working as expected.
func (td *SampleDatasource) CheckHealth(ctx context.Context, req *backend.CheckHealthRequest) (*backend.CheckHealthResult, error) {
	var status = backend.HealthStatusOk
	var message = "Data source is working !"

	if rand.Int()%2 == 0 {
		status = backend.HealthStatusError
		message = "randomized error"
	}

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

type instanceSettings struct {
	httpClient *http.Client
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	return &instanceSettings{
		httpClient: &http.Client{},
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}

func getToken() string {
	apiUrl := "https://data-hub.eglobalmark.com"
	resource := "/auth/realms/datahub/protocol/openid-connect/token"
	data := url.Values{}
	data.Set("client_id", "stelliograf")
	data.Set("client_secret", "412fff7a-a618-4313-a342-1b844d845b45")
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

	type Iot struct {
		Access_token string `json:"access_token"`
	}
	//Getting json from string
	in := []byte(buf.String())

	var iot Iot
	errr := json.Unmarshal(in, &iot)
	if errr != nil {
		panic(err)
	}
	//log.DefaultLogger.Info("TOKEN: ", string(iot.Access_token))
	return string(iot.Access_token)

}

func getEntityById(entity string, token string) {

	bToken := "Bearer " + token
	log.DefaultLogger.Info("BEARRRRRRRR", "request", bToken)
	apiUrl := "https://data-hub.eglobalmark.com"
	resource := "/ngsi-ld/v1/entities/" + entity

	u, _ := url.ParseRequestURI(apiUrl)
	u.Path = resource
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("GET", urlStr, nil)

	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Authorization", bToken)
	//r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	//r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)

	buf := new(strings.Builder)
	n, err := io.Copy(buf, resp.Body)
	if err != nil {
		log.DefaultLogger.Warn("err", err)
		log.DefaultLogger.Info("n:", n)
	}
	log.DefaultLogger.Info("response Status:", "request", resp.Status)
	log.DefaultLogger.Info("response Headers:", "request", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	log.DefaultLogger.Info("response Body:", "request", string(body))
	//log.DefaultLogger.Info("BUFF :", "request", buf.String())
	/*
		type Iot struct {
			Access_token string `json:"access_token"`
		}
		//Getting json from string
		in := []byte(buf.String())

		var iot Iot
		errr := json.Unmarshal(in, &iot)
		if errr != nil {
			panic(err)
		}

		fmt.Println("TOKEN: ", string(iot.Access_token))
		log.DefaultLogger.Info("TOKEN: ", string(iot.Access_token))
		return string(iot.Access_token)
	*/
}

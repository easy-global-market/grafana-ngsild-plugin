package main

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

//test jenkins 2
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

	token := getToken()
	log.DefaultLogger.Info("TOKEN : ", token)

	// create response struct
	response := backend.NewQueryDataResponse()

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

	entity := getEntityById(qm.QueryText, token)

	// create data frame response
	frame := data.NewFrame(qm.QueryText)

	//Store each value on a slice
	var value []string
	var createdAt []string
	var modifiedAt []string

	for k, v := range entity {
		var a Attribute
		if err := json.Unmarshal(v, &a); err == nil {

			//log.DefaultLogger.Warn("Got attribute : ", k, string(v))
			//If we have a value data, set value, else set the object data
			if string(a.Value) != "" {
				value = append(value, k+" : "+string(a.Value))
			} else {
				value = append(value, k+" : "+string(a.Object))
			}
			createdAt = append(createdAt, a.CreatedAt)
			modifiedAt = append(modifiedAt, a.ModifiedAt)

		}
	}

	frame.Fields = append(frame.Fields,
		data.NewField("Value :", nil, value),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("CreatedAt :", nil, createdAt),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("ModifiedAt :", nil, modifiedAt),
	)

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

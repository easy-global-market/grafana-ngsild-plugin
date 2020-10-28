package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

	// create response struct
	response := backend.NewQueryDataResponse()

	instance, err := td.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}
	instSetting, _ := instance.(*instanceSettings)

	log.DefaultLogger.Info("SETTINGS ", "authServerUrl", instSetting.authServerUrl, "resource", instSetting.resource, "clientId", instSetting.clientId, "clientSecret", instSetting.clientSecret, "contextBrokerUrl", instSetting.contextBrokerUrl)

	//Get token with settings param (url, resource, client_id, client_secret)
	token := getToken(instSetting)
	log.DefaultLogger.Info("TOKEN : ", "token", token)

	// loop over queries and execute them individually.
	for _, q := range req.Queries {
		res := td.query(ctx, q, instSetting, token)

		// save the response in a hashmap
		// based on with RefID as identifier
		response.Responses[q.RefID] = res
	}
	return response, nil
}

func (td *SampleDatasource) query(ctx context.Context, query backend.DataQuery, instSetting *instanceSettings, token string) backend.DataResponse {
	// Unmarshal the json into our queryModel
	var qm queryModel
	response := backend.DataResponse{}

	response.Error = json.Unmarshal(query.JSON, &qm)
	if response.Error != nil {
		return response
	}

	//entityID demandé
	log.DefaultLogger.Info("Query text ", "request", qm.QueryText)

	entity := getEntityById(qm.QueryText, token, instSetting)

	log.DefaultLogger.Info("Query Format ", "request", qm.Format)
	if qm.Format == "worldmap" {
		worldMapResponse := transformeToWorldMap(qm.QueryText, entity, response)
		return worldMapResponse
	}
	tableResponse := transformeToTable(qm.QueryText, entity, response)
	return tableResponse

}

func transformeToTable(QueryText string, entity map[string]json.RawMessage, response backend.DataResponse) backend.DataResponse {
	// create data frame response
	frame := data.NewFrame(QueryText)

	//Store each value on a slice
	var attribute []string
	var value []string
	var createdAt []string
	var modifiedAt []string

	for k, v := range entity {
		var a Attribute
		if err := json.Unmarshal(v, &a); err == nil {

			attribute = append(attribute, k)

			//If we have a value data, set value, else set the object data
			if string(a.Value) != "" {
				value = append(value, strings.Trim(string(a.Value), "\""))
			} else {
				value = append(value, strings.Trim(string(a.Object), "\""))
			}

			createdAt = append(createdAt, dateFormat(a.CreatedAt))
			modifiedAt = append(modifiedAt, dateFormat(a.ModifiedAt))
		}
	}

	frame.Fields = append(frame.Fields,
		data.NewField("Attribute", nil, attribute),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("Value ", nil, value),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("Created at", nil, createdAt),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("Modified at", nil, modifiedAt),
	)

	// add the frames to the response
	response.Frames = append(response.Frames, frame)
	return response
}

func transformeToWorldMap(QueryText string, entity map[string]json.RawMessage, response backend.DataResponse) backend.DataResponse {
	// create data frame response
	frame := data.NewFrame(QueryText)

	//Store each value on a slice
	var attribute []string
	var value []int64
	var latitude []string
	var longitude []string

	for _, v := range entity {
		var a Attribute
		if err := json.Unmarshal(v, &a); err == nil {

			if a.Type == "GeoProperty" {
				var location Location
				err := json.Unmarshal(a.Value, &location)
				if err != nil {
					log.DefaultLogger.Warn("error marshalling", "err", err)
				}
				log.DefaultLogger.Info("location ", "request", location.Coordinates)

				long := fmt.Sprintf("%f", location.Coordinates[0])
				lat := fmt.Sprintf("%f", location.Coordinates[1])

				attribute = append(attribute, QueryText)
				value = append(value, 1)
				longitude = append(longitude, long)
				latitude = append(latitude, lat)
			}
		}
	}

	frame.Fields = append(frame.Fields,
		data.NewField("attribute", nil, attribute),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("metric", nil, value),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("latitude", nil, latitude),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("longitude", nil, longitude),
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

	return &backend.CheckHealthResult{
		Status:  status,
		Message: message,
	}, nil
}

func newDataSourceInstance(setting backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
	//get settings (tokenUrl, resource, client_id)
	var settings settingsModel
	err := json.Unmarshal(setting.JSONData, &settings)
	if err != nil {
		log.DefaultLogger.Warn("error marshalling", "err", err)
		return nil, err
	}
	//get secure settings (client_secret)
	var secureData = setting.DecryptedSecureJSONData
	clientSecret := secureData["clientSecret"]

	//log.DefaultLogger.Info("SETTINGS ", "authServerUrl", settings.AuthServerUrl, "resource", settings.Resource, "clientId", settings.ClientId, "clientSecret", clientSecret, "contextBrokerUrl", settings.ContextBrokerUrl)

	return &instanceSettings{
		authServerUrl:    settings.AuthServerUrl,
		resource:         settings.Resource,
		clientId:         settings.ClientId,
		clientSecret:     clientSecret,
		contextBrokerUrl: settings.ContextBrokerUrl,
	}, nil
}

func (s *instanceSettings) Dispose() {
	// Called before creatinga a new instance to allow plugin authors
	// to cleanup.
}

func dateFormat(inputDate string) string {
	if inputDate != "" {
		//input format is like this layout
		layout := "2006-01-02T15:04:05.999999Z"
		time, _ := time.Parse(layout, inputDate)
		timeToDisplay := time.Format("2 Jan 2006 15:04:05")
		return timeToDisplay
	}
	return inputDate
}

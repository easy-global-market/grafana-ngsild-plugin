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

	//QueryText is the entityID that user set on query panel
	log.DefaultLogger.Info("Query text ", "request", qm.QueryText)
	//Context is the dashboard variable named : context
	log.DefaultLogger.Info("Context ", "request", qm.Context)
	//EntityType is the EntityType that user set on query panel
	log.DefaultLogger.Info("EntityType ", "request", qm.EntityType)
	//ValueFilterQuery to add in request
	log.DefaultLogger.Info("ValueFilterQuery ", "request", qm.ValueFilterQuery)

	var entity []map[string]json.RawMessage
	if qm.QueryText != "" {
		entity = getEntityById(qm.QueryText, qm.Context, token, instSetting)
	} else {
		entity = getEntitesByType(qm.EntityType, qm.ValueFilterQuery, qm.Context, token, instSetting)
	}

	log.DefaultLogger.Info("Query Format ", "request", qm.Format)
	if qm.Format == "worldmap" {
		worldMapResponse := transformToWorldMap(qm.QueryText, qm.MapMetric, entity, response)
		return worldMapResponse
	} else {
		tableResponse := transformToTable(qm.QueryText, entity, response)
		return tableResponse
	}

}

func transformToTable(QueryText string, entity []map[string]json.RawMessage, response backend.DataResponse) backend.DataResponse {
	// create data frame response
	frame := data.NewFrame(QueryText)

	//Store each value on a slice
	var attribute []string
	var value []string
	var createdAt []string
	var modifiedAt []string

	for _, element := range entity {
		for k, v := range element {
			var a Attribute
			if err := json.Unmarshal(v, &a); err == nil {

				attribute = append(attribute, k)

				//If attribute is a GeoProperty we set the coordinates as value
				if a.Type == "GeoProperty" {
					var location Location
					err := json.Unmarshal(a.Value, &location)
					if err == nil {
						log.DefaultLogger.Info("location ", "request", location.Coordinates)
						coord := fmt.Sprintf("%f", location.Coordinates)
						value = append(value, coord)
					} else {
						log.DefaultLogger.Warn("error marshalling", "err", err)
					}

				} else if a.Type == "Property" {
					value = append(value, strings.Trim(string(a.Value), "\""))
				} else if a.Type == "Relationship" {
					value = append(value, string(a.Object))
				}

				createdAt = append(createdAt, dateFormat(a.CreatedAt))
				modifiedAt = append(modifiedAt, dateFormat(a.ModifiedAt))
			}
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

func transformToWorldMap(QueryText string, MapMetric string, entity []map[string]json.RawMessage, response backend.DataResponse) backend.DataResponse {
	// create data frame response
	frame := data.NewFrame(QueryText)

	//Store each value on a slice
	var attribute []string
	var value []string
	var latitude []string
	var longitude []string

	for _, element := range entity {
		for k, v := range element {
			var a Attribute
			if err := json.Unmarshal(v, &a); err == nil {

				//MapMetric is the attribute that we want to display on the map
				if MapMetric != "" && MapMetric == k {
					attribute = append(attribute, MapMetric)
					value = append(value, string(a.Value))
				}

				if a.Type == "GeoProperty" {
					var location Location
					err := json.Unmarshal(a.Value, &location)
					if err != nil {
						log.DefaultLogger.Warn("error marshalling", "err", err)
					}

					long := fmt.Sprintf("%f", location.Coordinates[0])
					lat := fmt.Sprintf("%f", location.Coordinates[1])

					longitude = append(longitude, long)
					latitude = append(latitude, lat)
				}
			}
			//If no specific attribute has been asked for, we set the entityId and value to 1 to display it anyway
			if MapMetric == "" && k == "id" {
				attribute = append(attribute, strings.Trim(string(v), "\""))
				value = append(value, "1")
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

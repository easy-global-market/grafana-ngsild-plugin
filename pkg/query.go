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
	// create response struct
	response := backend.NewQueryDataResponse()

	instance, err := td.im.Get(req.PluginContext)
	if err != nil {
		return nil, err
	}
	instSetting, _ := instance.(*instanceSettings)

	//Get token with settings param (url, resource, client_id, client_secret)
	token := getToken(instSetting)

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

	var entity []byte
	if qm.EntityId != "" {
		entity = getEntityById(qm.EntityId, qm.Context, token, instSetting)
	} else {
		entity = getEntitesByType(qm.EntityType, qm.ValueFilterQuery, qm.Context, token, instSetting)
	}

	if qm.Format == "worldmap" {
		worldMapResponse := transformToWorldMap(qm, entity, response)
		return worldMapResponse
	} else {
		tableResponse := transformToTable(qm, entity, response)
		return tableResponse
	}
}

// Return a DataResponse to display data in table view
//(The dataResponse contains a frame with 5 fields : attributes, metrics, multiAttributeValues, createdAt, modifiedAt)
func transformToTable(qm queryModel, entitiesByte []byte, response backend.DataResponse) backend.DataResponse {
	var entityId = qm.EntityId
	var metadataSelector = qm.MetadataSelector
	var hasMetadataSelector = metadataSelector != ""

	// create data frame response
	frame := data.NewFrame(entityId)

	//Store each value on a slice
	var attributes []string
	var metrics []string
	var createdAt []string
	var modifiedAt []string
	var multiAttributeValues []string

	var entities []interface{}
	json.Unmarshal(entitiesByte, &entities)

	// Range over entities
	for entity := 0; entity < len(entities); entity++ {
		var foundMetadataSelector = false
		entityInterface := entities[entity].(map[string]interface{})
		// Range over attributes
		for k, v := range entityInterface {
			switch attribute := v.(type) {

			case string: // Handle case where attribute value is string (id, type, createdAt...)
			case []interface{}:
				if k != "@context" {
					//Range over attributes
					for _, multiAttribute := range attribute {
						var propertyInterface = multiAttribute.(map[string]interface{})
						var currentValue string
						var currentUnitCode string
						var currentCreatedAt string
						var currentModifiedAt string
						//Range over properties
						for propertyKey, propertyValue := range propertyInterface {
							//We get the value if it's a Property or object if it's a Relationship
							if propertyKey == "value" || propertyKey == "object" {
								currentValue = fmt.Sprintf("%v", propertyValue)
							}
							if propertyKey == "unitCode" {
								currentUnitCode = fmt.Sprintf("%v", propertyValue)
							}
							if propertyKey == "createdAt" {
								currentCreatedAt = fmt.Sprintf("%v", propertyValue)
							}
							if propertyKey == "modifiedAt" {
								currentModifiedAt = fmt.Sprintf("%v", propertyValue)
							}
							//Getting metadataSelector value and unitCode
							if propertyKey == metadataSelector {
								foundMetadataSelector = true
								var metadataSelectorPropertyInterface = propertyValue.(map[string]interface{})
								var metadataSelectorValueString = fmt.Sprintf("%v", metadataSelectorPropertyInterface["value"])
								var metadataSelectorUnitCodeString = fmt.Sprintf("%v", metadataSelectorPropertyInterface["unitCode"])

								var multiAttributeValue = buildString("", metadataSelectorValueString, metadataSelectorUnitCodeString, "", "", "")
								multiAttributeValues = append(multiAttributeValues, multiAttributeValue)
							}
						}
						//If current attribute don't have the metadataSelector
						if metadataSelector != "" && !foundMetadataSelector {
							multiAttributeValues = append(multiAttributeValues, "")
						}

						attributes = append(attributes, k)
						createdAt = append(createdAt, dateFormat(currentCreatedAt))
						modifiedAt = append(modifiedAt, dateFormat(currentModifiedAt))
						metrics = append(metrics, buildString("", currentValue, currentUnitCode, "", "", ""))
					}
				}

			case interface{}:
				var attributeValueInterface = attribute.(map[string]interface{})["value"]

				//If key is "location"
				if k == "location" {
					// convert map to json
					jsonString, _ := json.Marshal(attributeValueInterface)
					// convert json to struct
					location := Location{}
					json.Unmarshal(jsonString, &location)

					coordinates := fmt.Sprintf("%f", location.Coordinates)
					metrics = append(metrics, coordinates)
					createdAt = append(createdAt, "")
					modifiedAt = append(modifiedAt, "")

					if hasMetadataSelector {
						multiAttributeValues = append(multiAttributeValues, "")
					}

				} else {
					var currentValue string
					var currentUnitCode string
					var currentCreatedAt string
					var currentModifiedAt string

					var propertyInterface = attribute.(map[string]interface{})
					//Getting the property value and unitCode
					for propertyKey, propertyValue := range propertyInterface {
						//We get the value if it's a Property or object if it's a Relationship
						if propertyKey == "value" || propertyKey == "object" {
							currentValue = fmt.Sprintf("%v", propertyValue)
						}
						if propertyKey == "unitCode" {
							currentUnitCode = fmt.Sprintf("%v", propertyValue)
						}
						if propertyKey == "createdAt" {
							currentCreatedAt = fmt.Sprintf("%v", propertyValue)
						}
						if propertyKey == "modifiedAt" {
							currentModifiedAt = fmt.Sprintf("%v", propertyValue)
						}
						//Getting metadataSelector value and unitCode
						if propertyKey == metadataSelector {
							foundMetadataSelector = true
							var metadataSelectorPropertyInterface = propertyValue.(map[string]interface{})

							var metadataSelectorValueString = fmt.Sprintf("%v", metadataSelectorPropertyInterface["value"])
							var metadataSelectorUnitCodeString = fmt.Sprintf("%v", metadataSelectorPropertyInterface["unitCode"])

							mutltiAttributeValue := buildString("", metadataSelectorValueString, metadataSelectorUnitCodeString, "", "", "")
							multiAttributeValues = append(multiAttributeValues, mutltiAttributeValue)
						}
					}
					//If current attribute don't have the metadataSelector
					if metadataSelector != "" && !foundMetadataSelector {
						multiAttributeValues = append(multiAttributeValues, "")
					}

					metrics = append(metrics, buildString("", currentValue, currentUnitCode, "", "", ""))
					createdAt = append(createdAt, dateFormat(currentCreatedAt))
					modifiedAt = append(modifiedAt, dateFormat(currentModifiedAt))
				}
				attributes = append(attributes, k)

			default:
				log.DefaultLogger.Error(k, "is of a type I don't know how to handle")
			}
		}
	}

	frame.Fields = append(frame.Fields,
		data.NewField("Attribute", nil, attributes),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("Value ", nil, metrics),
	)
	if hasMetadataSelector {
		frame.Fields = append(frame.Fields,
			data.NewField(metadataSelector, nil, multiAttributeValues),
		)
	}
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

// Return a DataResponse to display data in map view
//(The dataResponse contains a frame with 6 fields : entitiesId, attributes, metrics, latitudes, longitudes, multiAttributeValues)
func transformToWorldMap(qm queryModel, entitiesByte []byte, response backend.DataResponse) backend.DataResponse {
	var VALUES_SEPARATOR = ","
	var entityId = qm.EntityId
	var metadataSelector = qm.MetadataSelector
	var mapMetric = qm.MapMetric
	var hasMetadataSelector = metadataSelector != ""

	// create data frame response
	frame := data.NewFrame(entityId)
	//Store each value on a slice
	var entitiesId []string
	var attributes []string
	var metrics []string
	var latitudes []float64
	var longitudes []float64
	var multiAttributeValues []string

	var entities []interface{}
	json.Unmarshal(entitiesByte, &entities)

	// Range over entities
	for entity := 0; entity < len(entities); entity++ {
		entityInterface := entities[entity].(map[string]interface{})
		var foundAttribute = false
		var hasLocation = false
		var foundMetadataSelector = false
		entityId = fmt.Sprintf("%v", entityInterface["id"])

		// Range over attributes
		for k, v := range entityInterface {
			switch attribute := v.(type) {

			case string: // Handle case where attribute value is string (id, type, createdAt...)
			case []interface{}:
				//We find the attribute
				if k == mapMetric {
					foundAttribute = true
					var allMultiAttributeValues = ""
					var currentValue string
					var currentUnitCode string
					var currentMetadataSelectorValue string
					var currentMetadataSelectorUnitCode string
					//Range over properties
					for _, multiAttribute := range attribute {
						var propertyInterface = multiAttribute.(map[string]interface{})
						//Getting the property value and unitCode
						for propertyKey, propertyValue := range propertyInterface {
							//We get the value if it's a Property or object if it's a Relationship
							if propertyKey == "value" || propertyKey == "object" {
								currentValue = fmt.Sprintf("%v", propertyValue)
							}
							if propertyKey == "unitCode" {
								currentUnitCode = fmt.Sprintf("%v", propertyValue)
							}
							//Getting metadataSelector value and unitCode
							if propertyKey == metadataSelector {
								var metadataSelectorPropertyInterface = propertyValue.(map[string]interface{})

								currentMetadataSelectorValue = fmt.Sprintf("%v", metadataSelectorPropertyInterface["value"])
								currentMetadataSelectorUnitCode = fmt.Sprintf("%v", metadataSelectorPropertyInterface["unitCode"])
							}
						}
						allMultiAttributeValues = buildString(allMultiAttributeValues, currentValue, currentUnitCode, currentMetadataSelectorValue, currentMetadataSelectorUnitCode, VALUES_SEPARATOR)
					}
					entitiesId = append(entitiesId, entityId)
					attributes = append(attributes, mapMetric)
					firstAttributeValue := strings.Split(allMultiAttributeValues, " ")
					metrics = append(metrics, firstAttributeValue[0])
					multiAttributeValues = append(multiAttributeValues, allMultiAttributeValues)
				}

			case interface{}:
				//We find the attribute
				if k == mapMetric {
					foundAttribute = true
					var currentValue string
					var currentUnitCode string

					var propertyInterface = attribute.(map[string]interface{})
					//Getting the property value and unitCode
					for propertyKey, propertyValue := range propertyInterface {
						//We get the value if it's a Property or object if it's a Relationship
						if propertyKey == "value" || propertyKey == "object" {
							currentValue = fmt.Sprintf("%v", propertyValue)
						}
						if propertyKey == "unitCode" {
							currentUnitCode = fmt.Sprintf("%v", propertyValue)
						}
						//Getting metadataSelector value and unitCode
						if propertyKey == metadataSelector {
							foundMetadataSelector = true
							var metadataSelectorPropertyInterface = propertyValue.(map[string]interface{})

							var metadataSelectorValueString = fmt.Sprintf("%v", metadataSelectorPropertyInterface["value"])
							var metadataSelectorUnitCodeString = fmt.Sprintf("%v", metadataSelectorPropertyInterface["unitCode"])

							mutltiAttributeValue := buildString("", currentValue, currentUnitCode, metadataSelectorValueString, metadataSelectorUnitCodeString, VALUES_SEPARATOR)
							multiAttributeValues = append(multiAttributeValues, mutltiAttributeValue)
						}
					}
					//If one attribute don't have the metadataSelector
					if metadataSelector != "" && !foundMetadataSelector {
						mutltiAttributeValue := buildString("", currentValue, currentUnitCode, "", "", VALUES_SEPARATOR)
						multiAttributeValues = append(multiAttributeValues, mutltiAttributeValue)
					}

					entitiesId = append(entitiesId, entityId)
					attributes = append(attributes, mapMetric)
					metrics = append(metrics, currentValue)
				}

				//If key is "location"
				if k == "location" {
					var attributeValueInterface = attribute.(map[string]interface{})["value"]
					hasLocation = true
					// convert map to json
					jsonString, _ := json.Marshal(attributeValueInterface)
					// convert json to struct
					location := Location{}
					json.Unmarshal(jsonString, &location)

					latitudes = append(latitudes, location.Coordinates[1])
					longitudes = append(longitudes, location.Coordinates[0])
				}
			default:
				log.DefaultLogger.Info(k, "is of a type I don't know how to handle")
			}
		}
		// If we don't have location for an entity, but we already find the desired attribute
		// We can't display an entity without location, so display nothing for this entity
		if foundAttribute && !hasLocation {
			entitiesId = entitiesId[:len(entitiesId)-1]
			attributes = attributes[:len(attributes)-1]
			metrics = metrics[:len(metrics)-1]
			if hasMetadataSelector && foundMetadataSelector {
				multiAttributeValues = multiAttributeValues[:len(multiAttributeValues)-1]
			}
		}
		//If we have location for an entity but not the desired attribute
		if hasLocation && !foundAttribute {
			if mapMetric != "" {
				latitudes = latitudes[:len(latitudes)-1]
				longitudes = longitudes[:len(longitudes)-1]
			} else {
				//That means user didn't enter MapMetric, but entity has a location. So just display the location
				entitiesId = append(entitiesId, entityId)
				attributes = append(attributes, "no metric")
				metrics = append(metrics, "0")
			}
		}
	}

	frame.Fields = append(frame.Fields,
		data.NewField("id", nil, entitiesId),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("attribute", nil, attributes),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("metric", nil, metrics),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("latitude", nil, latitudes),
	)
	frame.Fields = append(frame.Fields,
		data.NewField("longitude", nil, longitudes),
	)
	if hasMetadataSelector {
		multiAttributeColumnName := mapMetric + " (" + metadataSelector + ")"
		frame.Fields = append(frame.Fields,
			data.NewField(multiAttributeColumnName, nil, multiAttributeValues),
		)
	}

	// add the frames to the response
	response.Frames = append(response.Frames, frame)
	return response
}

//Build string with or without multiAttribute and unitCode
//Ex : 10 CEL (20 MTR) / 10 (20 MTR) / 10 CEL (20) / 10 CEL
func buildString(accumulator string, value string, valueUnitCode string, metadataSelectorValue string, metadataSelectorUnitCode string, VALUES_SEPARATOR string) (buildedString string) {
	var result string
	var metadataSelectorBuildString string
	var valueUnitCodeBuildString = valueUnitCode
	var metadataSelectorUnitCodeBuildString = metadataSelectorUnitCode

	if valueUnitCodeBuildString != "" {
		valueUnitCodeBuildString = " " + valueUnitCodeBuildString
	}
	if metadataSelectorUnitCodeBuildString != "" {
		metadataSelectorUnitCodeBuildString = " " + metadataSelectorUnitCodeBuildString
	}
	if metadataSelectorValue != "" {
		metadataSelectorBuildString = " (" + metadataSelectorValue + metadataSelectorUnitCodeBuildString + ")"
	}

	if accumulator != "" {
		result = accumulator + VALUES_SEPARATOR + " " + value + valueUnitCodeBuildString + metadataSelectorBuildString
	} else {
		result = value + valueUnitCodeBuildString + metadataSelectorBuildString
	}
	return result
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
		log.DefaultLogger.Error("error marshalling", "err", err)
		return nil, err
	}
	//get secure settings (client_secret)
	var secureData = setting.DecryptedSecureJSONData
	clientSecret := secureData["clientSecret"]

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

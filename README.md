# NGSI-LD plugin for Grafana

To use NGSI-LD plugin, create a new panel and select datasource NGSI-LD
Query fields :
* **EntityIdentifier** : to get an entity by an id
* **Entity Type** : to get entities by a type
* **Value Filter Query** : to filter entities when getting them by type
* **Format** : A list with "Table" and "World Map"
* **Attribute to use as a metric** : to set an attribute as the value of the entity in the map

/!\ When you use Entity Type to get entities by type, you need to set a context.
Click on the dashboard setting > Variables > add a variable with the Name "**context**"

**Table** : is the format to display entities in a table. You will see the attributes, the values, createdAt and modifedAt
    
**World Map** : is the format to display entity on a Map. It return the entityId, a metric value, latitude and longitude
 (If you display it as table, grafana will truncate longitude and latitude. You can go to the "Overrides" field and set an "decimals" override to 6 for latitude and longitude)

 To display entity location on a map, select Woldmap Panel. For the settings :
  - "Location Data" : table
  - "Table Query Format" : coordinates
  - "Location Name Field" : attribute
  - "Metric Field" : metric (by default)
  - "Latitude Field" : latitude (by default)
  - "Longitude Field" : longitude (by default)


## Frontend development

* Install dependencies

```BASH
yarn install
```

* Build plugin in production mode

```BASH
yarn build
```

or

```BASH
npm run-script build
```

## Backend development

* Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

* Build backend plugin binaries for Linux, Windows and Darwin:

```BASH
mage -v
```


# NGSI-LD plugin for Grafana

## Frontend

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

## Backend

* Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

* Build backend plugin binaries for Linux, Windows and Darwin:

```BASH
mage -v
```

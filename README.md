# Ngsild plugin for Grafana

[![CircleCI](https://circleci.com/gh/grafana/simple-datasource-backend/tree/master.svg?style=svg)](https://circleci.com/gh/grafana/simple-datasource-backend/tree/master)


### Frontend

1. Install dependencies
```BASH
yarn install
```

2. Build plugin in production mode
```BASH
yarn build
```
or
```BASH
npm run-script build
```

### Backend

1. Update [Grafana plugin SDK for Go](https://grafana.com/docs/grafana/latest/developers/plugins/backend/grafana-plugin-sdk-for-go/) dependency to the latest minor version:

```bash
go get -u github.com/grafana/grafana-plugin-sdk-go
```

2. Build backend plugin binaries for Linux, Windows and Darwin:
```BASH
mage -v
```
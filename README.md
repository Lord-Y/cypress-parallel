# cypress-parallel-api [![CircleCI](https://circleci.com/gh/Lord-Y/cypress-parallel-api/tree/main.svg?style=svg)](https://circleci.com/gh/Lord-Y/cypress-parallel-api?branch=main)

`cypress-parallel-api` is the api used by `cypress-parallel` frontend UI.
The api permit to create pods inside your kubernetes cluster triggered by curl commands or via UI.

## API url

By default, the api url is `http://127.0.0.1:8080` but it can be override with os environment variable `CYPRESS_PARALLEL_API_URL`.

## Database

Our api is developped with PostgresSQL database so the environment variable CYPRESS_PARALLEL_API_DB_URI must be set:
```bash
export CYPRESS_PARALLEL_API_DB_URI="postgres://USERNAME:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable"
```

## Development
### Kind

During you local development, you must set the variable `CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE` in order to create to make the api loggued in with your `.kube/config`

```bash
export CYPRESS_PARALLEL_API_K8S_CLIENT_OUTSIDE=true
```

Please read `Kind` setup [here](./_developments/README.md)

### Start your postgres sql instance

```bash
sudo docker-compose up -d -f docker-compose.yml.yaml
```

### Debugging

To enable the debug mode on the api:
```bash
export CYPRESS_PARALLEL_API_LOG_LEVEL=debug
```

To enable the debug mode on the cli:
```bash
export CYPRESS_PARALLEL_CLI_LOG_LEVEL=debug
```

## Content types

Supported content types are:
- application/x-www-form-urlencoded
- application/json

## Unit testing

Make sure to add unit testing for almost every new features in order to ensure the quality of the api.

Run tests with:
```bash
go test -v ./... -coverprofile=coverage.out
```

See covering in the browser with:
```bash
go tool cover -html=coverage.out
```

See covering in the shell with:
```bash
go tool cover -func=coverage.out
```

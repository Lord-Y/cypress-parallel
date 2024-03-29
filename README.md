# cypress-parallel [![CircleCI](https://circleci.com/gh/Lord-Y/cypress-parallel/tree/main.svg?style=svg)](https://circleci.com/gh/Lord-Y/cypress-parallel?branch=main)

`cypress-parallel` is a single binary that hold both frontend UI and API that permit to run `cypress unit testing`.

It create pods inside your kubernetes cluster triggered by api or via UI.

On this UI you will be able to:
- create, update and delete teams
- create, update and delete projects and also trigger new unit testing
- create, update and delete annotations that will be used by pods
- create, update and delete environments variables that will be used by pods
- See execution results of units testing

 `cypress-parallel` use also 2 sub-components:
 - [cypress-parallel-cli](https://github.com/Lord-Y/cypress-parallel-cli)
 - [cypress-parallel-docker-images](https://github.com/Lord-Y/cypress-parallel-docker-images)

## API

By default, the api url is `http://127.0.0.1:8080` but it can be override with os environment variable `CYPRESS_PARALLEL_URL`.

## Database

Our api is developped with PostgresSQL 10+ database so the environment variable `CYPRESS_PARALLEL_DB_URI` must be set:
```bash
export CYPRESS_PARALLEL_DB_URI="postgres://USERNAME:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable"
```

## Development

### Running UI and API seperatly during development

In a shell, start the `UI`:
```bash
export CYPRESS_PARALLEL_PORT=8080
cd ui
npm run dev
```

In an another shell, start the `API`:
```bash
export CYPRESS_PARALLEL_K8S_CLIENT_OUTSIDE=true
export CYPRESS_PARALLEL_PORT=8080
go run main.go
```

### Git hooks

Add githook like so:

```bash
git config core.hooksPath .githooks
```

### Kind

During you local development, you must set the variable `CYPRESS_PARALLEL_K8S_CLIENT_OUTSIDE` in order to make the api loggued in with your `.kube/config`

```bash
export CYPRESS_PARALLEL_K8S_CLIENT_OUTSIDE=true
```

Please read `Kind` setup [here](./_developments/README.md)

### Start your postgres sql instance

```bash
sudo docker-compose up -d
```

### Debugging

To enable the debug mode on the api:
```bash
export CYPRESS_PARALLEL_LOG_LEVEL=debug
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

## Linter
```bash
# https://golangci-lint.run/usage/install/
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

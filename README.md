# cypress-parallel-api

`cypress-parallel-api` is the api used by `cypress-parallel` frontend UI.

## API url

By default, the api url is `http://127.0.0.1` but it can be override with os environment variable `CYPRESS_PARALLEL_API_URL`.

## Database

Our api is developped with PostgresSQL database so the environment variable CYPRESS_PARALLEL_API_DB_URI must be set:
```bash
export CYPRESS_PARALLEL_API_DB_URI="postgres://USERNAME:PASSWORD@HOST:PORT/DB_NAME?sslmode=disable"
```

## Content types

Supported content types are:
- application/x-www-form-urlencoded
- application/json

## Unit testing

Unit testing must be close to 100% to avoid any shitty bug as the application will continue to grow. By reaching this percentage, we make sure that when a bug is discovered, that's not a dummy one.

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

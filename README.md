# lets-go-chat

[![Latest Stable Version](https://img.shields.io/github/v/release/brokeyourbike/lets-go-chat)](https://github.com/brokeyourbike/lets-go-chat/releases)
[![Go Reference](https://pkg.go.dev/badge/github.com/brokeyourbike/lets-go-chat.svg)](https://pkg.go.dev/github.com/brokeyourbike/lets-go-chat)
[![Go Report Card](https://goreportcard.com/badge/github.com/brokeyourbike/lets-go-chat)](https://goreportcard.com/report/github.com/brokeyourbike/lets-go-chat)
[![Maintainability](https://api.codeclimate.com/v1/badges/b477b1c392da70fdad27/maintainability)](https://codeclimate.com/github/brokeyourbike/lets-go-chat/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/b477b1c392da70fdad27/test_coverage)](https://codeclimate.com/github/brokeyourbike/lets-go-chat/test_coverage)

Let's Go Chat

## How to use

```bash
HOST=127.0.0.1 PORT=8080 go build && ./lets-go-chat
```

or with `reflex`

```bash
HOST=127.0.0.1 PORT=8080 reflex -r '\.go' -s -- sh -c "go build && ./lets-go-chat"
```

## Database

```bash
docker run -it --rm --name go-postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=secret postgres:14.0
```

## Generate code from openapi

```bash
oapi-codegen -generate types -o api/server/types.gen.go -package server api/openapi.yaml
```

```bash
oapi-codegen -generate chi-server -o api/server/server.gen.go -package server api/openapi.yaml
```

## How to test

```bash
mockery --all && MallocNanoZone=0 go test -race -shuffle=on ./...
```

## How to run load test

```bash
artillery run ./loadtest.yml --output result.json  
```

## Authors
- [Ivan Stasiuk](https://github.com/brokeyourbike) | [Twitter](https://twitter.com/brokeyourbike) | [LinkedIn](https://www.linkedin.com/in/brokeyourbike) | [stasi.uk](https://stasi.uk)

## License
[Mozilla Public License v2.0](https://github.com/brokeyourbike/lets-go-chat/blob/main/LICENSE)

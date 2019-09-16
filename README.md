# feedback

[![pipeline status](https://gitlab.com/mysteriumnetwork/feedback/badges/master/pipeline.svg)](https://gitlab.com/mysteriumnetwork/feedback/pipelines)
[![Docker Pulls](https://img.shields.io/docker/pulls/mysteriumnetwork/feedback)](https://hub.docker.com/r/mysteriumnetwork/feedback)

Service to collect user feedback

# Run

```
mage build
mage run
curl http://localhost:8080/api/v1/ping
```

# API

See [swagger.json](docs/swagger.json) for details.
Or `curl http://localhost:8080/api/v1/swagger.json` on a running server

Go API is available in `github.com/mysteriumnetwork/feedback/client` package.

# Updating API

Regenerate `swagger.json` (command below) and commit changes.

```
mage regen
```

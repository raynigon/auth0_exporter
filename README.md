# Auth0 exporter

[![GitHub Release](https://img.shields.io/github/release/raynigon/auth0_exporter.svg?style=flat)](https://github.com/raynigon/auth0_exporter/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/raynigon/auth0_exporter/v2)](https://goreportcard.com/report/github.com/raynigon/github_billing_exporter/v2)

This exporter exposes [Prometheus](https://prometheus.io/) metrics from [Auth0 Management API](https://auth0.com/docs/api/management/v2/).


## Installation

For pre-built binaries please take a look at the releases.
https://github.com/raynigon/auth0_exporter/releases

### Docker

```bash
docker pull ghcr.io/raynigon/auth0-exporter:latest
docker run --rm -p 9776:9776 ghcr.io/raynigon/auth0-exporter:latest
```

Example `docker-compose.yml`:

```yaml
auth0_exporter:
    image: ghcr.io/raynigon/auth0-exporter:latest
    command:
        - '--auth0.domain=<DOMAIN>.eu.auth0.com'
        - '--auth0.client-id=test123'
        - '--auth0.client-secret=secret-code'
    restart: always
    ports:
    - "127.0.0.1:9776:9776"
```

### Kubernetes

You can find an deployment definition at: https://github.com/raynigon/auth0_exporter/tree/main/examples/kubernetes/deployment.yaml .

## Building and running

## Token privileges

**TODO**

### Build

    make build

### Running

Running using an environment variable:

    export A0E_AUTH0_DOMAIN="<DOMAIN>.eu.auth0.com"
    export A0E_AUTH0_CLIENT_ID="test123"
    export A0E_AUTH0_CLIENT_SECRET="secret-code"
    ./auth0_exporter

Running using args:

    ./auth0_exporter \
        --auth0.domain=<DOMAIN>.eu.auth0.com 
        --auth0.client-id=test123 \
        --auth0.client-secret=secret-code

## Metrics

Name                            | Labels    | Type  | Description 
--------------------------------|-----------|-------|-------------
`auth0_stats_active_users`      | -         | Gauge | -
`auth0_stats_leaked_passwords`  | interval  | Gauge | -
`auth0_stats_logins`            | interval  | Gauge | -
`auth0_stats_signups`           | interval  | Gauge | -
`auth0_users_blocked`           | -         | Gauge | -
`auth0_users_email_verified`    | -         | Gauge | -
`auth0_users_total`             | -         | Gauge | -

### Labels
All metrics contain the labels `domain` and `tenant`.
The `domain` label contains the tenant domain given in the config.
The `tenant` label contains the friendly name of the tenant stored in auth0.

#### Interval
The interval label has three values:
1. `1` which represents a time interval of one day
1. `7` which represents a time interval of seven days
1. `30` which represents a time interval of 30 days

## Environment variables / args reference

Version    | Env		               | Arg		             | Description			                	       | Default
-----------|---------------------------|-------------------------|-------------------------------------------------|---------
\>=`0.0.1` | `A0E_AUTH0_DOMAIN`        | `--auth0.domain`	     | The domain under which the management api is reachable                                               | **REQUIRED**
\>=`0.0.1` | `A0E_AUTH0_CLIENT_ID`     | `--auth0.client-id`	 | The client id for the application which is used to monitor the tenant                                               | **REQUIRED**
\>=`0.0.1` | `A0E_AUTH0_CLIENT_SECRET` | `--auth0.client-secret` | The client secret for the application which is used to monitor the tenant                                               | **REQUIRED**
\>=`0.0.1` | `A0E_LISTEN_ADDRESS`      | `--web.listen-address`  | Address on which to expose metrics.             | `:9776`
\>=`0.0.1` | `A0E_METRICS_PATH`	       | `--web.telemetry-path`  | Path under which to expose metrics.             | `/metrics`
\>=`0.0.1` | `A0E_LOG_LEVEL`           | `--log.level`	         | -                                               | `"info"`
\>=`0.0.1` | `A0E_LOG_FORMAT`          | `--log.format`	         | -                                               | `"logfmt"`
\>=`0.0.1` | `A0E_LOG_OUTPUT`          | `--log.output`	         | -                                               | `"stdout"`



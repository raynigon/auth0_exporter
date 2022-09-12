# Auth0 exporter

[![GitHub Release](https://img.shields.io/github/release/raynigon/auth0_exporter.svg?style=flat)](https://github.com/raynigon/auth0_exporter/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/raynigon/auth0_exporter/v2)](https://goreportcard.com/report/github.com/raynigon/github_billing_exporter/v2)

This exporter exposes [Prometheus](https://prometheus.io/) metrics from [Auth0 Management API](https://auth0.com/docs/api/management/v2/).

## Metrics
```
# HELP auth0_exporter_build_info A metric with a constant '1' value labeled by version, revision, branch, and goversion from which auth0_exporter was built.
# TYPE auth0_exporter_build_info gauge
auth0_exporter_build_info{branch="",goversion="go1.17.13",revision="",version=""} 1
# HELP auth0_stats_active_users The number of active users that logged in during the last 30 days.
# TYPE auth0_stats_active_users gauge
auth0_stats_active_users{domain="<PLACEHOLDER>.auth0.com",tenant="<PLACEHOLDER>"} 123
# HELP auth0_stats_leaked_passwords Number of breached-password detections in the interval. The interval is given in days.
# TYPE auth0_stats_leaked_passwords gauge
auth0_stats_leaked_passwords{domain="<PLACEHOLDER>.auth0.com",interval="1",tenant="<PLACEHOLDER>"} 0
auth0_stats_leaked_passwords{domain="<PLACEHOLDER>.auth0.com",interval="30",tenant="<PLACEHOLDER>"} 0
auth0_stats_leaked_passwords{domain="<PLACEHOLDER>.auth0.com",interval="7",tenant="<PLACEHOLDER>"} 0
# HELP auth0_stats_logins Number of logins in the interval. The interval is given in days.
# TYPE auth0_stats_logins gauge
auth0_stats_logins{domain="<PLACEHOLDER>.auth0.com",interval="1",tenant="<PLACEHOLDER>"} 1234
auth0_stats_logins{domain="<PLACEHOLDER>.auth0.com",interval="30",tenant="<PLACEHOLDER>"} 123
auth0_stats_logins{domain="<PLACEHOLDER>.auth0.com",interval="7",tenant="<PLACEHOLDER>"} 12
# HELP auth0_stats_signups Number of signups in the interval. The interval is given in days.
# TYPE auth0_stats_signups gauge
auth0_stats_signups{domain="<PLACEHOLDER>.auth0.com",interval="1",tenant="<PLACEHOLDER>"} 12345
auth0_stats_signups{domain="<PLACEHOLDER>.auth0.com",interval="30",tenant="<PLACEHOLDER>"} 1234
auth0_stats_signups{domain="<PLACEHOLDER>.auth0.com",interval="7",tenant="<PLACEHOLDER>"} 123
# HELP auth0_users_blocked The number of blocked users.
# TYPE auth0_users_blocked gauge
auth0_users_blocked{domain="<PLACEHOLDER>.auth0.com",tenant="<PLACEHOLDER>"} 123
# HELP auth0_users_email_verified The number of users which have a verified email address.
# TYPE auth0_users_email_verified gauge
auth0_users_email_verified{domain="<PLACEHOLDER>.auth0.com",tenant="<PLACEHOLDER>"} 12345
# HELP auth0_users_total The number of existing users.
# TYPE auth0_users_total gauge
auth0_users_total{domain="<PLACEHOLDER>.auth0.com",tenant="<PLACEHOLDER>"} 12345
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 3160
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```


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

### Auth0 Application Privileges

The following privileges are needed for the Auth0 application:
**TODO**

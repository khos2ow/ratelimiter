# ratelimiter

[![Build Status](https://github.com/khos2ow/ratelimiter/workflows/build/badge.svg)](https://github.com/khos2ow/ratelimiter/actions) [![GoDoc](https://godoc.org/github.com/khos2ow/ratelimiter?status.svg)](https://godoc.org/github.com/khos2ow/ratelimiter) [![Go Report Card](https://goreportcard.com/badge/github.com/khos2ow/ratelimiter)](https://goreportcard.com/report/github.com/khos2ow/ratelimiter) [![codecov](https://codecov.io/gh/khos2ow/ratelimiter/branch/master/graph/badge.svg)](https://codecov.io/gh/khos2ow/ratelimiter) [![License](https://img.shields.io/github/license/khos2ow/ratelimiter)](https://github.com/khos2ow/ratelimiter/blob/master/LICENSE) [![Latest release](https://img.shields.io/github/v/release/khos2ow/ratelimiter)](https://github.com/khos2ow/ratelimiter/releases)

`ratelimiter` is an app to do dirstibuted rate limiting in front of backend services. It consists of:

- A command line interface ([`ratelimiter`](./cmd/main.go)) built on these packages.
- Docker [image](Dockerfile) to run ratelimiter in a containerized workload.
- Go package implementing [rate limiting](./limiter.go), which can directly be used in other projects.

## Table of Contents

- [Installation](#installation)
- [Running](#running)
  - [CLI binary](#cli-binary)
  - [Docker Container](#docker-container)
  - [Kubernetes](#kubernetes)
- [Go Package](#go-package)

## Installation

The latest version can be installed using `go get`:

``` bash
GO111MODULE="on" go get github.com/khos2ow/ratelimiter@v0.1.1
```

**NOTE:** please use the latest go to do this, ideally go 1.13.5 or greater. We use go 1.14.

This will put `ratelimiter` in `$(go env GOPATH)/bin`. If you encounter the error `ratelimiter: command not found` after installation then you may need to either add that directory to your `$PATH` as shown [here](https://golang.org/doc/code.html#GOPATH) or do a manual installation by cloning the repo and run `make build` from the repository which will put `ratelimiter` in:

```bash
$(go env GOPATH)/src/github.com/khos2ow/ratelimiter/bin/$(uname | tr '[:upper:]' '[:lower:]')-amd64/ratelimiter
```

Stable binaries are also available on the [releases](https://github.com/khos2ow/ratelimiter/releases) page. To install, download the binary for your platform from "Assets" and place this into your `$PATH`:

```bash
curl -Lo ./ratelimiter https://github.com/khos2ow/ratelimiter/releases/download/v0.1.1/ratelimiter-v0.1.1-$(uname | tr '[:upper:]' '[:lower:]')-amd64
chmod +x ./ratelimiter
mv ./ratelimiter /some-dir-in-your-PATH/ratelimiter
```

**NOTE:** Windows releases are in `EXE` format.

## Running

There are multiple ways of running `ratelimiter` service.

### CLI binary

You can run ratelimiter binary with provided flags as standalone binary:

```bash
ratelimiter \
    --rate-limit <number> \
    --rate-interval <number> \
    --rate-timeunit <time-unit> \
    --use-redis <use-redis-or-in-memory-cache> \
    --redis-url <ip-of-redis> \
    --redis-port <port-of-redis> \
    --redis-password <password-for-redis> \
    --backend-server <fdqn-or-ip-of-backend-service>
```

Note that you can provide multiple `--backend-server <string>` or one comma-separated list of servers. e.g:

```bash
ratelimiter --backend-server 1.2.3.4 --backend-server 5.6.7.8
or
ratelimiter --backend-server 1.2.3.4,5.6.7.8
```

You can also use environment variables defined on the host instead of using the flags.

| Name             | Flag               |
|------------------|--------------------|
| `RATE_LIMIT`     | `--rate-limit`     |
| `RATE_INTERVAL`  | `--rate-interval`  |
| `RATE_TIMEUNIT`  | `--rate-timeunit`  |
| `USE_REDIS`      | `--use-redis`      |
| `REDIS_URL`      | `--redis-url`      |
| `REDIS_PORT`     | `--redis-port`     |
| `REDIS_PASSWORD` | `--redis-password` |
| `BACKEND_SERVER` | `--backend-server` |

Note: You have to only use comma-separated value in `BACKEND_SERVER` environment variable.

### Docker Container

Docker images are created on each release with the following tag format:

```text
khos2ow/ratelimiter:<git-tag-without-leading-v>
e.g.
khos2ow/ratelimiter:0.1.1
```

and you can simply use the image:

```bash
docker run -d \
    -e RATE_LIMIT=<number> \
    -e RATE_INTERVAL=<number> \
    -e RATE_TIMEUNIT=<time-unit> \
    -e USE_REDIS=<use-redis-or-in-memory-cache> \
    -e REDIS_URL=<ip-of-redis> \
    -e REDIS_PORT=<port-of-redis> \
    -e REDIS_PASSWORD=<password-for-redis> \
    -e BACKEND_SERVER=<comma-separated-list-of-backend-service>
    -p 8000:8000 \
    khos2ow/ratelimiter:0.1.1
```

### Kubernetes

Prerequisites

1. Install [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)
2. Deploy [Redis cluster](https://engineering.bitnami.com/articles/deploy-and-scale-a-redis-cluster-on-kubernetes-with-bitnami-and-helm.html)
3. Update Redis URL and port and optionally backend_server in [deploy/config.yaml](./deploy/config.yaml):

    ```yaml
    ---
    apiVersion: v1
    kind: ConfigMap
    metadata:
      name: rate-limiter-config
      namespace: rate-limiter
      labels:
        app: rate-limiter
    data:
      RATE_LIMIT: "<RATE.LIMIT>"
      RATE_INTERVAL: "<RATE.INTERVAL>"
      RATE_TIMEUNIT: "<RATE.TIMEUNIT>"
      USE_REDIS: "true | false"
      REDIS_URL: "<REDIS.URL>"
      REDIS_PORT: "<REDIS.PORT>"
      BACKEND_SERVER: ""
      # BACKEND_SERVER: "<COMMA.SEPARATED.BACKEND.SERVERS>"
    ```

4. Update Redis password in [deploy/secret.yaml](./deploy/secret.yaml):

    ```yaml
    ---
    apiVersion: v1
    kind: Secret
    metadata:
      name: rate-limiter-secret
      namespace: rate-limiter
    type: Opaque
    data:
      REDIS_PASSWORD: "<REDIS.PASSWORD>"
   ```

5. Enable or update [deploy/ingress.yaml](./deploy/ingress.yaml):

    ```yaml
    ---
    apiVersion: networking.k8s.io/v1beta1
    kind: Ingress
    metadata:
      name: rate-limiter
      namespace: rate-limiter
      labels:
        app: rate-limiter
      annotations:
        kubernetes.io/ingress.class: nginx-internal
    spec:
      rules:
      - host: rate-limiter.example.com
        http:
        paths:
        - path: /
          backend:
            serviceName: rate-limiter-service
            servicePort: http
    ```

then you can deploy using `kubectl`:

```bash
kubectl apply -f deploy -n rate-limiter
```

## Go Package

ratelimiter exposes most of its functionality through Go Package which can be imported in other projects. To do that you can use package manager of your choice:

```bash
GO111MODULE="on" go get github.com/khos2ow/ratelimiter@v0.1.1
```

and then

```go
import "github.com/khos2ow/ratelimiter/pkg/ratelimiter"
```

[example/main.go](./example/main.go):

```go
package main

import (
    "fmt"
    "time"

    "github.com/khos2ow/ratelimiter/internal/data"
    "github.com/khos2ow/ratelimiter/pkg/ratelimiter"
)

func main() {
    resource := "foo"
    store := data.NewInMemory(&data.Options{})
    rule := ratelimiter.NewRule(10, 1, time.Second)
    limiter := ratelimiter.NewLimiter(rule, store)

    start := time.Now()
    fmt.Printf("limiting resource '%s' to %s\n\n", resource, rule.String())

    for i := 0; i < 25; i++ {
        allowed, err := limiter.IsAllowed(resource)
        if err != nil {
            fmt.Printf("hit #%-10derror: %-10velapsed: %f seconds\n", i+1, err.Error(), time.Now().Sub(start).Seconds())
        } else {
            fmt.Printf("hit #%-10dallowed: %-10velapsed: %f seconds\n", i+1, allowed, time.Now().Sub(start).Seconds())
        }
        time.Sleep(80 * time.Millisecond)
    }

    fmt.Printf("\ntook %f seconds\n", time.Now().Sub(start).Seconds())
}

// limiting resource 'foo' to 10 hits per second
//
// hit #1         allowed: true      elapsed: 0.000012 seconds
// hit #2         allowed: true      elapsed: 0.080239 seconds
// hit #3         allowed: true      elapsed: 0.160510 seconds
// hit #4         allowed: true      elapsed: 0.240883 seconds
// hit #5         allowed: true      elapsed: 0.321136 seconds
// hit #6         allowed: true      elapsed: 0.401298 seconds
// hit #7         allowed: true      elapsed: 0.481417 seconds
// hit #8         allowed: true      elapsed: 0.561576 seconds
// hit #9         allowed: true      elapsed: 0.641844 seconds
// hit #10        allowed: true      elapsed: 0.722082 seconds
// hit #11        allowed: false     elapsed: 0.802300 seconds
// hit #12        allowed: false     elapsed: 0.882519 seconds
// hit #13        allowed: false     elapsed: 0.962731 seconds
// hit #14        allowed: true      elapsed: 1.042958 seconds
// hit #15        allowed: true      elapsed: 1.123172 seconds
// hit #16        allowed: true      elapsed: 1.203390 seconds
// hit #17        allowed: true      elapsed: 1.283565 seconds
// hit #18        allowed: true      elapsed: 1.363771 seconds
// hit #19        allowed: true      elapsed: 1.443981 seconds
// hit #20        allowed: true      elapsed: 1.524204 seconds
// hit #21        allowed: true      elapsed: 1.604430 seconds
// hit #22        allowed: true      elapsed: 1.684739 seconds
// hit #23        allowed: true      elapsed: 1.764983 seconds
// hit #24        allowed: true      elapsed: 1.845355 seconds
// hit #25        allowed: false     elapsed: 1.925563 seconds
//
// took 2.009465 seconds
```

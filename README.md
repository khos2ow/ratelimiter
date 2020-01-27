# ratelimiter

`ratelimiter` is an app to do dirstibuted rate limiting in front of backend services. It consists of:

- A command line interface ([`ratelimiter`](./main.go)) built on these packages.
- Docker [image](./images) to run ratelimiter in a containerized workload.

## Table of Contents

- [Installation](#installation)
- [Running](#running)
  - [CLI binary](#cli-binary)
  - [Docker Container](#docker-container)
  - [Kubernetes](#kubernetes)

## Installation

The latest version can be installed using `go get`:

``` bash
GO111MODULE="off" go get github.com/khos2ow/ratelimiter@v0.0.1
```

**NOTE:** please use the latest go to do this, ideally go 1.13.5 or greater.

This will put `ratelimiter` in `$(go env GOPATH)/bin`. If you encounter the error `ratelimiter: command not found` after installation then you may need to either add that directory to your `$PATH` as shown [here](https://golang.org/doc/code.html#GOPATH) or do a manual installation by cloning the repo and run `make build` from the repository which will put `ratelimiter` in:

```bash
$(go env GOPATH)/src/github.com/khos2ow/ratelimiter/bin/$(uname | tr '[:upper:]' '[:lower:]')-amd64/ratelimiter
```

Stable binaries are also available on the [releases](https://github.com/khos2ow/ratelimiter/releases) page. To install, download the binary for your platform from "Assets" and place this into your `$PATH`:

```bash
curl -Lo ./ratelimiter https://github.com/khos2ow/ratelimiter/releases/download/v0.0.1/ratelimiter-v0.0.1-$(uname | tr '[:upper:]' '[:lower:]')-amd64
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
khos2ow/ratelimiter:0.0.1
```

and you can simply use the image:

```bash
docker run -d \
    -e REDIS_URL=<ip-of-redis> \
    -e REDIS_PORT=<port-of-redis> \
    -e REDIS_PASSWORD=<password-for-redis> \
    -e BACKEND_SERVER=<comma-separated-list-of-backend-service>
    -p 8000:8000 \
    khos2ow/ratelimiter:0.0.1
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

# Contents

1. [Overview](#1-overview)<br/>
   1.1. [Purpose](#11-purpose)<br/>
   1.2. [Definitions](#12-definitions)<br/>
2. [Configuration](#2-configuration)<br/>
3. [Deployment](#3-deployment)<br/>
   3.1. [Prerequisites](#31-prerequisites)<br/>
   3.2. [Bare](#32-bare)<br/>
   3.3. [Docker](#33-docker)<br/>
   3.4. [K8s](#34-k8s)<br/>
   &nbsp;&nbsp;&nbsp;3.4.1. [Helm](#341-helm) <br/>
4. [Usage](#4-usage)<br/>
5. [Design](#5-design)<br/>
   5.1. [Requirements](#51-requirements)<br/>
   5.2. [Approach](#52-approach)<br/>
   5.3. [Limitations](#53-limitations)<br/>
6. [Contributing](#6-contributing)<br/>
   6.1. [Versioning](#61-versioning)<br/>
   6.2. [Issue Reporting](#62-issue-reporting)<br/>
   6.3. [Building](#63-building)<br/>
   6.4. [Testing](#64-testing)<br/>
   &nbsp;&nbsp;&nbsp;6.4.1. [Functional](#641-functional)<br/>
   &nbsp;&nbsp;&nbsp;6.4.2. [Performance](#642-performance)<br/>
   6.5. [Releasing](#65-releasing)<br/>

# 1. Overview

TBD

## 1.1. Purpose

TBD

## 1.2. Definitions

TBD

# 2. Configuration

The service is configurable using the environment variables:

| Variable               | Example value  | Description                                                          |
|------------------------|----------------|----------------------------------------------------------------------|
| API_PORT               | `8080`         | gRPC API port                                                        |
| API_MATCHES_URI        | `matches:8080` | [Matches](https://github.com/awakari/matches) dependency service URI |
| API_MATCHES_BATCH_SIZE | `100`          | Matches query results size limit                                     |
| API_OUTPUT_URI         | `output:8080`  | Output dependency service URI                                        |

# 3. Deployment

## 3.1. Prerequisites

[Matches](https://github.com/awakari/matches) dependency service should be deployed.

## 3.2. Bare

Preconditions:
1. Build patterns executive using ```make build```

Then run the command:
```shell
API_PORT=8082 \
API_MATCHES_URI=localhost:8080 \
API_OUTPUT_URI=localhost:8081 \
./router
```

## 3.3. Docker

```shell
make run
```

## 3.4. K8s

### 3.4.1. Helm

Create a helm package from the sources:
```shell
helm package helm/router/
```

Install the helm chart:
```shell
helm install router ./router-<CHART_VERSION>.tgz
```

where
* `<CHART_VERSION>` is the helm chart version

# 4. Usage

The service provides a gRPC API for routing a message.

Example command:
```shell
grpcurl \
  -plaintext \
  -proto api/grpc/service.proto \
  -d @ \
  localhost:8080 \
  router.Service/Route
```
Payload:
```json
{
   "id": "3426d090-1b8a-4a09-ac9c-41f2de24d5ac",
   "metadata": {
      "source": {
         "ce_string": "com.lipsum"
      },
      "specversion": {
         "ce_string": "1.0"
      },
      "type": {
         "ce_string": "com.lipsum"
      }
   },
   "text_data": "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."
}
```

# 5. Design

## 5.1. Requirements

| #     | Summary                                | Description                                                                                                 |
|-------|----------------------------------------|-------------------------------------------------------------------------------------------------------------|
| REQ-1 | TODO                                   | TODO                                                                                                        |

## 5.2. Approach

TODO

## 5.3. Limitations

| #     | Summary | Description |
|-------|---------|-------------|
| LIM-1 | TODO    | TODO        |

# 6. Contributing

## 6.1. Versioning

The service uses the [semantic versioning](http://semver.org/).
The single source of the version info is the git tag:
```shell
git describe --tags --abbrev=0
```

## 6.2. Issue Reporting

TODO

## 6.3. Building

```shell
make build
```
Generates the sources from proto files, compiles and creates the `router` executable.

## 6.4. Testing

### 6.4.1. Functional

```shell
make test
```

### 6.4.2. Performance

TODO

## 6.5. Releasing

To release a new version (e.g. `1.2.3`) it's enough to put a git tag:
```shell
git tag -v1.2.3
git push --tags
```

The corresponding CI job is started to build a docker image and push it with the specified tag (+latest).

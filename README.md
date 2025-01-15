# Gno Alerters

<img src="https://github.com/gnoverse/gno-alerters/blob/master/.github/assets/banner.png?raw=true" width="400" height="400"/>

A suite of alerting tools for monitoring and responding to events within the Gno blockchain ecosystem. Gno Alerters helps developers and operators maintain system reliability by detecting anomalies and notifying stakeholders in real-time.

## Features

* **Alerters**
    - Slack integration for real-time notifications.

* **Alerts**
    - Detect when the network is stalled for a defined period.
    - Identify validators who miss a consecutive number of defined blocks.


## Installation

### Prerequisites

- [Go](https://golang.org/) 1.22+
- A Gno blockchain node endpoint
- [Docker](https://www.docker.com/) (optional for containerized deployment)

### Clone the Repository

```bash
$ git clone https://github.com/gnoverse/gno-alerters.git
$ cd gno-alerters
```

### Build the Project

```bash
$ go build -o build/gno-alerter ./cmd/gno-alerter
# or
$ make build
```

### Run Tests

```bash
$ go test -v ./...
```

## Configuration

Gno Alerters uses a configuration file to define alerting rules and notification channels.

```yaml
[rpc]
endpoint = "http://localhost:26657"

[slack]
token      = "xoxb-4242"
channel_id = "XXXXXXX"

[alerts]
stalled_period = "30s"
consecutive_missed = [20, 100, 500]
```

## Usage

### Running Locally

Start the alerter with your configuration file:

```bash
$ ./gno-alerter -config config.yaml
```

---

Happy monitoring with **Gno Alerters**! 🚀


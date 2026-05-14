# URL Lookup Service

A small URL safety checker service that tells an HTTP proxy whether a requested URL is known to be malicious.

Run the service in one terminal, then send lookup requests from another terminal using `curl`. The service responds with a JSON verdict: `malicious` if the URL is present in the local blocklist, or `safe` if it is not.

## Contents

- [URL Lookup Service](#url-lookup-service)
  - [Contents](#contents)
  - [Prerequisites](#prerequisites)
  - [Quick Start](#quick-start)
  - [Blocked URLs](#blocked-urls)
  - [Testing](#testing)
  - [Run The Service](#run-the-service)
  - [Configuration](#configuration)
  - [Build](#build)
  - [Troubleshooting](#troubleshooting)
  - [How It Works](#how-it-works)
  - [Design](#design)
  - [Part 2: Discussion and Final Result](#part-2-discussion-and-final-result)
## Prerequisites

<details>
<summary>macOS</summary>

1. Install [Homebrew](https://brew.sh) if you don't have it:

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

2. Install Go 1.26 and Git:

```sh
brew install go@1.26 git
```

3. Verify:

```sh
go version
git --version
```

`curl` is pre-installed on macOS.

</details>

<details>
<summary>Linux (Debian/Ubuntu)</summary>

1. Install Go 1.26, Git, and curl:

```sh
sudo snap install go --channel=1.26/stable --classic
sudo apt update && sudo apt install git curl
```

2. Verify:

```sh
go version
git --version
```

</details>

## Quick Start

Clone the repository:

```sh
git clone https://github.com/mvillla/url-safety-checker.git
cd url-safety-checker
```

Start the service:

```sh
go run ./cmd/urlinfo
```

In another terminal, check a URL:

```sh
curl localhost:8080/urlinfo/1/malware.test/bad
```

To check any URL, paste it after `/urlinfo/1/` without the `http://` or `https://`:

| URL you want to check | Command |
| --- | --- |
| `http://malware.test/bad` | `curl localhost:8080/urlinfo/1/malware.test/bad` |
| `https://bad.example:443/download` | `curl localhost:8080/urlinfo/1/bad.example:443/download` |
| `https://google.com/search?q=test` | `curl 'localhost:8080/urlinfo/1/google.com/search?q=test'` |

Expected response:

```json
{"normalized_url":"malware.test/bad","verdict":"malicious","matched":true,"reason":"known malware URL"}
```

## Blocked URLs

The service loads its malware list from [`data/malware_urls.txt`](data/malware_urls.txt). Sample entries:

```text
malware.test/bad
malware.test/phishing/login
bad.example:443/download
phishing.example/password-reset
payload.invalid/install?os=macos
ransom.test/download/locker
```

## Testing

The smoke test is the primary way to validate the project. It runs the Go test
suite, starts the service on port `18080`, and checks live HTTP responses.

From the `url-safety-checker` directory:

```sh
cd url-safety-checker
./scripts/smoke-test.sh
```

Expected output:

```text
Running smoke checks...
PASS probe      healthz
PASS probe      readyz
PASS malicious  malware.test/bad
PASS malicious  bad.example:443/download
PASS safe       google.com/search?q=test
PASS safe       example.com/docs

Summary: 6 passed, 0 failed
```

If port `18080` is already in use:

```sh
PORT=19090 ./scripts/smoke-test.sh
```

To run the unit tests on their own:

```sh
go test ./...
```

## Run The Service

```sh
cd url-safety-checker
go run ./cmd/urlinfo
```

Expected startup logs:

```text
loaded 50 malware URLs from data/malware_urls.txt
starting url safety checker on :8080
```

The server listens on port `8080` by default. Example lookup:

```sh
curl localhost:8080/urlinfo/1/malware.test/bad
```

Shell note: quote URLs containing `?` when using shells such as `zsh`:

```sh
curl 'localhost:8080/urlinfo/1/google.com/search?q=test'
```

Stop the service with `Ctrl+C`.

## Configuration

| Name | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP port used by the service |
| `MALWARE_URLS_FILE` | `data/malware_urls.txt` | Local malware URL dataset |

Example:

```sh
cd url-safety-checker
PORT=9090 go run ./cmd/urlinfo
```

## Build

Build is optional. It creates a standalone binary so the service can run
without `go run`:

```sh
cd url-safety-checker
go build -o bin/urlinfo ./cmd/urlinfo
./bin/urlinfo
```

## Troubleshooting

<details>
<summary>zsh: no matches found</summary>

This happens when a URL contains `?` and is not quoted.

Use:

```sh
curl 'localhost:8080/urlinfo/1/google.com/search?q=test'
```

</details>

<details>
<summary>Port already in use</summary>

Use another port:

```sh
PORT=9090 go run ./cmd/urlinfo
```

Or for the smoke test:

```sh
cd url-safety-checker
PORT=19090 ./scripts/smoke-test.sh
```

</details>

<details>
<summary>Missing malware URL file</summary>

Make sure you are in the `url-safety-checker` directory, or set the path explicitly:

```sh
cd url-safety-checker
MALWARE_URLS_FILE=data/malware_urls.txt go run ./cmd/urlinfo
```

</details>

## How It Works

![URL lookup service request flow](docs/assets/url-lookup-flow.svg)

1. **HTTP Proxy** intercepts a request and sends a blocking lookup to the service and waits for a verdict before allowing or denying traffic.
2. **Lookup API** validates the HTTP method and parses the route (`GET /urlinfo/1/...`).
3. **URL Normalization** lowercases the host and preserves port, path, and query string to produce a scheme-less lookup key.
4. **Verdict Engine** does an exact-match check against the in-memory malware URL set and returns `"safe"` or `"malicious"`, always as HTTP 200.
5. The **Local Dataset File** (`data/malware_urls.txt`) is read at startup and loaded into the **Malware URL Set** as an in-memory map. No database needed.
6. The **JSON verdict** travels back to the proxy, which uses it to allow or block the request.

## Design

See [docs/design.md](docs/design.md) for the service design and architecture notes.

## Part 2: Discussion and Final Result

> The thought exercise answers covering production scaling, observability, deployment, and lifecycle are in [docs/discussion.md](docs/discussion.md).

# Design: URL Lookup Service

## 1. Problem

An HTTP proxy checks a requested URL before allowing access. This service
receives the URL lookup request and returns a safety verdict.

## 2. Requirements

- Expose `GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}`.
- Return whether the URL is known to be malicious.
- Return `200 OK` for both safe and malicious verdicts.
- Return `400` for malformed lookup requests.
- Return `405` for unsupported methods.
- Run locally on macOS and Linux.
- Include automated tests and run instructions.
- Include notes for scaling, operations, lifecycle, and deployment.

## 3. Project Structure

```text
cmd/urlinfo/
  main.go              server startup and dependency wiring

internal/httpapi/
  handler.go           HTTP routes, request parsing, JSON responses
  handler_test.go      handler tests

internal/lookup/
  normalize.go         URL key normalization
  normalize_test.go    normalization tests
  service.go           verdict model and lookup service
  store.go             store interface and in-memory/file-backed store
  store_test.go        store tests

data/
  malware_urls.txt     local malware URL list

docs/
  design.md            service design
  architecture.md      scaling and operational notes
```

The application has one executable entrypoint and two internal packages.
`httpapi` owns HTTP behavior. `lookup` owns URL normalization, verdicts, and
storage. `main` wires the components together.

## 4. Request Flow

```text
GET /urlinfo/1/example.com/path?x=1
  -> HTTP handler
  -> URL normalizer
  -> lookup service
  -> URL store
  -> JSON response
```

## 5. API Contract

Lookup endpoint:

```text
GET /urlinfo/1/{hostname_and_port}/{original_path_and_query_string}
```

Malicious response:

```json
{
  "normalized_url": "malware.test/bad",
  "verdict": "malicious",
  "matched": true,
  "reason": "known malware URL"
}
```

Safe response:

```json
{
  "normalized_url": "example.com/path",
  "verdict": "safe",
  "matched": false
}
```

Health endpoints:

```text
GET /healthz
GET /readyz
```

## 6. URL Normalization

Lookup keys are scheme-less.

Rules:

- Lowercase the host.
- Preserve the port.
- Preserve path case.
- Preserve query string case.
- Do not infer a scheme from the port.

Example:

```text
Example.COM:443/Path?Token=ABC -> example.com:443/Path?Token=ABC
```

## 7. Data Source

The service loads malware URLs from `data/malware_urls.txt`.

Format:

- one normalized URL per line
- empty lines ignored
- lines starting with `#` ignored

Example:

```text
# malware URLs
malware.test/bad
example.com:443/phishing?campaign=1
```

## 8. Error Handling

- Missing or unreadable malware URL file: startup error.
- Malformed lookup path: `400 Bad Request`.
- Unknown route: `404 Not Found`.
- Unsupported method: `405 Method Not Allowed`.
- Unexpected server error: `500 Internal Server Error`.

## 9. Configuration

Environment variables:

| Name | Default | Description |
| --- | --- | --- |
| `PORT` | `8080` | HTTP listen port |
| `MALWARE_URLS_FILE` | `data/malware_urls.txt` | Malware URL list path |

## 10. Testing

Planned coverage:

- malicious URL lookup
- safe URL lookup
- host normalization
- port preservation
- path preservation
- query preservation
- malformed lookup request
- unsupported method
- comment and blank line handling in the file loader

Command:

```sh
go test ./...
```

## 11. Scaling Considerations

The initial store loads a local file into memory. Lookup behavior depends on a
store interface, so later storage changes do not affect the HTTP contract.

The architecture notes will cover larger datasets, horizontal scaling, regional
deployment, update ingestion, operations, lifecycle, and deployment.

## 12. Delivery Plan

1. Create the Go module and package structure.
2. Implement URL normalization and tests.
3. Add verdict model, lookup service, and store.
4. Load malware URLs from the local file.
5. Add HTTP handlers and handler tests.
6. Add application startup and configuration.
7. Add run, test, and build instructions.
8. Add architecture notes for operations and deployment.

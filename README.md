# URL Safety Checker

URL lookup service for checking whether a requested URL is known to be
malicious.

The service exposes an HTTP endpoint for a proxy to query before allowing user
traffic. The application uses Go, the standard library, a local malware URL
file, and automated tests.

## Project Status

This repository currently contains the initial service design. Implementation
details will be added as the project progresses.

## Design

See [docs/design.md](docs/design.md) for the MVP design.


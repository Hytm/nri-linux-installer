# nri-linux-installer

nri-linux-installer is a tiny program that install New Relic Infrastructure agent on all supported Linux distribution automatically.

This is not supported or developed by New Relic.

## Installation

    go get -u github.com/hytm/nri-linux-installer

    You can also use the binary from this repository directly.

## Usage

nri-linux-installer detects the distribution and runs appropriate commands.

```
Usage: nri-linux-installer [options...] <license key>

Options:
  -mode  Installation mode (default as root). Use R for root, P for Privileged and U for Unprivileged.

```
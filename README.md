# Frankie

A command-line interface for [Frank Energie](https://www.frankenergie.nl/).

## Why

Frank Energie provides dynamic energy prices but doesn't offer an easy way to export historical usage and price data. Frankie lets you log this data locally for analysis, budgeting, or integration with other tools.

## Installation

```bash
go install github.com/pietern/frankie@latest
```

Or build from source:

```bash
go build -o frankie
```

## Usage

```bash
# Login to Frank Energie
frankie login

# View account summary
frankie summary

# View current prices
frankie prices

# View usage data
frankie usage

# View invoices
frankie invoices

# View connected sites, vehicles, chargers, batteries
frankie sites
frankie vehicles
frankie chargers
frankie batteries
frankie connections
```

### Output formats

```bash
# Table output (default)
frankie prices

# JSON output
frankie prices -o json
```

## Configuration

Configuration is stored in `~/.config/frankie/config.yaml`.

## Acknowledgements

The initial implementation of the API interactions was based on [python-frank-energie](https://github.com/HiDiHo01/python-frank-energie). Built with the assistance of [Claude Code](https://claude.ai/code).

# dist-plugin

Metric provider plugins for [dist](https://github.com/farishan-labs/dist).

Provides `gumroad` and `goatcounter` subcommands that print JSON envelopes consumed by `dist metrics`.

## Build

```sh
go build -o dist-plugin ./cmd/dist-plugin
```

## Usage

```sh
dist-plugin gumroad       # Gumroad products/sales/revenue (+ views via MCP)
dist-plugin goatcounter   # GoatCounter traffic
```

Each prints `{"metrics": {...}, "notes": [...]}`.

## Adding new providers

Create a new subcommand or a separate binary named `dist-plugin-<name>` that prints the same envelope; register it with `dist plugin add <name> <binary> [args...]`.
EOF
echo written